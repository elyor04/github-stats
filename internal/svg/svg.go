package svg

import (
	"fmt"
	"sort"
	"strings"

	"github-stats/internal/stats"
)

func GenerateStatsSVG(s *stats.UserStats) string {
	return fmt.Sprintf(`<svg width="465" height="195" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <style>
      .header { font: 600 18px 'Segoe UI', Ubuntu, Sans-Serif; fill: #fff; }
      .stat { font: 400 14px 'Segoe UI', Ubuntu, Sans-Serif; fill: #9f9f9f; }
      .stat-value { font: 600 14px 'Segoe UI', Ubuntu, Sans-Serif; fill: #fff; }
      .icon { fill: #79ff97; }
    </style>
  </defs>

  <rect width="465" height="195" fill="#151515" rx="4.5"/>

  <text x="25" y="35" class="header">%s's GitHub Stats</text>

  <!-- Stats -->
  <g transform="translate(0, 55)">
    <!-- Left column -->
    <g transform="translate(25, 0)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M8 .25a.75.75 0 01.673.418l1.882 3.815 4.21.612a.75.75 0 01.416 1.279l-3.046 2.97.719 4.192a.75.75 0 01-1.088.791L8 12.347l-3.766 1.98a.75.75 0 01-1.088-.79l.72-4.194L.818 6.374a.75.75 0 01.416-1.28l4.21-.611L7.327.668A.75.75 0 018 .25z"/>
      </svg>
      <text x="25" y="12" class="stat">Total Stars: <tspan class="stat-value">%d</tspan></text>
    </g>

    <g transform="translate(25, 30)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M5 3.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0zm0 2.122a2.25 2.25 0 10-1.5 0v.878A2.25 2.25 0 005.75 8.5h1.5v2.128a2.251 2.251 0 101.5 0V8.5h1.5a2.25 2.25 0 002.25-2.25v-.878a2.25 2.25 0 10-1.5 0v.878a.75.75 0 01-.75.75h-4.5A.75.75 0 015 6.25v-.878z"/>
      </svg>
      <text x="25" y="12" class="stat">Total Forks: <tspan class="stat-value">%d</tspan></text>
    </g>

    <g transform="translate(25, 60)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M2 2.5A2.5 2.5 0 014.5 0h8.75a.75.75 0 01.75.75v12.5a.75.75 0 01-.75.75h-2.5a.75.75 0 110-1.5h1.75v-2h-8a1 1 0 00-.714 1.7.75.75 0 01-1.072 1.05A2.495 2.495 0 012 11.5v-9zm10.5-1V9h-8c-.356 0-.694.074-1 .208V2.5a1 1 0 011-1h8z"/>
      </svg>
      <text x="25" y="12" class="stat">Total Repos: <tspan class="stat-value">%d</tspan></text>
    </g>

    <g transform="translate(25, 90)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
      </svg>
      <text x="25" y="12" class="stat">Total Commits: <tspan class="stat-value">%d</tspan></text>
    </g>

    <!-- Right column -->
    <g transform="translate(260, 0)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
      </svg>
      <text x="25" y="12" class="stat">Total PRs: <tspan class="stat-value">%d</tspan></text>
    </g>

    <g transform="translate(260, 30)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M8 9.5a1.5 1.5 0 100-3 1.5 1.5 0 000 3z"/>
        <path d="M8 0a8 8 0 100 16A8 8 0 008 0zM1.5 8a6.5 6.5 0 1113 0 6.5 6.5 0 01-13 0z"/>
      </svg>
      <text x="25" y="12" class="stat">Total Issues: <tspan class="stat-value">%d</tspan></text>
    </g>

    <g transform="translate(260, 60)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M5.5 3.5a2 2 0 100 4 2 2 0 000-4zM2 5.5a3.5 3.5 0 115.898 2.549 5.507 5.507 0 013.034 4.084.75.75 0 11-1.482.235 4.001 4.001 0 00-7.9 0 .75.75 0 01-1.482-.236A5.507 5.507 0 013.102 8.05 3.49 3.49 0 012 5.5zM11 4a.75.75 0 100 1.5 1.5 1.5 0 01.666 2.844.75.75 0 00-.416.672v.352a.75.75 0 00.574.73c1.2.289 2.162 1.2 2.522 2.372a.75.75 0 101.434-.44 5.01 5.01 0 00-2.56-3.012A3 3 0 0011 4z"/>
      </svg>
      <text x="25" y="12" class="stat">Followers: <tspan class="stat-value">%d</tspan></text>
    </g>

    <g transform="translate(260, 90)">
      <svg class="icon" y="0" width="16" height="16" viewBox="0 0 16 16">
        <path d="M5.5 3.5a2 2 0 100 4 2 2 0 000-4zM2 5.5a3.5 3.5 0 115.898 2.549 5.507 5.507 0 013.034 4.084.75.75 0 11-1.482.235 4.001 4.001 0 00-7.9 0 .75.75 0 01-1.482-.236A5.507 5.507 0 013.102 8.05 3.49 3.49 0 012 5.5zM11 4a.75.75 0 100 1.5 1.5 1.5 0 01.666 2.844.75.75 0 00-.416.672v.352a.75.75 0 00.574.73c1.2.289 2.162 1.2 2.522 2.372a.75.75 0 101.434-.44 5.01 5.01 0 00-2.56-3.012A3 3 0 0011 4z"/>
      </svg>
      <text x="25" y="12" class="stat">Following: <tspan class="stat-value">%d</tspan></text>
    </g>
  </g>
</svg>`, s.Name, s.TotalStars, s.TotalForks, s.TotalRepos, s.TotalCommits,
		s.TotalPRs, s.TotalIssues, s.Followers, s.Following)
}

