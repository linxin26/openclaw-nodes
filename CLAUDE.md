# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
go build -o openclaw-node.exe ./cmd      # Build Windows executable
go test ./...                              # Run all tests
go test -v ./internal/protocol/           # Run protocol tests with verbose output
golangci-lint run --no-config --disable-all -E errcheck ./...  # Lint (v1.48.0 required for Go 1.20)
```

## Architecture Overview

OpenClaw Node is a desktop agent that connects to a gateway via WebSocket and exposes device capabilities (camera, location, screen, etc.) through a command invocation protocol.

### Core Flow

```
cmd/main.go → protocol.Client (WebSocket) → gateway
                    ↓
              protocol.Dispatch()
                    ↓
         [device.* handlers in internal/device/]
```

### Key Packages

- **`internal/protocol`** - WebSocket client, protocol frames, command dispatch
  - `client.go` - Connection management, challenge-response handshake, ping/pong
  - `invoke.go` - Command dispatch via `protocol.Handlers` map
  - Commands registered via `protocol.RegisterHandler()` in each device package

- **`internal/crypto`** - Ed25519 identity for device authentication
  - Generates keypairs, signs auth payloads, exports public key base64

- **`internal/device/`** - Device capability handlers (camera, location, screen, etc.)
  - Each file registers its handlers via `init()` (e.g., `camera.go` → `camera.list`, `camera.snap`)
  - **External dependencies**: `ffmpeg` (camera listing/capture), `powershell` (screen capture via System.Windows.Forms)

- **`internal/discovery`** - mDNS service registration (currently simplified TCP listen on port 18789)

- **`internal/tray`** - System tray integration with status indicators

- **`store`** - Data directory management (`%APPDATA%\OpenClaw`)

### Protocol Details

- **Authentication**: Challenge-response using Ed25519 signatures
  1. Gateway sends `connect.challenge` with nonce
  2. Node signs auth payload with Ed25519 private key
  3. Node sends `connect` request with signature
- **Frames**: JSON with `type` (`req`, `res`, `event`), `id`, `method`, `params`, `payload`
- **Invoke**: Bidirectional - gateway can invoke commands on node via `node.invoke.request`

### Data Storage

- `identity.json` - Ed25519 keypair (stored in `%APPDATA%\OpenClaw`)
- `config.yaml` - Gateway address, TLS settings, discovery mode, capabilities

### Flags

- `-gateway` - Gateway address (host:port)
- `-token` - Gateway auth token
- `-tls` - Use TLS connection
- `-no-mdns` - Disable mDNS discovery

## Git Conventions

**Commit messages**: [Conventional Commits](https://www.conventionalcommits.org/)
```
<type>: <description>

<optional body>
```
Types: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `perf`, `ci`

Example:
```
feat: add Ed25519 identity layer for device authentication
```
