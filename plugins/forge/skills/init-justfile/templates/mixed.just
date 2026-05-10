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
      "")       cd "{{backend_dir}}" && BACKEND_RUN & _bp=$!; _cleanup() { kill "$_bp" 2>/dev/null || true; }; trap _cleanup EXIT INT TERM; sleep 2; cd "{{frontend_dir}}" && FRONTEND_RUN ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac

dev scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) cd "{{frontend_dir}}" && FRONTEND_DEV ;;
      backend)  cd "{{backend_dir}}" && BACKEND_DEV ;;
      "")       cd "{{backend_dir}}" && BACKEND_DEV & _bp=$!; _cleanup() { kill "$_bp" 2>/dev/null || true; }; trap _cleanup EXIT INT TERM; sleep 2; cd "{{frontend_dir}}" && FRONTEND_DEV ;;
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
    # --- Server Lifecycle (idempotent, per-service) ---
    mkdir -p tests/e2e/results
    _root="$(pwd)"
    probe_output=$(just probe 2>&1) || true
    if echo "$probe_output" | grep -q "FAIL:"; then
        for svc in frontend backend; do
            pid_file="$_root/tests/e2e/results/.pid-$svc"
            # Layer 1: tracked process alive?
            if [ -f "$pid_file" ] && kill -0 "$(tr -d '\r' < "$pid_file")" 2>/dev/null; then continue; fi
            # Layer 2: probe reported this service OK? → skip
            echo "$probe_output" | grep -q "FAIL: $svc" || continue
            # Layer 3: start only the failed service
            just run "$svc" > /dev/null 2>&1 &
            printf '%s\n' "$!" > "$pid_file"
        done
        # trap: cleanup on exit, interrupt, terminate
        _cleanup() {
            for f in "$_root"/tests/e2e/results/.pid-*; do
                [ -f "$f" ] && { kill "$(tr -d '\r' < "$f")" 2>/dev/null || true; rm -f "$f"; }
            done
        }
        trap _cleanup EXIT INT TERM
        # Early crash detection: check if any started service died immediately
        sleep 1
        for pid_f in "$_root"/tests/e2e/results/.pid-*; do
            if [ -f "$pid_f" ] && ! kill -0 "$(tr -d '\r' < "$pid_f")" 2>/dev/null; then
                echo "e2e: service process exited immediately ($(basename "$pid_f" | sed 's/.pid-//'))" >&2
                rm -f "$pid_f"
                exit 1
            fi
        done
        # Health check (wait for ALL services, 30s timeout)
        ready=false
        for i in {1..10}; do
            if just probe > /dev/null 2>&1; then ready=true; break; fi
            sleep 3
        done
        if [ "$ready" = false ]; then
            echo "e2e: health check failed after 30s" >&2
            just probe || true
            exit 1
        fi
    elif echo "$probe_output" | grep -q "OK:"; then
        : # all healthy, skip startup
    else
        echo "e2e: unexpected probe output: $probe_output" >&2; exit 1
    fi
    # --- Run Tests ---
    if [ "{{feature}}" != "" ]; then
        cd tests/e2e && E2E_FEATURE=1 npx playwright test features/{{feature}}/
    else
        if [ ! -d tests/e2e/node_modules ]; then npm install --prefix tests/e2e; fi
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
probe path="":
    #!/usr/bin/env bash
    set -euo pipefail
    config="tests/e2e/config.yaml"
    if [ ! -f "$config" ]; then
        echo "OK: no config.yaml (CLI-only project)"
        exit 0
    fi
    fail=0
    frontend=$(sed -n 's/^baseUrl:[[:space:]]*\(.*\)/\1/p' "$config" | head -1)
    backend=$(sed -n 's/^apiBaseUrl:[[:space:]]*\(.*\)/\1/p' "$config" | head -1)
    [ -n "$frontend" ] && frontend="${frontend}{{path}}"
    [ -n "$backend" ] && backend="${backend}{{path}}"
    for label in "frontend:$frontend" "backend:$backend"; do
        url="${label#*:}"
        [ -z "$url" ] && continue
        STATUS=$(curl -s -o /dev/null -w '%{http_code}' --max-time 5 "$url" 2>/dev/null || echo "000")
        STATUS=${STATUS:-000}
        if [ "$STATUS" != "000" ] && [ "$STATUS" -lt 500 ]; then
            echo "OK: ${label%%:*} ($url)"
        else
            echo "FAIL: ${label%%:*} ($url) status=$STATUS" >&2
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