var languageColors = map[string]string{
	"JavaScript": "#f1e05a",
	"TypeScript": "#3178c6",
	"Python":     "#3572A5",
	"Java":       "#b07219",
	"C":          "#555555",
	"C++":        "#f34b7d",
	"C#":         "#178600",
	"Go":         "#00ADD8",
	"Rust":       "#dea584",
	"Ruby":       "#701516",
	"PHP":        "#4F5D95",
	"Swift":      "#ffac45",
	"Kotlin":     "#A97BFF",
	"Dart":       "#00B4AB",
	"HTML":       "#e34c26",
	"CSS":        "#563d7c",
	"Shell":      "#89e051",
}

func GenerateLanguagesSVG(languages stats.LanguageStats) string {
	type kv struct {
		lang    string
		percent float64
	}
	entries := make([]kv, 0, len(languages))
	for lang, pct := range languages {
		entries = append(entries, kv{lang, pct})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].percent > entries[j].percent })
	if len(entries) > 8 {
		entries = entries[:8]
	}

	width := 300
	height := 45 + len(entries)*40

	var rows strings.Builder
	for i, e := range entries {
		color := languageColors[e.lang]
		if color == "" {
			color = "#858585"
		}
		fmt.Fprintf(&rows, `
      <g transform="translate(25, %d)">
        <circle cx="6" cy="8" r="6" fill="%s"/>
        <text x="20" y="12" class="lang-name">%s</text>
        <text x="%d" y="12" class="percentage">%v%%</text>
      </g>`, i*40, color, e.lang, width-80, e.percent)
	}

	return fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <style>
      .header { font: 600 18px 'Segoe UI', Ubuntu, Sans-Serif; fill: #fff; }
      .lang-name { font: 400 14px 'Segoe UI', Ubuntu, Sans-Serif; fill: #9f9f9f; }
      .percentage { font: 400 12px 'Segoe UI', Ubuntu, Sans-Serif; fill: #9f9f9f; }
    </style>
  </defs>

  <rect width="%d" height="%d" fill="#151515" rx="4.5"/>

  <text x="25" y="35" class="header">Most Used Languages</text>

  <g transform="translate(0, 50)">%s
  </g>
</svg>`, width, height, width, height, rows.String())
}
