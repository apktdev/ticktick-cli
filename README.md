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
go build -o tick ./cmd/tick
```

3. Install globally (pick one):

```bash
go install github.com/apktdev/ticktick-cli/cmd/tick@latest
```

or

```bash
sudo install -m 0755 ./tick /usr/local/bin/tick
```

## OAuth setup

Create a TickTick app and get:

- `client_id`
- `client_secret`
- `redirect_uri`

Then:

```bash
./tick auth set-client --client-id <id> --client-secret <secret> --redirect-uri <uri>
./tick auth login-url
```

Open the URL, authorize, then copy the returned `code` value:

```bash
./tick auth exchange --code <oauth_code>
./tick auth status
```

## Commands

```bash
./tick projects list
./tick projects get <project_id>
./tick projects add --name "Inbox 2" --color "#F18181" --view-mode list --kind TASK
./tick projects update <project_id> --name "Renamed"
./tick projects delete <project_id>
./tick tasks list --project-id <project_id>
./tick tasks get <project_id> <task_id>
./tick tasks add --project-id <project_id> --title "Buy milk" --due 2026-02-16
./tick tasks update <task_id> --project-id <project_id> --title "Buy oat milk"
./tick tasks complete <project_id> <task_id>
./tick tasks delete <project_id> <task_id>
```

Use JSON mode when scripting:

```bash
./tick --json projects list
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
