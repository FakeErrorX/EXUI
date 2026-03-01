<div align="center">

# EXUI Documentation

**Full technical reference for the EXUI Xray management panel**

[Quick Start](#quick-start) · [API Reference](#api-reference) · [Configuration](#configuration) · [Docker](#docker-deployment) · [Telegram Bot](#telegram-bot) · [Development](#development)

</div>

---

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Environment Variables](#environment-variables)
- [Panel Settings](#panel-settings)
- [API Reference](#api-reference)
  - [Authentication](#authentication)
  - [Inbound Endpoints](#inbound-endpoints)
  - [Server Endpoints](#server-endpoints)
- [Subscription System](#subscription-system)
- [Telegram Bot](#telegram-bot)
- [Docker Deployment](#docker-deployment)
- [Database](#database)
- [Project Structure](#project-structure)
- [Development](#development)

---

## Overview

EXUI is an open-source, web-based management panel built on top of [Xray-core](https://github.com/XTLS/Xray-core). It provides a full-featured GUI and REST API for configuring inbound proxies, managing clients, monitoring traffic, and automating notifications via Telegram.

| Property | Value |
|---|---|
| Language | Go 1.24+ |
| Web Framework | Gin |
| Database | SQLite (via GORM) |
| Xray Core | v1.250608.0+ |
| Default Port | `2053` |
| Default Base Path | `/` |

---

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                     EXUI Server                     │
│                                                      │
│  ┌──────────────┐   ┌──────────────────────────────┐ │
│  │  Gin Router  │──▶│         Controllers           │ │
│  └──────────────┘   │  ┌──────────────────────────┐│ │
│                      │  │  IndexController  (auth)  ││ │
│  ┌──────────────┐   │  │  ServerController (sys)   ││ │
│  │  Middleware  │   │  │  EXUIController   (panel) ││ │
│  │  - Auth      │   │  │  APIController    (api)   ││ │
│  │  - Sessions  │   │  └──────────────────────────┘│ │
│  │  - i18n      │   └──────────────────────────────┘ │
│  │  - GZIP      │                  │                  │
│  │  - Domain    │   ┌──────────────▼───────────────┐ │
│  └──────────────┘   │           Services            │ │
│                      │  InboundService / XrayService │ │
│  ┌──────────────┐   │  SettingService / TgbotService│ │
│  │  Cron Jobs   │   │  ServerService / UserService  │ │
│  │  - Xray watchdog  └──────────────┬───────────────┘ │
│  │  - Traffic stats        │                          │
│  │  - Client IP check  ┌───▼──────────────────────┐  │
│  │  - Tg notifications │     SQLite Database       │  │
│  │  - Log cleanup      │  Inbounds / Clients /     │  │
│  └──────────────┘      │  Settings / Users         │  │
│                         └──────────────────────────┘  │
└─────────────────────────────────────────────────────┘
```

### Key Components

| Component | Path | Responsibility |
|---|---|---|
| Web Server | `web/web.go` | Router init, TLS, cron scheduler |
| Controllers | `web/controller/` | HTTP handlers, request validation |
| Services | `web/service/` | Business logic layer |
| Xray Process | `xray/process.go` | Xray lifecycle management |
| Xray API | `xray/api.go` | gRPC communication with Xray |
| Database | `database/db.go` | GORM + SQLite, auto-migration |
| Jobs | `web/job/` | Scheduled background tasks |
| Sub system | `sub/` | Subscription link generation |

---

## Quick Start

### One-line Install (Linux)

```bash
bash <(curl -Ls https://raw.githubusercontent.com/fakeerrorx/EXUI/master/install.sh)
```

The installer will:
1. Detect OS and CPU architecture
2. Check GLIBC ≥ 2.32
3. Download latest release binary
4. Set up systemd service `ex-ui`
5. Generate random admin credentials and base path

### Management Commands

```bash
ex-ui start       # Start service
ex-ui stop        # Stop service
ex-ui restart     # Restart service
ex-ui status      # Show status
ex-ui enable      # Enable on boot
ex-ui disable     # Disable on boot
ex-ui log         # Show logs
ex-ui update      # Update to latest
ex-ui install     # Re-install
ex-ui uninstall   # Remove
```

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `XUI_LOG_LEVEL` | `info` | Log level: `debug`, `info`, `notice`, `warn`, `error` |
| `XUI_DEBUG` | `false` | Enable debug mode (`true`/`false`) |
| `XUI_BIN_FOLDER` | `bin` | Path to folder containing Xray binary |
| `XUI_DB_FOLDER` | `/etc/ex-ui` | Path to SQLite database folder |
| `XRAY_VMESS_AEAD_FORCED` | — | Force VMess AEAD encryption |

---

## Panel Settings

All settings are stored in the database and configurable via the **Settings** page or directly via API. Below are all recognized setting keys with their defaults.

### Web / Panel

| Key | Default | Description |
|---|---|---|
| `webListen` | `` | Bind address (empty = all interfaces) |
| `webDomain` | `` | Restrict panel to this domain only |
| `webPort` | `2053` | Panel HTTP port |
| `webBasePath` | `/` | URL prefix for all panel routes |
| `webCertFile` | `` | Path to TLS certificate file |
| `webKeyFile` | `` | Path to TLS private key file |
| `sessionMaxAge` | `60` | Session lifetime in minutes |
| `pageSize` | `50` | Items per page in tables |
| `timeLocation` | `Local` | Timezone for display |

### Client / Traffic

| Key | Default | Description |
|---|---|---|
| `expireDiff` | `0` | Days before expiry to send warning |
| `trafficDiff` | `0` | Traffic (GB) threshold for warning |
| `remarkModel` | `-ieo` | Client remark pattern |

### Two-Factor Authentication

| Key | Default | Description |
|---|---|---|
| `twoFactorEnable` | `false` | Enable TOTP 2FA on login |
| `twoFactorToken` | `` | TOTP secret token |

### Subscription Service

| Key | Default | Description |
|---|---|---|
| `subEnable` | `false` | Enable subscription endpoint |
| `subPort` | `2096` | Subscription server port |
| `subPath` | `/sub/` | Subscription URL path |
| `subDomain` | `` | Domain for subscription links |
| `subCertFile` | `` | TLS cert for subscription server |
| `subKeyFile` | `` | TLS key for subscription server |
| `subUpdates` | `12` | Subscription update interval (hours) |
| `subEncrypt` | `true` | Encrypt subscription content |
| `subShowInfo` | `true` | Show traffic/expiry info in subscription |
| `subURI` | `` | Override subscription base URI |
| `subTitle` | `` | Subscription profile title |
| `subJsonPath` | `/json/` | JSON subscription path |
| `subJsonURI` | `` | Override JSON subscription URI |
| `subJsonFragment` | `` | Fragment settings for JSON sub |
| `subJsonMux` | `` | Mux settings for JSON sub |
| `subJsonRules` | `` | Routing rules for JSON sub |
| `subJsonNoises` | `` | DNS noise settings for JSON sub |
| `datepicker` | `gregorian` | Calendar style (`gregorian`/`jalali`) |

### Telegram Bot

| Key | Default | Description |
|---|---|---|
| `tgBotEnable` | `false` | Enable Telegram bot |
| `tgBotToken` | `` | Bot API token |
| `tgBotChatId` | `` | Comma-separated admin chat IDs |
| `tgBotProxy` | `` | SOCKS5 proxy for bot (e.g. `socks5://host:port`) |
| `tgBotAPIServer` | `` | Custom Telegram API server URL |
| `tgRunTime` | `@daily` | Cron expression for stats notification |
| `tgBotBackup` | `false` | Send DB backup via bot on schedule |
| `tgBotLoginNotify` | `true` | Notify admins on panel login |
| `tgCpu` | `80` | CPU usage (%) threshold for alert |

### External Traffic Reporting

| Key | Default | Description |
|---|---|---|
| `externalTrafficInformEnable` | `false` | Enable webhook on traffic exhaustion |
| `externalTrafficInformURI` | `` | Webhook URL to POST traffic events |

### WARP

| Key | Default | Description |
|---|---|---|
| `warp` | `` | Warp configuration JSON |

---

## API Reference

All API endpoints are under the **base path** (default `/`) + `/panel/api/inbounds/`.

### Authentication

The API uses **cookie-based session authentication** — the same session created by the web panel login. To use the API programmatically:

**Step 1 — Login**

```http
POST /login
Content-Type: application/x-www-form-urlencoded

username=admin&password=yourpassword&twoFactorCode=
```

**Response:**

```json
{ "success": true, "msg": "...", "obj": null }
```

The server sets a session cookie (`EXUI`). Include this cookie in all subsequent API requests.

**Step 2 — Use API**

```http
GET /panel/api/inbounds/list
Cookie: EXUI=<session_cookie>
```

---

### Inbound Endpoints

All routes are prefixed with `/panel/api/inbounds`.

#### `GET /list`

List all inbounds for the current user.

```bash
curl -b cookies.txt http://localhost:2053/panel/api/inbounds/list
```

**Response:**
```json
{
  "success": true,
  "obj": [
    {
      "id": 1,
      "userId": 1,
      "up": 0,
      "down": 0,
      "total": 0,
      "remark": "my-inbound",
      "enable": true,
      "expiryTime": 0,
      "clientStats": [],
      "listen": "",
      "port": 443,
      "protocol": "vless",
      "settings": "...",
      "streamSettings": "...",
      "tag": "inbound-443",
      "sniffing": "..."
    }
  ]
}
```

---

#### `GET /get/:id`

Get a single inbound by ID.

```bash
curl -b cookies.txt http://localhost:2053/panel/api/inbounds/get/1
```

---

#### `GET /getClientTraffics/:email`

Get traffic stats for a client by email.

```bash
curl -b cookies.txt http://localhost:2053/panel/api/inbounds/getClientTraffics/user@example.com
```

---

#### `GET /getClientTrafficsById/:id`

Get traffic stats for a client by UUID.

```bash
curl -b cookies.txt http://localhost:2053/panel/api/inbounds/getClientTrafficsById/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

---

#### `POST /add`

Add a new inbound.

```bash
curl -b cookies.txt -X POST http://localhost:2053/panel/api/inbounds/add \
  -H "Content-Type: application/json" \
  -d '{
    "remark": "my-vless",
    "port": 443,
    "protocol": "vless",
    "settings": "{\"clients\":[],\"decryption\":\"none\"}",
    "streamSettings": "{\"network\":\"tcp\"}",
    "sniffing": "{\"enabled\":true,\"destOverride\":[\"http\",\"tls\"]}",
    "enable": true,
    "expiryTime": 0,
    "listen": "",
    "total": 0
  }'
```

---

#### `POST /del/:id`

Delete an inbound by ID.

```bash
curl -b cookies.txt -X POST http://localhost:2053/panel/api/inbounds/del/1
```

---

#### `POST /update/:id`

Update an existing inbound.

```bash
curl -b cookies.txt -X POST http://localhost:2053/panel/api/inbounds/update/1 \
  -H "Content-Type: application/json" \
  -d '{ "remark": "updated-name", "enable": true, ... }'
```

---

#### `POST /addClient`

Add a client to an existing inbound.

```bash
curl -b cookies.txt -X POST http://localhost:2053/panel/api/inbounds/addClient \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "settings": "{\"clients\":[{\"id\":\"NEW-UUID\",\"email\":\"user@example.com\",\"limitIp\":0,\"totalGB\":0,\"expiryTime\":0,\"enable\":true,\"tgId\":\"\",\"subId\":\"\"}]}"
  }'
```

---

#### `POST /:id/delClient/:clientId`

Remove a client from an inbound.

```bash
curl -b cookies.txt -X POST \
  http://localhost:2053/panel/api/inbounds/1/delClient/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

---

#### `POST /updateClient/:clientId`

Update an existing client.

```bash
curl -b cookies.txt -X POST \
  http://localhost:2053/panel/api/inbounds/updateClient/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx \
  -H "Content-Type: application/json" \
  -d '{ "id": 1, "settings": "{ ... updated client JSON ... }" }'
```

---

#### `POST /:id/resetClientTraffic/:email`

Reset traffic counters for a specific client.

```bash
curl -b cookies.txt -X POST \
  http://localhost:2053/panel/api/inbounds/1/resetClientTraffic/user@example.com
```

---

#### `POST /resetAllTraffics`

Reset traffic counters for all inbounds.

```bash
curl -b cookies.txt -X POST http://localhost:2053/panel/api/inbounds/resetAllTraffics
```

---

#### `POST /resetAllClientTraffics/:id`

Reset traffic counters for all clients of a specific inbound.

```bash
curl -b cookies.txt -X POST http://localhost:2053/panel/api/inbounds/resetAllClientTraffics/1
```

---

#### `POST /delDepletedClients/:id`

Delete all clients with exhausted traffic or expired time on an inbound. Use `-1` for all inbounds.

```bash
curl -b cookies.txt -X POST http://localhost:2053/panel/api/inbounds/delDepletedClients/1
```

---

#### `POST /clientIps/:email`

Get the recorded IP addresses for a client.

```bash
curl -b cookies.txt -X POST \
  http://localhost:2053/panel/api/inbounds/clientIps/user@example.com
```

---

#### `POST /clearClientIps/:email`

Clear recorded IPs for a client.

```bash
curl -b cookies.txt -X POST \
  http://localhost:2053/panel/api/inbounds/clearClientIps/user@example.com
```

---

#### `POST /onlines`

Get a list of currently online client emails.

```bash
curl -b cookies.txt -X POST http://localhost:2053/panel/api/inbounds/onlines
```

---

#### `GET /createbackup`

Trigger a manual database backup sent to Telegram admins.

```bash
curl -b cookies.txt http://localhost:2053/panel/api/inbounds/createbackup
```

---

### Server Endpoints

All routes are prefixed with `/server`.

| Method | Path | Description |
|---|---|---|
| `POST` | `/status` | Get server + Xray status (CPU, RAM, uptime, traffic) |
| `POST` | `/getXrayVersion` | List available Xray versions (cached 1 min) |
| `POST` | `/stopXrayService` | Stop Xray process |
| `POST` | `/restartXrayService` | Restart Xray process |
| `POST` | `/installXray/:version` | Install a specific Xray version |
| `POST` | `/updateGeofile/:fileName` | Update a geo data file (`geoip.dat`, `geosite.dat`) |
| `POST` | `/logs/:count` | Get last N log lines |
| `POST` | `/getConfigJson` | Get current Xray running configuration |
| `GET` | `/getDb` | Download the SQLite database file |
| `POST` | `/importDB` | Upload and import a database file |
| `POST` | `/getNewX25519Cert` | Generate a new X25519 key pair |

---

### Response Format

Every API response follows this structure:

```json
{
  "success": true,
  "msg": "Optional human-readable message",
  "obj": { ... }
}
```

| Field | Type | Description |
|---|---|---|
| `success` | `bool` | Whether the operation succeeded |
| `msg` | `string` | Status or error message |
| `obj` | `any` | Returned data (null if not applicable) |

---

## Subscription System

When `subEnable` is `true`, a separate HTTP(S) server starts on `subPort`.

### Endpoints

| Path | Description |
|---|---|
| `{subPath}{subId}` | Standard subscription (base64-encoded links) |
| `{subJsonPath}{subId}` | JSON subscription (clash-compatible config) |

### Subscription Info Header

Responses include:

```
profile-update-interval: 12
content-disposition: attachment; filename="ProfileTitle"
subscription-userinfo: upload=X; download=X; total=X; expire=X
```

---

## Telegram Bot

The Telegram bot provides full panel management from a chat interface.

### Setup

1. Create a bot via [@BotFather](https://t.me/BotFather) — get the token
2. Get your Telegram user ID (e.g. via [@userinfobot](https://t.me/userinfobot))
3. Enter both in **Settings → Telegram Bot**

### Bot Commands

| Command | Description |
|---|---|
| `/start` | Show main menu |
| `/help` | Show available commands |
| Status menu | View server CPU/RAM/network stats |
| Inbounds menu | List, add, delete, update inbounds |
| Client menu | Add/remove/update clients, reset traffic |
| Backup | Send DB backup to admin |
| Usage stats | Per-client traffic reports |

### Notifications

The bot sends automatic notifications for:

- **Login events** — panel web login (success/fail)
- **Scheduled stats** — traffic summary at configured cron time
- **CPU alert** — when usage exceeds `tgCpu` threshold
- **Client expiry / traffic warning** — when approaching limits

---

## Docker Deployment

### Using Docker Compose

```yaml
services:
  eEXUI:
    build:
      context: .
      dockerfile: ./Dockerfile
    container_name: eEXUI_app
    volumes:
      - $PWD/db/:/etc/ex-ui/
      - $PWD/cert/:/root/cert/
    environment:
      XRAY_VMESS_AEAD_FORCED: "false"
      eEXUI_ENABLE_FAIL2BAN: "true"
    network_mode: host
    restart: unless-stopped
```

```bash
docker compose up -d
```

### Volumes

| Host Path | Container Path | Purpose |
|---|---|---|
| `./db/` | `/etc/ex-ui/` | Database & config persistence |
| `./cert/` | `/root/cert/` | TLS certificates |

---

## Database

EXUI uses **SQLite** via GORM with automatic schema migration on startup.

### Location

```
/etc/ex-ui/ex-ui.db
```

Override with env var `XUI_DB_FOLDER`.

### Key Models

#### `Inbound`

| Column | Type | Description |
|---|---|---|
| `id` | int | Primary key |
| `userId` | int | Owner user ID |
| `up` | int64 | Upload bytes |
| `down` | int64 | Download bytes |
| `total` | int64 | Traffic limit (0 = unlimited) |
| `remark` | string | Display name |
| `enable` | bool | Active/inactive |
| `expiryTime` | int64 | Unix timestamp (0 = no expiry) |
| `listen` | string | Bind address |
| `port` | int | Bind port |
| `protocol` | string | Xray protocol |
| `settings` | string | JSON: clients, decryption, etc. |
| `streamSettings` | string | JSON: network, security, etc. |
| `tag` | string | Xray inbound tag |
| `sniffing` | string | JSON: sniffing config |
| `clientStats` | []ClientTraffic | Related client traffic rows |

---

## Project Structure

```
EXUI/
├── config/              # App name, version, env config
├── database/            # DB connection, models, migrations
│   └── model/           # GORM model definitions
├── logger/              # Logging wrapper
├── sub/                 # Subscription server & link generators
├── util/                # Shared utilities
│   ├── common/          # Error helpers, formatting
│   ├── crypto/          # Encryption helpers
│   ├── json_util/       # JSON helpers
│   ├── random/          # Secure random strings
│   ├── reflect_util/    # Reflection helpers
│   └── sys/             # OS-level info (cross-platform)
├── web/                 # Main panel web server
│   ├── assets/          # Static frontend files (embedded)
│   ├── controller/      # HTTP controllers
│   ├── entity/          # Response data structures
│   ├── global/          # Shared server state
│   ├── html/            # Go templates (embedded)
│   ├── job/             # Cron job definitions
│   ├── locale/          # i18n loader & middleware
│   ├── middleware/       # Gin middleware (auth, domain, redirect)
│   ├── network/         # Auto-HTTPS listener/connector
│   ├── service/         # Business logic
│   ├── session/         # Session management
│   └── translation/     # i18n translation files (.toml)
├── xray/                # Xray-core process & gRPC API client
├── main.go              # Application entrypoint
├── go.mod               # Go module definition
├── Dockerfile
├── docker-compose.yml
└── install.sh           # Auto-installer for Linux
```

---

## Development

### Prerequisites

- Go **1.24+**
- A Linux/macOS/Windows environment
- Xray binary in `bin/` folder (or set `XUI_BIN_FOLDER`)

### Run in Debug Mode

```bash
XUI_DEBUG=true go run main.go
```

In debug mode:
- Gin runs in `DebugMode` (logs all requests)
- HTML templates are loaded from disk (hot reload)
- Static assets served from `web/assets/` directly

### Build

```bash
go build -o ex-ui main.go
```

### Cron Jobs

| Job | Interval | Purpose |
|---|---|---|
| `CheckXrayRunningJob` | Every 1s | Ensure Xray is alive |
| Xray restart check | Every 30s | Restart if flagged |
| `XrayTrafficJob` | Every 10s | Pull & store traffic stats |
| `CheckClientIpJob` | Every 10s | Parse logs for client IPs |
| `CheckCpuUsageJob` | — | Alert on CPU spike |
| `ClearLogsJob` | Daily | Purge old log files |
| `StatsNotifyJob` | Configurable | Telegram stats report |
| `CheckHashStorageJob` | Periodic | Clean expired hash cache |

### Tech Stack

| Layer | Library |
|---|---|
| HTTP Router | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) + SQLite |
| Xray Core | [xtls/xray-core](https://github.com/XTLS/Xray-core) |
| Telegram Bot | [mymmrac/telego](https://github.com/mymmrac/telego) |
| Scheduler | [robfig/cron](https://github.com/robfig/cron) |
| Sessions | [gin-contrib/sessions](https://github.com/gin-contrib/sessions) |
| i18n | [nicksnyder/go-i18n](https://github.com/nicksnyder/go-i18n) |
| System Stats | [shirou/gopsutil](https://github.com/shirou/gopsutil) |
| TOTP (2FA) | [xlzd/gotp](https://github.com/xlzd/gotp) |
| HTTP Client | [valyala/fasthttp](https://github.com/valyala/fasthttp) |

---

<div align="center">

**EXUI** · GPL-3.0 License · by [FakeErrorX](https://github.com/FakeErrorX)

</div>
