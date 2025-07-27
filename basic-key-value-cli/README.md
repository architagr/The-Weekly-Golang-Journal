# Basic Key-Value Store CLI in Go

This project is a minimal in-memory key-value store built with Go.
It supports:

- PUT: Store a key-value pair with a TTL (time-to-live).
- GET: Retrieve a value by key.
- DELETE: Remove a key.
- Automatic expiration when TTL expires.
- Simple CLI for interactive commands.

## 📂 Project Structure

```csharp
basic-key-value-cli/
│── cmd/
│   └── cli.go        # CLI entry point
│
└── internal/
    └── controller/
        └── store.go  # Store implementation

```

## ▶️ How to Run

1. Clone the repository:

```bash
git clone https://github.com/architagr/The-Weekly-Golang-Journal.git
cd The-Weekly-Golang-Journal/basic-key-value-cli
```

2. Build and run:

```bash
go run cmd/cli.go
```

## 💻 Example Usage

```vbnet
Simple KV Store CLI
Commands: PUT <key> <value> <ttl_sec> | GET <key> | DELETE <key> | EXIT

> PUT name Archit 10
Key stored

> GET name
Value: Archit

(wait 10 seconds…)

> GET name
Error: key expired

> DELETE name
Error: key not found
```

## 🛠 Future Enhancements

- Background goroutine to clean expired keys automatically.
- Persistent storage (e.g., file or database).
- Support for batch operations.
