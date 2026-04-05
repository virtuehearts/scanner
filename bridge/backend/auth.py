import os
from fastapi import Header, HTTPException, Depends
from typing import Optional

API_KEY = os.getenv("BRIDGE_API_KEY", "your-secret-api-key")

async def verify_api_key(x_api_key: Optional[str] = Header(None)):
    if not x_api_key or x_api_key != API_KEY:
        raise HTTPException(status_code=403, detail="Could not validate credentials")
    return x_api_key
