# Go-Nemo

This is a Go implementation of the Nemo project, migrated from Java.

## Modules

- **User Service**: Handles user initialization and login.
- **EntLive Service**: Handles entertainment live rooms (creation, list, info, close).
- **SocialChat Service**: Handles 1v1 social chat, reporter, and rewards.
- **Game Service**: Handles game rooms (create, join, start, end).

## Prerequisites

- Go 1.18+
- MySQL
- Redis

## Configuration

Configuration is located in `config/config.yaml`.
You can override settings using environment variables.

## Running

```bash
go run cmd/server/main.go
```

## API Endpoints

### User

- POST `/nemo/app/initAppAndUser`

### Entertainment Live

- POST `/nemo/entertainmentLive/createLive`
- POST `/nemo/entertainmentLive/closeLive`
- GET `/nemo/entertainmentLive/list`
- GET `/nemo/entertainmentLive/info`

### Social Chat

- POST `/nemo/socialChat/user/reporter`
- GET `/nemo/socialChat/user/getOnLineUser`
- POST `/nemo/socialChat/user/reward`

### Game

- GET `/nemo/game/list`
- POST `/nemo/game/create`
- POST `/nemo/game/join`
- POST `/nemo/game/start`
- POST `/nemo/game/end`
