# Numbernama Go

A small **number-matching board game** backend written in Go, migrated from the original Node.js + Socket.IO service in this repository. Game state lives **in RAM per WebSocket connection**; when the socket disconnects, that session is discarded. HTTP routes are wired with **[Gorilla mux](https://github.com/gorilla/mux)**; real-time play uses **[go-socket](https://github.com/dibyaranjan-pradhan/go-socket)** (JSON events over WebSocket, built on `gorilla/websocket`).

## The game (short)

- You get a grid of digits (modes **1–18** or **1–19**).
- Pick **two cells** whose numbers are **equal** or **sum to 10**.
- They may clear only if they match one of the **path rules** from the original game (same column with empty cells between, same row with empty cells between, or the “L-shaped empty corridor” case described in the legacy `utils/numbers-server.js` comments).
- **Check** reflows remaining digits into rows of nine, pulling from the bottom.
- **Clear** removes completely empty full rows.
- **Undo** restores the last successful pair removal.
- **Restart** keeps the same mode and rebuilds the starting board.

Open **`/numbers-game`** in a browser after starting the server; the UI and script are embedded from `web/`.

## Project layout

| Path | Role |
|------|------|
| `cmd/main.go` | `ListenAndServe` bootstrap |
| `cmd/wire.go` | Wire injector (`//go:build wireinject`) |
| `cmd/wire_gen.go` | Generated composition root (check in after `make wire`) |
| `cmd/providers.go` | zerolog + go-socket + router wiring |
| `middleware/middleware.go` | shared HTTP middleware registration |
| `router/gameplay.go` | `/ws/gameplay`, `/static/*`, `/numbers-game`, `/health` |
| `router/user.go` | `/api/user/me` route registration |
| `router/user_setting.go` | `/api/settings` route registration |
| `handler/gameplay_socket.go` | go-socket gameplay handlers (`initiateGamePlay`, `match`, `check`, `clear`, `undo`) |
| `handler/gameplay.go` | gameplay HTTP handlers |
| `handler/user.go`, `handler/user_setting.go` | sanity HTTP handlers |
| `service/gameplay.go` | core gameplay business logic (ported from Node) |
| `service/user.go`, `service/user_setting.go` | stubs for future persistence |
| `repo/memory_gameplay.go` | `map[clientID]*GameState` RAM store |
| `model/*.go` | feature structs and payload DTOs |
| `utils/logger.go` | zerolog wrapper + go-socket diag adapter |
| `web/` | `index.html`, `app.js`, `embed.go` |

Dependency injection uses **[Google Wire](https://github.com/google/wire)**. After changing constructor graphs in `cmd/wire.go`, run:

```bash
make wire
```

## WebSocket wire format (go-socket)

Client → server (text frame, JSON):

```json
{ "event": "match", "payload": [[0,0],[1,0]] }
```

Server → client:

```json
{ "event": "match", "payload": { "matched": true, "selectedElems": [[0,0],[1,0]] } }
```

On connect, the server emits **`go_gameplay_connected`** with the connection id string as `payload` (replacing the old `node_gameplay_connected` name so you can tell Node vs Go clients apart).

`clear` and `undo` payloads mirror the old tuple responses: `[boolean, data]`.

## Run locally

```bash
cd numbernama-go
go mod tidy
go run ./cmd
```

Defaults to **port 7002** (same as the Node app). Override with `PORT`.

- UI: [http://localhost:7002/numbers-game](http://localhost:7002/numbers-game)
- Health: [http://localhost:7002/health](http://localhost:7002/health)
- WebSocket: `ws://localhost:7002/ws/gameplay`

## Docker

```bash
make docker-build
docker run --rm -p 7002:7002 numbernama-go:local
```

## Module path

The Go module is declared as **`numbernama-go`** so the tree builds without a remote Git import path. If you publish this module, change the first line of `go.mod` to your canonical path (for example `module github.com/you/numbernama/numbernama-go`) and replace imports accordingly.

## Next steps (not in this minimal port)

- Persist users, settings, scores, and match history.
- Share `Board` state in a room for multiplayer.
- Harden CORS and auth on `/ws/gameplay` using middleware that wraps `s.Handler()`.
