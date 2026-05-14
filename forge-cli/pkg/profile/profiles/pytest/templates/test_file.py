"""E2E test template for pytest.

Place in tests/e2e/features/<slug>/test_<feature>.py
Run with: python -m pytest tests/e2e/features/<slug>/ -v
"""

import pytest
import subprocess
import os


def test_tc_nnn_description():
    """Traceability: TC-NNN → {PRD Source}"""
    # Step 1: Setup
    # Step 2: Execute
    result = subprocess.run(
        ["binary", "--flag", "value"],
        capture_output=True, text=True
    )

    # Expected: ...
    assert result.returncode == 0, f"command failed: {result.stderr}"
    assert "expected text" in result.stdout
