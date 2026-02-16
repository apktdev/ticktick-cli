# ticktick-cli

A Go-based TickTick CLI inspired by the command-oriented style of `gogcli`.

## Features (MVP)

- OAuth setup and token exchange
- Automatic token refresh (when refresh token is available)
- Project listing
- Task listing by project
- Task create / complete / delete
- Human-readable output or `--json`

## Install

1. Install Go 1.22+
2. Build:

```bash
go build -o tick ./cmd/tick
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
./tick tasks list --project-id <project_id>
./tick tasks add --project-id <project_id> --title "Buy milk" --due 2026-02-16
./tick tasks complete <project_id> <task_id>
./tick tasks delete <project_id> <task_id>
```

Use JSON mode when scripting:

```bash
./tick --json projects list
```

## Notes

- Config is stored at your OS config dir, typically `~/.config/tickcli/config.json`.
- Due dates accept RFC3339 or `YYYY-MM-DD`.
- This is an initial foundation; additional commands (tags, sections, batch update, sync helpers) can be layered on top.
