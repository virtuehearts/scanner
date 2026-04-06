# Fortune Bitcoin Scanner 💰🚀

Fortune is an elite Bitcoin wallet cracker specializing in P2PKH private keys. It combines high-performance brute forcing with intelligent "brain forcing" across a dataset of over **323,000** non-zero balance wallets.

Designed for the modern era, Fortune is **AI Agent Native**, allowing bots like **Claude Code**, **Hermes**, or other autonomous agents to control scanning operations via a robust FastAPI bridge and containerized environment.

![btc cracker telegram screenshot](/docs/screenshot.webp?raw=true)

## 🌟 Key Features

- **High-Performance Core:** Written in Go for maximum efficiency in brute-forcing private keys.
- **GPU Acceleration:** Optional OpenCL-based GPU acceleration for NVIDIA RTX and other compatible hardware.
- **Bloom Filter Support:** Ultra-fast address lookup using memory-efficient Bloom filters.
- **AI Agent Integration:** Native support for AI agents (Claude, Hermes, etc.) to start, stop, and monitor scans autonomously.
- **FastAPI Bridge:** A modern REST API layer that bridges the Go scanner with web and AI interfaces.
- **Containerized Architecture:** Fully Dockerized for easy deployment in any environment.
- **Real-time Monitoring:** Modern Web UI dashboard for tracking IOPS, tried keys, and found wallets.
- **Multi-Strategy Scanning:**
    - `bruteforce`: Traditional random key generation against rich address datasets.
    - `brainforce`: Dictionary and permutation-based attacks on brain wallets.
- **Instant Notifications:** Integrated Telegram bot support for real-time alerts on found keys.
- **Rich Datasets:** Includes pre-packaged datasets of wallets with non-zero balances.

---

## 🚀 Quick Start with Docker

The fastest way to deploy the full stack (Scanner + API Bridge + Web UI) is using the bridge Dockerfile:

```bash
# Build the image
docker build -t fortune-scanner -f Dockerfile.bridge .

# Run the container
docker run -d -p 8000:8000 --name fortune-app fortune-scanner
```

- **Web UI:** [http://localhost:8000/static/index.html](http://localhost:8000/static/index.html)
- **API Docs:** [http://localhost:8000/docs](http://localhost:8000/docs)

---

## 🤖 AI Agent & API Integration

Fortune is built to be operated by AI bots. Agents can manage the entire lifecycle of a scan through the API.

### Authentication
All API requests require the `x-api-key` header.
- **Default Key:** `your-secret-api-key` (Configure via `BRIDGE_API_KEY` env var)

### API Skills for Agents

| Skill | Endpoint | Description |
|---|---|---|
| **Start Scan** | `POST /api/start` | Launch a `bruteforce` or `brainforce` task with custom workers. |
| **Stop Scan** | `POST /api/stop` | Terminate the active scanning process. |
| **Status** | `GET /api/status` | Get real-time IOPS, total keys tried, and found keys. |
| **Logs** | `GET /api/logs` | Retrieve the latest 100 lines of scanner output. |
| **Found** | `GET /api/found` | List all discovered private keys and their addresses. |
| **Files** | `GET /api/files` | List available address datasets for scanning. |

Detailed integration guides can be found in [agent_skills.MD](agent_skills.MD).

---

## 💻 CLI Usage (Traditional)

For users who prefer the command line, the Go binary can be used directly:

```bash
# Run brute force with 4 workers
./fortune --workers 4 bruteforce

# Run brain wallet permutation attack
./fortune --workers 2 brainforce --pass-length 6 --pass-alphabet english-lower
```

### Global Options
- `--workers`: Number of parallel execution threads.
- `--gpu`: Enable GPU acceleration using OpenCL.
- `--bloom`: Enable Bloom filter for fast address lookup.
- `--batch-size`: Size of key generation batches (default 1024, recommended for GPU).
- `--heartbit-sec`: Frequency of status updates to STDOUT.
- `--telegram-token`: Token for Telegram notifications.
- `--telegram-channel`: Channel name for findings alerts.

### High-Performance Usage (NVIDIA RTX)
To reach millions of operations per second, use the GPU and Bloom filter:
```bash
./fortune --gpu --bloom --batch-size 4096 --workers 8 bruteforce
```

---

## 🛠 Configuration

| Environment Variable | Default | Description |
|---|---|---|
| `BRIDGE_API_KEY` | `your-secret-api-key` | Security key for API access. |
| `PORT` | `8000` | Port for the FastAPI bridge. |

## ⚖️ License & Disclaimer

This tool is for educational and authorized security testing purposes only. The authors are not responsible for any misuse or damage caused by this software.

© [github.com/shlima/fortune](https://github.com/shlima/fortune)
