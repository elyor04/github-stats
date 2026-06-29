# github-stats

A small HTTP service that generates GitHub profile stats and language usage as SVG images (and JSON), similar to GitHub readme stats cards.

This is the original **Deno/TypeScript** implementation. A Go rewrite of the same service lives on the [`golang`](https://github.com/elyor04/github-stats/tree/golang) branch.

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
deno task start
```

## Running with Docker Compose

```bash
docker compose up
```

## Caching

Responses are cached in memory for one hour per `(username, visibility)` pair to avoid hitting GitHub API rate limits. A `Deno.cron` job refreshes the cache for the default user (`elyor04`, including private repos) every 5 minutes.
