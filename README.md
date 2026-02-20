# ticktick-cli

A Go-based TickTick CLI inspired by the command-oriented style of `gogcli`.

## Features

- OAuth setup and token exchange
- Automatic token refresh (when refresh token is available)
- All documented TickTick OpenAPI `/open/v1` endpoints
- Project list/get/create/update/delete
- Task list/get/create/update/complete/delete
- Human-readable output or `--json`

## Install

1. Install Go 1.22+
2. Build:

```bash
go build -o ticktick ./cmd/ticktick
```

3. Install globally (pick one):

```bash
go install github.com/apktdev/ticktick-cli/cmd/ticktick@latest
```

or

```bash
sudo install -m 0755 ./ticktick /usr/local/bin/ticktick
```

## CI Builds

GitHub Actions builds binaries automatically on every push to `main`, tags (`v*`), pull requests, and manual dispatch.

- Open the repository's **Actions** tab.
- Select the latest **Build** workflow run.
- Download the artifact matching your platform (for example `ticktick-linux-amd64`).

## Releases (Public HTTP Downloads)

Version tags automatically publish GitHub Release assets (stable URLs, no `gh` CLI required).

1. Create and push a tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

2. GitHub Actions runs the **Release** workflow and publishes assets on:
   - `https://github.com/apktdev/ticktick-cli/releases/tag/v0.1.0`

3. Download directly via URL pattern:
   - `https://github.com/apktdev/ticktick-cli/releases/download/v0.1.0/ticktick-linux-amd64.tar.gz`

## OAuth setup

Create a TickTick app and get:

- `client_id`
- `client_secret`
- `redirect_uri`

Then:

```bash
./ticktick auth set-client --client-id <id> --client-secret <secret> --redirect-uri <uri>
./ticktick auth login-url
```

Open the URL, authorize, then copy the returned `code` value:

```bash
./ticktick auth exchange --code <oauth_code>
./ticktick auth status
```

## Commands

```bash
./ticktick projects list
./ticktick projects get <project_id>
./ticktick projects add --name "Inbox 2" --color "#F18181" --view-mode list --kind TASK
./ticktick projects update <project_id> --name "Renamed"
./ticktick projects delete <project_id>
./ticktick tasks list --project-id <project_id>
./ticktick tasks get <project_id> <task_id>
./ticktick tasks add --project-id <project_id> --title "Buy milk" --due 2026-02-16
./ticktick tasks update <task_id> --project-id <project_id> --title "Buy oat milk"
./ticktick tasks complete <project_id> <task_id>
./ticktick tasks delete <project_id> <task_id>
```

Use JSON mode when scripting:

```bash
./ticktick --json projects list
```

## Notes

- Config is stored at your OS config dir, typically `~/.config/tickcli/config.json`.
- You can run fully from environment variables instead of config file:
  - `TICKTICK_ACCESS_TOKEN`
  - `TICKTICK_CLIENT_ID`
  - `TICKTICK_CLIENT_SECRET`
  - `TICKTICK_REDIRECT_URI`
  - `TICKTICK_REFRESH_TOKEN`
  - `TICKTICK_TOKEN_TYPE`
  - `TICKTICK_SCOPE`
  - `TICKTICK_TOKEN_EXPIRY` (RFC3339, e.g. `2026-08-15T14:47:20Z`)
- When any env var above is set, env values override file config and config is not auto-saved.
- Date flags (`--start`, `--due`) accept RFC3339 or `YYYY-MM-DD`.
- Official endpoints covered:
  - `GET /open/v1/project`
  - `GET /open/v1/project/{projectId}`
  - `GET /open/v1/project/{projectId}/data`
  - `POST /open/v1/project`
  - `POST /open/v1/project/{projectId}`
  - `DELETE /open/v1/project/{projectId}`
  - `GET /open/v1/project/{projectId}/task/{taskId}`
  - `POST /open/v1/task`
  - `POST /open/v1/task/{taskId}`
  - `POST /open/v1/project/{projectId}/task/{taskId}/complete`
  - `DELETE /open/v1/project/{projectId}/task/{taskId}`
