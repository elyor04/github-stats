# github-stats

A small HTTP service that generates GitHub profile stats and language usage as SVG images (and JSON), similar to GitHub readme stats cards.

This is the **Go** rewrite of the service. The original Deno/TypeScript implementation lives on the [`deno`](https://github.com/elyor04/github-stats/tree/deno) branch.

## Endpoints

| Path             | Description                                  |
| ---------------- | --------------------------------------------- |
| `/stats`         | Stats card as SVG                             |
| `/languages`     | Most used languages card as SVG               |
| `/api/stats`     | Stats as JSON                                 |
| `/api/languages` | Language percentages as JSON                  |

All endpoints accept:

- `username` (optional, defaults to `elyor04`)
- `private` (optional, `true`/`false`) — include private repos; only takes effect when `username` matches the authenticated user for `GITHUB_TOKEN`

Example:

```
GET /stats?username=octocat
GET /api/languages?username=octocat&private=true
```

## Configuration

Copy `.env.example` to `.env` and fill in the values, or set the equivalent environment variables directly.

| Variable      | Required | Default | Description                                  |
| ------------- | -------- | ------- | --------------------------------------------- |
| `GITHUB_TOKEN`| yes      | —       | GitHub personal access token                  |
| `PORT`        | no       | `8000`  | Port the HTTP server listens on               |

## Running locally

```bash
go run .
```

## Running with Docker Compose

```bash
docker compose up
```

## Running tests

```bash
go test ./...
```

## Project layout

```
main.go                    entrypoint: logger, config, cron job, HTTP server
internal/config             env/.env loading (Viper)
internal/cache               in-memory TTL cache used to avoid GitHub rate limits
internal/githubapi           GitHub REST API client (Resty)
internal/stats               stats/language aggregation logic
internal/svg                 SVG card rendering
internal/server              HTTP routes and handlers (Gin)
```

## Libraries

- [Gin](https://github.com/gin-gonic/gin) — HTTP routing
- [Resty](https://github.com/go-resty/resty) — GitHub API HTTP client
- [Viper](https://github.com/spf13/viper) — configuration
- [Zap](https://github.com/uber-go/zap) — structured logging
- [robfig/cron](https://github.com/robfig/cron) — periodic cache refresh for the default user
- [Testify](https://github.com/stretchr/testify) — unit tests

## Caching

Responses are cached in memory for one hour per `(username, visibility)` pair to avoid hitting GitHub API rate limits. A cron job refreshes the cache for the default user (`elyor04`, including private repos) every 5 minutes.
