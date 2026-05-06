# --- forge standard recipes ---

frontend_dir := "FRONTEND_DIR"
backend_dir  := "BACKEND_DIR"

project-type:
    @echo "mixed"

compile scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npx tsc --noEmit ;;
      backend)  cd "{{backend_dir}}" && BACKEND_COMPILE ;;
      "")       (cd "{{frontend_dir}}" && npx tsc --noEmit) && (cd "{{backend_dir}}" && BACKEND_COMPILE) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

build scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npm run build ;;
      backend)  cd "{{backend_dir}}" && BACKEND_BUILD ;;
      "")       (cd "{{frontend_dir}}" && npm run build) && (cd "{{backend_dir}}" && BACKEND_BUILD) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

run scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && FRONTEND_RUN ;;
      backend)  cd "{{backend_dir}}" && BACKEND_RUN ;;
      "")       cd "{{backend_dir}}" && BACKEND_RUN & backend_pid=$!; trap "kill $backend_pid 2>/dev/null" EXIT; sleep 2; cd "{{frontend_dir}}" && FRONTEND_RUN ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

dev scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && FRONTEND_DEV ;;
      backend)  cd "{{backend_dir}}" && BACKEND_DEV ;;
      "")       cd "{{backend_dir}}" && BACKEND_DEV & backend_pid=$!; trap "kill $backend_pid 2>/dev/null" EXIT; sleep 2; cd "{{frontend_dir}}" && FRONTEND_DEV ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

test scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npm test ;;
      backend)  cd "{{backend_dir}}" && BACKEND_TEST ;;
      "")       (cd "{{frontend_dir}}" && npm test) && (cd "{{backend_dir}}" && BACKEND_TEST) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

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

lint scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npm run lint ;;
      backend)  cd "{{backend_dir}}" && BACKEND_LINT ;;
      "")       (cd "{{frontend_dir}}" && npm run lint) && (cd "{{backend_dir}}" && BACKEND_LINT) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

fmt scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npx prettier --write . ;;
      backend)  cd "{{backend_dir}}" && BACKEND_FMT ;;
      "")       (cd "{{frontend_dir}}" && npx prettier --write .) && (cd "{{backend_dir}}" && BACKEND_FMT) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

check scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npm run lint && npx tsc --noEmit ;;
      backend)  cd "{{backend_dir}}" && BACKEND_LINT && BACKEND_COMPILE ;;
      "")       (cd "{{frontend_dir}}" && npm run lint && npx tsc --noEmit) && (cd "{{backend_dir}}" && BACKEND_LINT && BACKEND_COMPILE) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

clean scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && rm -rf dist ;;
      backend)  cd "{{backend_dir}}" && BACKEND_CLEAN ;;
      "")       (cd "{{frontend_dir}}" && rm -rf dist) && (cd "{{backend_dir}}" && BACKEND_CLEAN) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

install scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && npm install ;;
      backend)  cd "{{backend_dir}}" && BACKEND_INSTALL ;;
      "")       (cd "{{frontend_dir}}" && npm install) && (cd "{{backend_dir}}" && BACKEND_INSTALL) ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

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
