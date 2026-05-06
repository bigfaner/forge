# --- forge standard recipes ---

project-type:
    @echo "backend"

compile:
    python -m py_compile src/

build:
    python -m build

run:
    python ENTRY_POINT

dev:
    DEV_COMMAND

test:
    pytest

[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ "{{feature}}" != "" ]; then
        cd tests/e2e && E2E_FEATURE=1 npx playwright test features/{{feature}}/
    else
        [ ! -d tests/e2e/node_modules ] && npm install --prefix tests/e2e
        cd tests/e2e && npx playwright test
    fi

lint:
    ruff check .

fmt:
    ruff format .

check:
    ruff check . && python -m py_compile src/

clean:
    rm -rf build/ dist/ *.egg-info

install:
    pip install -e .

ci: install compile build test lint

e2e-setup force="":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ ! -f tests/e2e/package.json ]; then
        echo "Error: tests/e2e/package.json not found" >&2
        exit 1
    fi
    case "{{force}}" in
      force) npm install --prefix tests/e2e ;;
      "")
        if [ ! -d tests/e2e/node_modules ]; then
            npm install --prefix tests/e2e
        fi
        ;;
      *) echo "[forge] invalid value '{{force}}'; expected 'force' or empty" >&2; exit 1 ;;
    esac
    npx --prefix tests/e2e playwright install chromium
    echo "OK: e2e dependencies ready"

# probe: check if configured services are healthy
probe path="/health":
    #!/usr/bin/env bash
    set -euo pipefail
    config="tests/e2e/config.yaml"
    if [ ! -f "$config" ]; then
        echo "OK: no config.yaml (CLI-only project)"
        exit 0
    fi
    fail=0
    for url in $(grep 'Url:' "$config" | grep -oE 'https?://[^ "]+'); do
        probe_url="${url}{{path}}"
        if curl -sf --max-time 5 "$probe_url" > /dev/null 2>&1; then
            echo "OK: $probe_url"
        else
            echo "FAIL: $probe_url not responding" >&2
            fail=$((fail+1))
        fi
    done
    [ "$fail" -eq 0 ] || exit 1

[arg("feature", long)]
e2e-verify feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -z "{{feature}}" ]; then
        echo "Usage: just e2e-verify --feature <slug>" >&2
        exit 1
    fi
    if [ ! -d "tests/e2e/features/{{feature}}" ]; then
        echo "Error: tests/e2e/features/{{feature}}/ not found" >&2
        exit 1
    fi
    matches=$(grep -rn '// VERIFY:' "tests/e2e/features/{{feature}}/" --include='*.spec.ts' || true)
    if [ -n "$matches" ]; then
        count=$(echo "$matches" | wc -l | tr -d ' ')
        echo "Error: $count unresolved // VERIFY: marker(s) in tests/e2e/features/{{feature}}/" >&2
        echo "" >&2
        echo "$matches" >&2
        echo "" >&2
        echo "Replace each // VERIFY: comment with a real assertion before running tests." >&2
        exit 1
    fi
    echo "OK: no unresolved // VERIFY: markers in tests/e2e/features/{{feature}}/"
    if [ -f "tests/e2e/helpers.ts" ]; then
        hmatches=$(grep -n '// VERIFY:\|// TEMPLATE:' "tests/e2e/helpers.ts" || true)
        if [ -n "$hmatches" ]; then
            echo "Warning: unresolved markers in tests/e2e/helpers.ts:" >&2
            echo "$hmatches" >&2
        fi
    fi

# --- end forge standard recipes ---
