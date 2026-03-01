
<div align="center">

<img src="https://raw.githubusercontent.com/FakeErrorX/EXUI/master/media/logo.png" alt="EXUI Logo" width="120" />

# EXUI

**Advanced Xray Management Panel**

*A powerful, open-source web panel for managing Xray-core — clean UI, full REST API, and Telegram automation.*

<br/>

[![Release](https://img.shields.io/github/v/release/fakeerrorx/EXUI.svg?style=for-the-badge&logo=github&color=4A90D9)](https://github.com/FakeErrorX/EXUI/releases)
[![Build](https://img.shields.io/github/actions/workflow/status/fakeerrorx/EXUI/release.yml.svg?style=for-the-badge&logo=githubactions&logoColor=white)](https://github.com/FakeErrorX/EXUI/actions)
[![Go](https://img.shields.io/github/go-mod/go-version/fakeerrorx/EXUI.svg?style=for-the-badge&logo=go&logoColor=white&color=00ADD8)](#)
[![Downloads](https://img.shields.io/github/downloads/fakeerrorx/EXUI/total.svg?style=for-the-badge&logo=github&color=28a745)](https://github.com/FakeErrorX/EXUI/releases/latest)
[![License](https://img.shields.io/badge/license-GPL%20V3-blue.svg?longCache=true&style=for-the-badge&logo=gnu)](https://www.gnu.org/licenses/gpl-3.0.en.html)

<br/>

[**📖 Documentation**](DOCUMENTATION.md) · [**🚀 Quick Start**](#-quick-start) · [**🐳 Docker**](#-docker) · [**🤖 Telegram Bot**](#-telegram-bot) · [**📡 REST API**](#-rest-api)

</div>

---

> [!IMPORTANT]
> This project is intended for **personal use only**. Do not use it for illegal purposes or in production environments without understanding the risks.

---

## ✨ Features

<table>
<tr>
<td>

- 🔌 **Multi-Protocol** — VLESS, VMess, Trojan, Shadowsocks, SOCKS, HTTP and more
- 🖥️ **Live Dashboard** — real-time CPU, RAM, network, and Xray status
- 👥 **Client Management** — add, edit, bulk-delete, traffic reset
- 📊 **Traffic Monitoring** — per-client up/down with expiry tracking
- 🔐 **Two-Factor Auth** — TOTP-based 2FA for panel login

</td>
<td>

- 🤖 **Telegram Bot** — full panel management from chat
- 📬 **Subscription Links** — standard & JSON/Clash-compatible
- 🌐 **REST API** — complete programmatic control
- 🔒 **Auto TLS** — built-in HTTPS with custom cert support
- 🌍 **i18n** — multi-language interface support

</td>
</tr>
</table>

---

## 🚀 Quick Start

### One-line Install (Linux)

```bash
bash <(curl -Ls https://raw.githubusercontent.com/fakeerrorx/EXUI/master/install.sh)
```

> Supports: Ubuntu · Debian · CentOS · Fedora · Arch · Alpine and more  
> Requires: GLIBC ≥ 2.32 · Root access

After installation, access the panel at:

```
http://YOUR_IP:2053
```

### Management Commands

```bash
ex-ui start         # Start the service
ex-ui stop          # Stop the service
ex-ui restart       # Restart the service
ex-ui status        # Show current status
ex-ui update        # Update to latest version
ex-ui log           # View live logs
ex-ui uninstall     # Remove EXUI
```

---

## 🐳 Docker

```yaml
services:
  eEXUI:
    image: ghcr.io/fakeerrorx/exui:latest
    container_name: eEXUI_app
    volumes:
      - ./db/:/etc/ex-ui/
      - ./cert/:/root/cert/
    environment:
      XRAY_VMESS_AEAD_FORCED: "false"
    network_mode: host
    restart: unless-stopped
```

```bash
docker compose up -d
```

| Volume | Purpose |
|--------|---------|
| `./db/` → `/etc/ex-ui/` | Database & config persistence |
| `./cert/` → `/root/cert/` | TLS certificates |

---

## 📡 REST API

EXUI exposes a full REST API under `/panel/api/inbounds/`. All endpoints require session authentication.

**Login first:**

```bash
curl -c cookies.txt -X POST http://localhost:2053/login \
  -d "username=admin&password=yourpassword"
```

**Then use the API:**

```bash
# List all inbounds
curl -b cookies.txt http://localhost:2053/panel/api/inbounds/list

# Add a client
curl -b cookies.txt -X POST http://localhost:2053/panel/api/inbounds/addClient \
  -H "Content-Type: application/json" \
  -d '{ "id": 1, "settings": "{ ... }" }'

# Get client traffic by email
curl -b cookies.txt http://localhost:2053/panel/api/inbounds/getClientTraffics/user@example.com
```

> 📖 For the complete API reference with all endpoints, request/response formats, and examples, see the [**Documentation**](DOCUMENTATION.md#api-reference).

---

## 🤖 Telegram Bot

Enable the bot in **Settings → Telegram Bot** to get:

- 📋 Inbound & client management from chat
- 📈 Scheduled traffic reports
- 🔔 Login notifications
- ⚠️ CPU spike & traffic limit alerts
- 💾 Automated database backups

**Setup**:
1. Create a bot via [@BotFather](https://t.me/BotFather)
2. Get your user ID via [@userinfobot](https://t.me/userinfobot)
3. Paste both into the panel settings

---

## ⚙️ Configuration

Key environment variables:

| Variable | Default | Description |
|---|---|---|
| `XUI_LOG_LEVEL` | `info` | Log verbosity |
| `XUI_DEBUG` | `false` | Enable debug mode |
| `XUI_BIN_FOLDER` | `bin` | Xray binary location |
| `XUI_DB_FOLDER` | `/etc/ex-ui` | Database path |

> 📖 All panel settings, subscription options, and Telegram configuration are documented in [**DOCUMENTATION.md**](DOCUMENTATION.md).

---

## 🏗️ Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.24+ |
| Framework | Gin · GORM |
| Database | SQLite |
| Proxy Core | Xray-core |
| Scheduler | robfig/cron |
| 2FA | TOTP (gotp) |
| Telegram | telego |
| Frontend | Vue 2 · Ant Design Vue |

---

## 🤝 Contributing

Contributions are welcome! Please open an issue before submitting large pull requests.

```bash
# Clone and run in debug mode
git clone https://github.com/FakeErrorX/EXUI.git
cd EXUI
XUI_DEBUG=true go run main.go
```

---

## ⭐ Support

**If EXUI saves you time, consider giving it a star!** ⭐

<div align="center">

[![Star History](https://img.shields.io/github/stars/fakeerrorx/EXUI?style=for-the-badge&logo=github&color=yellow)](https://github.com/FakeErrorX/EXUI/stargazers)

---

GPL-3.0 License · © [FakeErrorX](https://github.com/FakeErrorX)

</div>
