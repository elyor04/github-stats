package githubapi

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	rc *resty.Client
}

type Repo struct {
	Name            string `json:"name"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	Private         bool   `json:"private"`
	Fork            bool   `json:"fork"`
	Owner           struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type User struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	PublicRepos int    `json:"public_repos"`
}

type searchResult struct {
	TotalCount int `json:"total_count"`
}

func NewClient(token string) *Client {
	rc := resty.New().
		SetBaseURL("https://api.github.com").
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Accept", "application/vnd.github+json").
		SetTimeout(30 * time.Second)
	return &Client{rc: rc}
}

func (c *Client) GetAuthenticatedUser() (*User, error) {
	var user User
	resp, err := c.rc.R().SetResult(&user).Get("/user")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("could not get authenticated user: status %d", resp.StatusCode())
	}
	return &user, nil
}

func (c *Client) GetUser(username string) (*User, error) {
	var user User
	resp, err := c.rc.R().SetResult(&user).Get("/users/" + username)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("could not fetch user %s: status %d", username, resp.StatusCode())
	}
	return &user, nil
}

func (c *Client) ListReposForUser(username string) ([]Repo, error) {
	return c.paginateRepos("/users/"+username+"/repos", map[string]string{})
}

func (c *Client) ListReposForAuthenticatedUser() ([]Repo, error) {
	return c.paginateRepos("/user/repos", map[string]string{
		"affiliation": "owner",
		"visibility":  "all",
	})
}

func (c *Client) paginateRepos(path string, params map[string]string) ([]Repo, error) {
	var all []Repo
	page := 1
	for {
		var pageRepos []Repo
		req := c.rc.R().SetResult(&pageRepos).SetQueryParams(params).
			SetQueryParam("page", fmt.Sprintf("%d", page)).
			SetQueryParam("per_page", "100")

		resp, err := req.Get(path)
		if err != nil {
			return nil, err
		}
		if resp.IsError() {
			return nil, fmt.Errorf("github api error: status %d for %s", resp.StatusCode(), path)
		}

		all = append(all, pageRepos...)
		if len(pageRepos) < 100 {
			break
		}
		page++
	}
	return all, nil
}

func (c *Client) ListLanguages(owner, repo string) (map[string]int, error) {
	var languages map[string]int
	resp, err := c.rc.R().SetResult(&languages).Get(fmt.Sprintf("/repos/%s/%s/languages", owner, repo))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("could not fetch languages for %s: status %d", repo, resp.StatusCode())
	}
	return languages, nil
}

func (c *Client) SearchCount(path, query string) (int, error) {
	var result searchResult
	resp, err := c.rc.R().SetResult(&result).
		SetQueryParam("q", query).
		SetQueryParam("per_page", "1").
		Get(path)
	if err != nil {
		return 0, err
	}
	if resp.IsError() {
		return 0, fmt.Errorf("search error: status %d", resp.StatusCode())
	}
	return result.TotalCount, nil
}
