# Fortune Scanner AI Agent Skills

This document describes how an AI agent (e.g., Hermes) can interact with the Fortune Scanner Bridge API to manage Bitcoin address cracking tasks.

## API Configuration

- **Base URL:** `http://<host>:8000`
- **Auth:** Requires `x-api-key` header for all requests.

## Skills

### 1. Start Scanner
Launch a new scanning process with specific parameters.

- **Endpoint:** `POST /api/start`
- **Payload:**
```json
{
  "command": "bruteforce", // or "brainforce"
  "workers": 4,
  "gpu": true,
  "bloom": true,
  "batch_size": 4096,
  "files": ["addresses/Bitcoin/2023/04/p2pkh_Rich_Max_1.txt"],
  "night": false,
  "telegram_token": "optional_token",
  "telegram_channel": "@optional_channel"
}
```

### 2. Stop Scanner
Terminate the currently running scanning process.

- **Endpoint:** `POST /api/stop`

### 3. Get Status
Retrieve real-time metrics and running state.

- **Endpoint:** `GET /api/status`
- **Response:**
```json
{
  "running": true,
  "iops": 15420,
  "total_tried": 120500,
  "found_keys": [],
  "last_heartbeat": "tried: 120500, iops: 15420"
}
```

### 4. List Address Files
Find available address datasets to scan.

- **Endpoint:** `GET /api/files`

### 5. View Logs
Read the last 100 lines of the scanner output.

- **Endpoint:** `GET /api/logs`

### 6. Get Found Wallets
Retrieve any successfully found private keys.

- **Endpoint:** `GET /api/found`

## Example Hermes Integration

The agent can monitor the scanner and report findings:

1. **Check status:** `GET /api/status`
2. **If `found_keys` is not empty:** Parse the findings and alert the user.
3. **If `iops` is low:** Suggest increasing workers or checking system resources.
