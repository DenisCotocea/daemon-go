# DAEMON GO: Tailscale Device Monitor & Telegram Alerter

A lightweight and efficient Go service that monitors the status of devices within a Tailscale network (Tailnet) in real-time and sends instant Telegram alerts whenever a node transitions ONLINE or goes OFFLINE.

---

## 🚀 Features

- Native OAuth2 Authentication: Secure integration with Tailscale API v2 using Go's golang.org/x/oauth2/clientcredentials package.
- Zero-Memory Leak Streaming: Direct JSON decoding from HTTP response streams (json.NewDecoder) without intermediate buffer allocations.
- In-Memory State Diffing: Tracks device statuses in memory and fires Telegram notifications strictly on state changes (connectedToControl).
- Clean Alert Formatting: Clear, structured Telegram messages displaying HostName and associated User accounts.

---

## 🛠️ Built With

- Go (Golang) — Core runtime and standard library
- Tailscale API v2 — Device state polling
- Telegram Bot API — Real-time notification dispatch
- OAuth2 Client Credentials — Machine-to-Machine (M2M) authentication

---

## 📁 Project Structure


main.go      # Application entry point, state initialization, and polling ticker

functions.go    # OAuth2 HTTP client, API fetching logic, and Telegram dispatcher

structs.go     # JSON payload structures for Tailscale and Telegram APIs

go.mod       # Go module declaration and dependency lockfile

.env         # Environment variable configuration (git-ignored)

---

## ⚙️ Setup & Installation

### 1. Clone the repository
git clone https://github.com/DenisCotocea/daemon-go.git
cd tailscale-telegram-monitor

### 2. Configure Environment Variables
Create a .env file in the project root containing your API credentials:

CLIENT_ID=your_tailscale_oauth_client_id

CLIENT_SECRET=your_tailscale_oauth_client_secret

TAILNET_ID=your_tailnet_name (e.g., example.com or user@gmail.com)

TELEGRAM_BOT_SECRET=your_telegram_bot_token

TELEGRAM_CHAT_ID=your_telegram_chat_id

---

## 🏃 Running the Service

Run directly during development:
go run .

Compile a production binary:
go build -o daemon-go .
./daemon-go

---

## 📲 Telegram Alert Examples

When a device changes connection state, the bot sends a formatted notification directly to your chat:

🟢 TAILSCALE ALERT

Host: Hostname

Status: ONLINE

User: dev@example.com


🔴 TAILSCALE ALERT

Host: Hostname

Status: OFFLINE

User: dev@example.com