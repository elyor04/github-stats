package stats

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github-stats/internal/cache"
	"github-stats/internal/githubapi"

	"go.uber.org/zap"
)

type UserStats struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	TotalStars   int    `json:"totalStars"`
	TotalForks   int    `json:"totalForks"`
	TotalRepos   int    `json:"totalRepos"`
	PublicRepos  int    `json:"publicRepos"`
	PrivateRepos int    `json:"privateRepos"`
	Followers    int    `json:"followers"`
	Following    int    `json:"following"`
	TotalCommits int    `json:"totalCommits"`
	TotalIssues  int    `json:"totalIssues"`
	TotalPRs     int    `json:"totalPRs"`
}

type LanguageStats map[string]float64

type Service struct {
	client *githubapi.Client
	cache  *cache.Cache
	log    *zap.Logger
}

func NewService(client *githubapi.Client, log *zap.Logger) *Service {
	return &Service{client: client, cache: cache.New(time.Hour), log: log}
}

func (s *Service) isAuthenticatedUser(username string) bool {
	cacheKey := "auth-user"
	var authUser *githubapi.User
	if cached, ok := s.cache.Get(cacheKey); ok {
		u := cached.(githubapi.User)
		authUser = &u
	} else {
		u, err := s.client.GetAuthenticatedUser()
		if err != nil {
			s.log.Warn("could not get authenticated user", zap.Error(err))
			return false
		}
		s.cache.Set(cacheKey, *u)
		authUser = u
	}
	return strings.EqualFold(authUser.Login, username)
}

func (s *Service) listRepos(username string, includePrivate bool) ([]githubapi.Repo, error) {
	if includePrivate && s.isAuthenticatedUser(username) {
		return s.client.ListReposForAuthenticatedUser()
	}
	return s.client.ListReposForUser(username)
}

func (s *Service) GetUserStats(username string, includePrivate bool) (*UserStats, error) {
	cacheKey := fmt.Sprintf("user-stats-%s-%s", username, visibilityKey(includePrivate))
	if cached, ok := s.cache.Get(cacheKey); ok {
		stats := cached.(UserStats)
		return &stats, nil
	}

	user, err := s.client.GetUser(username)
	if err != nil {
		return nil, err
	}

	repos, err := s.listRepos(username, includePrivate)
	if err != nil {
		return nil, err
	}

	totalStars, totalForks, publicRepos, privateRepos := 0, 0, 0, 0
	for _, r := range repos {
		totalStars += r.StargazersCount
		totalForks += r.ForksCount
		if r.Private {
			privateRepos++
		} else {
			publicRepos++
		}
	}

	totalCommits := 0
	if count, err := s.client.SearchCount("/search/commits", "author:"+username); err != nil {
		s.log.Warn("could not fetch commits", zap.Error(err))
	} else {
		totalCommits = count
	}

	totalIssues, totalPRs := 0, 0
	if count, err := s.client.SearchCount("/search/issues", "author:"+username+" type:issue"); err != nil {
		s.log.Warn("could not fetch issues", zap.Error(err))
	} else {
		totalIssues = count
	}
	if count, err := s.client.SearchCount("/search/issues", "author:"+username+" type:pr"); err != nil {
		s.log.Warn("could not fetch PRs", zap.Error(err))
	} else {
		totalPRs = count
	}

	finalPublicRepos := user.PublicRepos
	finalPrivateRepos := 0
	if includePrivate {
		finalPublicRepos = publicRepos
		finalPrivateRepos = privateRepos
	}

	name := user.Name
	if name == "" {
		name = user.Login
	}

	result := UserStats{
		Username:     user.Login,
		Name:         name,
		TotalStars:   totalStars,
		TotalForks:   totalForks,
		TotalRepos:   len(repos),
		PublicRepos:  finalPublicRepos,
		PrivateRepos: finalPrivateRepos,
		Followers:    user.Followers,
		Following:    user.Following,
		TotalCommits: totalCommits,
		TotalIssues:  totalIssues,
		TotalPRs:     totalPRs,
	}

	s.cache.Set(cacheKey, result)
	return &result, nil
}

func (s *Service) GetLanguageStats(username string, includePrivate bool) (LanguageStats, error) {
	cacheKey := fmt.Sprintf("lang-stats-%s-%s", username, visibilityKey(includePrivate))
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(LanguageStats), nil
	}

	repos, err := s.listRepos(username, includePrivate)
	if err != nil {
		return nil, err
	}

	totals := map[string]int{}
	for _, r := range repos {
		if r.Fork {
			continue
		}
		langs, err := s.client.ListLanguages(r.Owner.Login, r.Name)
		if err != nil {
			s.log.Warn("could not fetch languages for repo", zap.String("repo", r.Name), zap.Error(err))
			continue
		}
		for lang, bytes := range langs {
			totals[lang] += bytes
		}
	}

	result := computePercentages(totals)
	s.cache.Set(cacheKey, result)
	return result, nil
}

func computePercentages(totals map[string]int) LanguageStats {
	total := 0
	for _, bytes := range totals {
		total += bytes
	}

	type kv struct {
		lang    string
		percent float64
	}
	var sorted []kv
	if total > 0 {
		for lang, bytes := range totals {
			pct, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(bytes)/float64(total)*100), 64)
			sorted = append(sorted, kv{lang, pct})
		}
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].percent > sorted[j].percent })

	result := LanguageStats{}
	for _, item := range sorted {
		result[item.lang] = item.percent
	}
	return result
}

func visibilityKey(includePrivate bool) string {
	if includePrivate {
		return "private"
	}
	return "public"
}
