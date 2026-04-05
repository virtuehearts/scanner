# Fortune Scanner Bridge 💰

Elite Bitcoin wallet cracker (P2PKH private keys) with a modern Web UI and AI Agent compatible FastAPI bridge.

## Features
- **FastAPI Backend:** A high-performance bridge for controlling the scanner.
- **Modern Web UI:** Real-time dashboard for monitoring IOPS, tried keys, and found wallets.
- **AI Agent Native:** Built-in skills for agents like Hermes to operate the scanner autonomously.
- **Rich Wallet Datasets:** Integrated datasets of over 300,000 wallets with non-zero balances.

## Deployment with Docker

The fastest way to get started is by using the bridge Dockerfile:

```bash
docker build -t fortune-bridge -f Dockerfile.bridge .
docker run -d -p 8000:8000 --name fortune-scanner fortune-bridge
```

Access the UI at: `http://localhost:8000/static/index.html`

## AI Agent Integration

The bridge provides a standard API for AI agents. See [agent_skills.MD](agent_skills.MD) for the complete API reference and integration guide.

## Configuration

When using the Web UI or API, you must provide the `x-api-key` header (default key is `your-secret-api-key` unless overridden via environment variable `BRIDGE_API_KEY`).

### Environment Variables
| Variable | Description |
|---|---|
| BRIDGE_API_KEY | Secret key for API authentication |
| PORT | Port for the bridge server (default 8000) |

## Traditional CLI Usage

The original Go command line tool is still available:

```bash
./fortune bruteforce --workers 4
```

For more CLI details, see the original [README.md](README.md).
