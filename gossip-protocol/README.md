# Gossip Protocol (Go + Gin)

## Features

- Pluggable gossip initiation strategies: anti-entropy, rumor-mongering, aggregation
- Pluggable spread strategies: push, pull, push-pull
- Configurable fanout, gossip interval, and buffer size
- Peer auto-discovery via /join
- Gin-based HTTP server
- Structured logs for debugging

## 🛠 Prerequisites

- Go 1.21+
- Terminal (for running multiple instances)
- Optional: curl for testing endpoints

## ▶️ How to Run

1. Build

```sh
   go mod tidy
   go build -o gossip-node main.go
```

2. Start Multiple Nodes
   Run 3 nodes in separate terminals:

```bash
# Terminal 1: Start first node (port 8000)

./gossip-node --port=8000

# Terminal 2: Start second node (port 8001, peer with node 1)

./gossip-node --port=8001 --peers=<http://localhost:8000>

# Terminal 3: Start third node (port 8002, peer with node 1 and 2)

./gossip-node --port=8002 --peers=<http://localhost:8000,http://localhost:8001>

```

Nodes will announce themselves and exchange gossip messages periodically.

## 🌐 Available Endpoints

1. Check Node Health

   ```bash
   curl <http://localhost:8000/health>
   ```

2. Add a New Peer (Dynamic Join)

   ```bash
   curl -X POST <http://localhost:8000/join> -H "Content-Type: application/json" \
    -d '{"url": "<http://localhost:8003"}>'
   ```

## 🧩 Logs & Debugging

Each node prints structured logs:

```csharp
[BOOT] Gossip node 45f8 started on port 8000 | Strategy: anti-entropy | Spread: push
[JOIN] Peer added: <http://localhost:8001> | Current peers: [http://localhost:8001]
[GOSSIP] Spreading state to 2 peers: [http://localhost:8001 http://localhost:8002]
[RECV] Gossip received from 45f8 | Health: map[node-1:true]
```

## 🔄 Testing Spread Strategies

Change spread strategy:

```bash
config.SpreadMethod = "push-pull" // or "pull"
```

Change gossip strategy:

```bash
config.Strategy = "rumor-mongering"
```
