"""Shared pytest fixtures for E2E tests.

Place in tests/e2e/conftest.py
"""

import pytest
import os


@pytest.fixture(scope="session")
def base_url():
    return os.environ.get("API_URL", "http://localhost:8080")


@pytest.fixture(scope="session")
def api_client(base_url):
    """Provides a requests.Session with base URL configured."""
    import requests
    session = requests.Session()
    session.base_url = base_url
    yield session
    session.close()


@pytest.fixture(scope="session")
def auth_token(base_url):
    """Acquire and cache an auth token for the session."""
    import requests
    resp = requests.post(
        f"{base_url}/api/auth/login",
        json={
            "username": os.environ.get("TEST_USERNAME", "admin"),
            "password": os.environ.get("TEST_PASSWORD", "password"),
        }
    )
    assert resp.status_code == 200, f"auth failed: {resp.text}"
    return resp.json()["token"]
