# Pytest Generate Strategy

Profile-specific test generation rules for the `gen-test-scripts` skill.

## Test Runner & Imports

| Item | Value |
|------|-------|
| Test runner | pytest (`python -m pytest`) |
| Assertions | Built-in `assert` statements |
| Fixtures | `@pytest.fixture` from `conftest.py` |
| Markers | `@pytest.mark.e2e`, `@pytest.mark.api`, `@pytest.mark.cli` |

## Spec Template Mapping

| Test type | Template file | Output filename |
|-----------|--------------|-----------------|
| Mixed | `templates/test_file.py` | `test_<feature>.py` |

Pytest discovers files matching `test_*.py` or `*_test.py`.

## CLI Testing

Use `subprocess.run`:

```python
result = subprocess.run(
    ["binary", "--flag", "value"],
    capture_output=True, text=True
)
assert result.returncode == 0
assert "expected text" in result.stdout
```

## API Testing

Use `requests` or `httpx`:

```python
import requests

resp = requests.get("http://localhost:8080/api/health")
assert resp.status_code == 200
assert resp.json()["status"] == "ok"
```

## Auth

- API: `requests` session with auth headers or cookies
- CLI: environment variables via `os.environ`
- Fixtures in `conftest.py` for auth setup/teardown

## Import Conventions

```python
import pytest
import subprocess
import os
from helpers import run_cli, api_client
```

## Anti-Patterns (Forbidden)

- No `time.sleep()` — use `pytest`'s retry mechanisms or poll loops with timeout
- No hardcoded URLs — use environment variables or config from `conftest.py`
- No bare `except:` — always catch specific exceptions

## Compilation Check

```bash
just e2e-compile
```

## Traceability

Each test function includes a docstring with TC ID and PRD source:

```python
def test_tc_nnn_description():
    """Traceability: TC-NNN → {PRD Source}"""
    ...
```

## Conftest

Shared fixtures go in `tests/e2e/conftest.py`:

```python
import pytest

@pytest.fixture(scope="session")
def api_client():
    base_url = os.environ.get("API_URL", "http://localhost:8080")
    # ... setup
    yield client
    # ... teardown
```
