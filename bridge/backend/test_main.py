import pytest
from fastapi.testclient import TestClient
from bridge.backend.main import app
import os

client = TestClient(app)
API_KEY = "your-secret-api-key"

def test_read_root():
    response = client.get("/")
    assert response.status_code == 200
    assert "Fortune Scanner API is running" in response.json()["message"]

def test_unauthorized():
    response = client.get("/api/status")
    assert response.status_code == 403

def test_get_status():
    response = client.get("/api/status", headers={"x-api-key": API_KEY})
    assert response.status_code == 200
    assert "running" in response.json()

def test_list_files():
    # Ensure addresses directory exists for the test
    os.makedirs("addresses", exist_ok=True)
    with open("addresses/test_addresses.txt", "w") as f:
        f.write("1LQoWist8KkaUXSPKZHNvEyfrEkPHzSsCd\n")

    response = client.get("/api/files", headers={"x-api-key": API_KEY})
    assert response.status_code == 200
    assert any("test_addresses.txt" in f for f in response.json()["files"])
