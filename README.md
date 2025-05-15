# 📦 Jiboia Tunnel

Jiboia Tunnel is a reverse HTTP tunneling tool based on WebSocket, inspired by solutions like `ngrok` and `inlets`.

## 📁 Project Structure

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

## 🚀 Available Commands (post-build)

### Start the relay
```bash
./jiboia-relay
```
- WebSocket listening on `/ws`
- HTTP requests are routed to connected clients

### Start the client (local user)
```bash
./jiboia-client http 3000 --name myapp
```
Exposes your local service (`localhost:3000`) as:
```
http://myapp.jiboia.local
```

#### With a defined name and remote relay:
```bash
./jiboia-client http 3000 --name myapp --relay wss://relay.jiboia.io/ws
```

#### Additional available flags:
| Flag             | Type     | Description                                                       |
|------------------|----------|-------------------------------------------------------------------|
| `--name`         | string   | Tunnel name (subdomain).                                          |
| `--relay`        | string   | WebSocket address of the relay.                                   |
| `--proto`        | string   | Protocol to expose (`http`, `tcp`).                               |
| `--hostname`     | string   | Full custom domain (e.g., `mywebsite.com`).                       |
| `--inspect`      | bool     | Shows detailed traffic (debug mode).                              |
| `--authtoken`    | string   | Authentication token with the server.                             |
| `--config`       | string   | Path to external configuration file.                              |
| `--region`       | string   | Relay region (e.g., `us`, `sa-east`).                             |
| `--label`        | string   | Friendly tunnel identifier (used in logs/future API).             |
| `--log-level`    | string   | Log level (`debug`, `info`, `warn`, `error`).                     |

### Start local test server
```bash
./jiboia-server
```
Responds with simple HTML at `http://localhost:3000`

## 🧪 Testing Locally with `go run`
```bash
# Terminal 1
sudo go run ./cmd/jiboia-relay/main.go

# Terminal 2
go run ./cmd/jiboia-server/main.go

# Terminal 3
go run ./cmd/jiboia-client/main.go http 3000 --name jiboia
```

Open in the browser:
```
http://jiboia.jiboia.local
```

Add to your `/etc/hosts`:
```
127.0.0.1 jiboia.jiboia.local
```

## 🛠 Building the Binaries
```bash
# Build all manually
GOOS=linux GOARCH=amd64 go build -o jiboia-relay ./cmd/jiboia-relay
GOOS=linux GOARCH=amd64 go build -o jiboia-client ./cmd/jiboia-client
GOOS=linux GOARCH=amd64 go build -o jiboia-server ./cmd/jiboia-server
```
Or with `goreleaser`, defining multiple builds per binary.

## 🧱 How the Application Works
- **relay**: receives HTTP requests, extracts subdomain, redirects via WebSocket to a connected client.
- **client**: listens to WebSocket messages and acts as a reverse proxy for a local server.
- **server**: mock application for testing the tunneling chain.

## 🔮 Future Expansions
- JWT token authentication
- HTTPS support with Let's Encrypt / Caddy
- Web dashboard with tunnel panel
- REST API in `jiboia-server` for administrative control
- Load balancing between multiple relays

This separation by binaries improves control, facilitates segmented deployment (e.g., relay in the cloud, client on a local machine), and follows good modularity practices. Ready for production use or extension with new features. 