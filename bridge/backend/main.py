from fastapi import FastAPI, UploadFile, File, Depends, HTTPException
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel
from typing import List, Optional
import os
import shutil
from bridge.backend.process_manager import ProcessManager
from bridge.backend.auth import verify_api_key

app = FastAPI(title="Fortune Scanner Bridge")
pm = ProcessManager(executable_path="./fortune")

# Models
class ScannerConfig(BaseModel):
    command: str  # e.g., "bruteforce", "brainforce", "wordlist"
    workers: Optional[int] = 1
    files: Optional[List[str]] = None
    telegram_token: Optional[str] = None
    telegram_channel: Optional[str] = None
    night: Optional[bool] = False
    gpu: Optional[bool] = False
    bloom: Optional[bool] = False
    batch_size: Optional[int] = 1024
    pass_length: Optional[int] = None
    pass_alphabet: Optional[List[str]] = None
    pass_shuffle: Optional[int] = None
    pass_file: Optional[str] = None
    args: Optional[List[str]] = []

# Serve static files for the UI
app.mount("/static", StaticFiles(directory="bridge/frontend"), name="static")

@app.get("/")
async def root():
    return {"message": "Fortune Scanner API is running. Access UI at /static/index.html"}

@app.post("/api/start")
async def start_scanner(config: ScannerConfig, api_key: str = Depends(verify_api_key)):
    # Global flags must come BEFORE the command in urfave/cli
    args = []

    if config.workers:
        args += ["--workers", str(config.workers)]

    # We always set heartbit-sec to 1 for the UI
    args += ["--heartbit-sec", "1"]

    if config.files:
        for f in config.files:
            args += ["--file", f]

    if config.telegram_token:
        args += ["--telegram-token", config.telegram_token]

    if config.telegram_channel:
        args += ["--telegram-channel", config.telegram_channel]

    if config.night:
        args += ["--night"]

    if config.gpu:
        args += ["--gpu"]

    if config.bloom:
        args += ["--bloom"]

    if config.batch_size:
        args += ["--batch-size", str(config.batch_size)]

    # Now add the command
    args += [config.command]

    if config.command == "brainforce":
        if config.pass_length:
            args += ["--pass-length", str(config.pass_length)]
        if config.pass_alphabet:
            for a in config.pass_alphabet:
                args += ["--pass-alphabet", a]
        if config.pass_shuffle is not None:
            args += ["--pass-shuffle", str(config.pass_shuffle)]
    elif config.command == "wordlist":
        if config.pass_file:
            args += ["--pass-file", config.pass_file]

    if config.args:
        args += config.args

    success, msg = pm.start(args)
    if not success:
        raise HTTPException(status_code=400, detail=msg)
    return {"message": msg}

@app.post("/api/stop")
async def stop_scanner(api_key: str = Depends(verify_api_key)):
    success, msg = pm.stop()
    if not success:
        raise HTTPException(status_code=400, detail=msg)
    return {"message": msg}

@app.get("/api/status")
async def get_status(api_key: str = Depends(verify_api_key)):
    return pm.get_status()

@app.get("/api/files")
async def list_files(api_key: str = Depends(verify_api_key)):
    # Explore the addresses directory to list available files
    files_list = []
    base_dir = "addresses"
    for root, dirs, files in os.walk(base_dir):
        for f in files:
            if f.endswith(".txt"):
                files_list.append(os.path.join(root, f))
    return {"files": files_list}

@app.post("/api/upload")
async def upload_file(file: UploadFile = File(...), api_key: str = Depends(verify_api_key)):
    # Upload custom address files
    upload_dir = "addresses/uploads"
    os.makedirs(upload_dir, exist_ok=True)

    # Sanitize filename to prevent path traversal
    filename = os.path.basename(file.filename)
    if not filename:
        raise HTTPException(status_code=400, detail="Invalid filename")

    file_path = os.path.join(upload_dir, filename)

    with open(file_path, "wb") as buffer:
        shutil.copyfileobj(file.file, buffer)

    return {"filename": filename, "path": file_path}

@app.get("/api/logs")
async def get_logs(api_key: str = Depends(verify_api_key)):
    if not os.path.exists("scanner.log"):
        return {"logs": []}

    # Efficiently read the last 100 lines without loading the whole file
    lines = []
    chunk_size = 1024
    with open("scanner.log", "rb") as f:
        f.seek(0, os.SEEK_END)
        file_size = f.tell()

        buffer = b""
        pointer = file_size

        while len(lines) <= 100 and pointer > 0:
            step = min(pointer, chunk_size)
            pointer -= step
            f.seek(pointer)
            buffer = f.read(step) + buffer
            lines = buffer.decode('utf-8', errors='ignore').splitlines()

    return {"logs": lines[-100:]}

@app.get("/api/found")
async def get_found_keys(api_key: str = Depends(verify_api_key)):
    import json
    found = []
    if os.path.exists("found_keys.json"):
        with open("found_keys.json", "r") as f:
            for line in f:
                if line.strip():
                    found.append(json.loads(line))
    return {"found": found}
