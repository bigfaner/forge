---
title: "Python pytest Testing Convention"
---

# Python pytest Testing Convention

Convention for generating Python test code using the pytest framework with its rich assertion introspection.

## framework

- **name**: pytest
- **version**: pytest 7.0+
- **language**: Python
- **runner_command**: `pytest -v`

## discovery

- **test_dir**: `tests/`
- **file_pattern**: `test_*.py`, `*_test.py`
- **exclude_pattern**: `__pycache__/`, `.tox/`, `node_modules/`

## structure

- **suite_pattern**: Files act as suites — each `test_*.py` file is a collection of test functions
- **case_pattern**: `def test_<description>()` — top-level test functions with `test_` prefix
- **hook_pattern**: `conftest.py` with `pytest.fixture`, `session`/`module`/`class` scoped fixtures

### Test Function Naming

Pattern: `test_<description>` using snake_case with descriptive names.

```python
def test_login_with_valid_credentials():
    ...
```

### Class-Based Tests

Group related tests using classes (no inheritance required):

```python
class TestTaskLifecycle:
    def test_claim_task(self):
        ...

    def test_complete_task(self):
        ...
```

### Parametrized Tests

Use `@pytest.mark.parametrize` for data-driven testing:

```python
@pytest.mark.parametrize("input,expected", [
    ("hello", "HELLO"),
    ("", ""),
    ("123", "123"),
])
def test_uppercase(input, expected):
    assert input.upper() == expected
```

### Fixtures

Use `conftest.py` for shared setup:

```python
# conftest.py
import pytest

@pytest.fixture
def project_dir(tmp_path):
    forge_dir = tmp_path / ".forge"
    forge_dir.mkdir()
    (forge_dir / "config.yaml").write_text("{}")
    return tmp_path
```

### CLI Testing

Use `subprocess.run` to invoke CLI binaries:

```python
import subprocess

def test_cli_command():
    result = subprocess.run(
        ["forge", "subcommand", "--flag", "value"],
        capture_output=True,
        text=True,
    )
    assert result.returncode == 0
    assert "expected output" in result.stdout
```

### API Testing

Use `httpx` or `requests` for HTTP integration testing:

```python
import httpx

def test_api_endpoint():
    response = httpx.get("http://localhost:8080/api/resource")
    assert response.status_code == 200
    assert response.json()["data"] is not None
```

### Traceability

Each test function should include a traceability comment:

```python
def test_login_with_valid_credentials():
    # Traceability: TC-001 -> PRD User Auth section
    ...
```

## assertions

- **style**: assert statement
- **library**: Python built-in `assert` with pytest's assertion introspection
- **custom_matchers**: none (plain assert statements)

### Key Patterns

- `assert actual == expected` — equality check
- `assert actual != expected` — inequality check
- `assert actual in collection` — membership check
- `assert actual is None` — identity check
- `assert actual is not None` — not-None check
- `assert condition` — boolean assertion
- `assert not condition` — negated boolean
- `isinstance(obj, cls)` — type check
- `len(collection) == n` — length check
- `pytest.raises(Exception)` — expected exception context manager

### Rich Comparison (pytest enhancement)

pytest rewrites assert statements to provide detailed failure output:

```python
assert user.name == "Alice"  # Shows full diff on failure
assert "key" in response.json()  # Shows available keys on failure
assert len(items) == 3  # Shows actual length on failure
```

**Rule**: Use plain `assert` statements. Do not import `unittest.AssertionError` or use `unittest` assert methods.

## Tags

- **Format**: `@pytest.mark.<name>` decorators
- **Built-in marks**: `@pytest.mark.slow`, `@pytest.mark.skip`, `@pytest.mark.xfail`
- **Custom marks**: Register in `pyproject.toml` or `pytest.ini`

```python
import pytest

@pytest.mark.slow
def test_large_dataset_processing():
    ...
```

### Mark Registration

```toml
# pyproject.toml
[tool.pytest.ini_options]
markers = [
    "slow: marks tests as slow",
    "e2e: end-to-end tests",
]
```

## Result Format

- **Output flags**: `-v` (verbose), `--tb=short` (traceback format)
- **Format type**: text (default) or `--junitxml=report.xml` for CI integration
- **JSON report**: Use `pytest-json-report` plugin (`--json-report`)

### JUnit XML Structure

```xml
<testsuite name="tests" tests="10" failures="0" errors="0">
  <testcase classname="tests.test_feature" name="test_action" time="0.123" />
</testsuite>
```

## Import Patterns

Standard imports for pytest e2e tests:

```python
import subprocess
import os
import pytest
```

- HTTP tests add: `import httpx` or `import requests`
- File tests add: `from pathlib import Path`
- Fixture tests use `conftest.py` (auto-discovered)

## Anti-patterns (Forbidden)

| Forbidden | Replacement |
|-----------|-------------|
| `time.sleep()` for synchronization | Retry loop with `tenacity` or custom retry |
| `unittest.TestCase` base class | Plain functions with `assert` statements |
| `self.assertEqual()` style | Plain `assert` statements |
| Hardcoded ports | `port = 0` with dynamic allocation or environment variables |
| Real secrets/tokens in code | `os.environ["E2E_API_TOKEN"]` |
| `print()` for debug output | `capsys` fixture or remove entirely |
| `pytest.skip()` without condition | Implement properly or don't generate |
| Mixed assertion styles | Use only `assert` statements, never `unittest` |

## Helpers

### run_cli helper

```python
import subprocess

def run_cli(*args: str, env: dict | None = None) -> subprocess.CompletedProcess:
    """Run a CLI command and return the result."""
    run_env = {**os.environ, **(env or {})}
    return subprocess.run(
        ["forge", *args],
        capture_output=True,
        text=True,
        env=run_env,
    )
```

### retry helper

```python
import time

def retry(fn, max_attempts=3, interval=1.0):
    """Retry a function until it succeeds or max attempts reached."""
    for i in range(max_attempts):
        try:
            return fn()
        except Exception:
            if i == max_attempts - 1:
                raise
            time.sleep(interval)
```
