# 📦 Jiboia Tunnel — Structure and Usage

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

🌍 This documentation is also available in:
- [Português (Brasil)](README.pt-BR.md)
- [Español](README.es.md)

Jiboia Tunnel is a reverse HTTP tunneling tool based on WebSocket, inspired by solutions like `ngrok` and `inlets`. Below is the detailed documentation for developers who want to understand and replicate the project with precision.

---

## 📁 Project Structure (separated binaries)

```
jiboia-tunnel/
├── cmd/
│   ├── jiboia-client/main.go     # Client: connects to the relay and sends local traffic
│   ├── jiboia-relay/main.go      # Relay: accepts WebSocket and forwards requests
│   ├── jiboia-server/main.go     # Mock: local test server
├── shared/
│   └── message.go                # Common types: TunnelMessage, TunnelResponse
├── go.mod
├── go.sum
└── README.md
```

Each `cmd/<name>` directory defines a separate binary:
- `jiboia-client`
- `jiboia-relay`
- `jiboia-server`

---

## 🚀 Available Commands (post-build)

### Start the relay
```bash
./jiboia-relay
```
- WebSocket listening on `/ws`
- HTTP requests are routed to connected clients

---

### Start the client (local user)
```bash
./jiboia-client http 3000
```
Shortcut to expose `localhost:3000` via default relay (`ws://localhost:80/ws`).

Your local service (`localhost:3000`) becomes:
```
http://<generated-name>.jiboia.local
```

#### With defined name and remote relay:
```bash
./jiboia-client http 3000 --name myapp --relay wss://relay.jiboia.io/ws
```

#### Additional available flags:
| Flag             | Type     | Description                                                  |
|------------------|----------|--------------------------------------------------------------|
| `--name`         | string   | Tunnel name (subdomain).                                    |
| `--relay`        | string   | Relay WebSocket address.                                    |
| `--proto`        | string   | Protocol to expose (`http`, `tcp`).                         |
| `--hostname`     | string   | Custom full domain (e.g., `mydomain.com`).                  |
| `--inspect`      | bool     | Enable detailed traffic view (debug mode).                  |
| `--authtoken`    | string   | Authentication token with the server.                       |
| `--config`       | string   | Path to external config file.                               |
| `--region`       | string   | Relay region (e.g., `us`, `sa-east`).                       |
| `--label`        | string   | Friendly tunnel label (used in logs/API).                   |
| `--log-level`    | string   | Logging level (`debug`, `info`, `warn`, `error`).           |

---

### Start local test server
```bash
./jiboia-server
```
Simple HTML server on `http://localhost:3000`

---

## 🧪 Local Testing with `go run`
```bash
# Terminal 1
sudo go run ./cmd/jiboia-relay/main.go

# Terminal 2
go run ./cmd/jiboia-server/main.go

# Terminal 3
go run ./cmd/jiboia-client/main.go --name jiboia --local http://localhost:3000
```

Open in browser:
```
http://jiboia.jiboia.local
```

Add to your `/etc/hosts`:
```
127.0.0.1 jiboia.jiboia.local
```

---

## 🛠 Building binaries
```bash
# Manual build
GOOS=linux GOARCH=amd64 go build -o jiboia-relay ./cmd/jiboia-relay
GOOS=linux GOARCH=amd64 go build -o jiboia-client ./cmd/jiboia-client
GOOS=linux GOARCH=amd64 go build -o jiboia-server ./cmd/jiboia-server
```
Or use `goreleaser` to define multiple builds per binary.

---

## 📦 Terminal Installation (Linux/macOS/Windows)

### curl (Linux/macOS)
```bash
curl -s https://raw.githubusercontent.com/valcinei/jiboia-tunnel/main/install.sh | bash
```

### PowerShell (Windows)
```powershell
iwr https://raw.githubusercontent.com/valcinei/jiboia-tunnel/main/install.ps1 -useb | iex
```

The script detects your platform, downloads the latest binaries, and places them in:
- Linux/macOS: `/usr/local/bin`
- Windows: `%ProgramFiles%\JiboiaTunnel\`

---

## 🧱 How the application works
- **relay:** receives HTTP requests, extracts subdomain, redirects via WebSocket to the connected client.
- **client:** listens to WebSocket messages and acts as a reverse proxy to the local server.
- **server:** mock app to test the tunneling chain.

---

## ✅ Remaining steps for persistence and authentication

### 🔐 JWT Authentication
1. Create middleware `RequireAuth()` to protect routes (`/tunnels`, etc).
2. Apply middleware to REST routes in `jiboia-server`.
3. Add JWT token validation via cookie or `Authorization: Bearer`.
4. Create optional `logout` endpoint.
5. (future) Add `/users` route with user persistence.

### 💾 Real SQLite persistence
1. Create `Migrate()` function to run `CREATE TABLE IF NOT EXISTS tunnels (...)`.
2. Create `SQLiteStore` type implementing the `TunnelStore` interface.
3. Replace `InMemoryStore` with `SQLiteStore`.
4. Add error handling for opening database (path, permissions).

### 🔑 Client token support
1. Add `--authtoken` support to `jiboia-client`.
2. Include token in the `Authorization` header for API calls.
3. Validate token in `relay` to allow/deny tunnel connection.

---

## 🌐 Support for custom user domains

### Goal
Allow authenticated users to register and use custom domains for their tunnels.

### Implementation steps

7. **Domain registration by the user:**
   - Create `POST /domains` endpoint in `jiboia-server`.
   - Associate the domain with the tunnel and authenticated user.
   - Example payload:
     ```json
     { "hostname": "mydomain.com" }
     ```
   - Optionally require DNS TXT verification in the future.

8. **Relay validation:**
   - Accept full domain requests in `Host:` header, not just subdomains.
   - Check if domain exists in database and is linked to an active tunnel.

9. **Client flags:**
   - `--hostname` to support external domains.

10. **Save in backend:**
    - Link domain to tunnel in the database.

11. **DNS setup:**
    - User should point domain to relay IP (A or CNAME record).

12. **HTTPS/TLS (future):**
    - Support via Let's Encrypt or Nginx/Caddy config.

13. **DNS validation (future):**
    - `/verify-domain` endpoint for token-based DNS validation.

---

## 🤝 Contributing

Contributions are welcome! You can:
- Open issues with ideas, bugs, or improvements
- Submit pull requests
- Join discussions

### Getting started
1. Fork the repository
2. Clone your fork locally
3. Create a branch:
   ```bash
   git checkout -b my-feature
   ```
4. Make changes and commit:
   ```bash
   git commit -m "feat: add custom hostname support"
   ```
5. Push to your fork:
   ```bash
   git push origin my-feature
   ```
6. Open a pull request to the main repository

See `CONTRIBUTING.md` for more details.

---

## 🔮 Future improvements
- JWT authentication
- HTTPS via Let's Encrypt
- Web dashboard
- REST API in `jiboia-server` for administration
- Load balancing between relays

---

This binary separation improves modularity, allows segmented deployment (e.g., relay in the cloud and client locally), and is ready for production use or future expansion.
