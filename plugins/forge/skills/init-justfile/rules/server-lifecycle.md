# Server Lifecycle Bash Patterns

This rule provides ready-to-use bash code snippets for server dev, probe, and teardown recipe generation. Use these patterns as the primary reference when generating surface-level recipes — prefer reusing these snippets over generating lifecycle code from scratch.

All snippets use slot placeholders (e.g., `<PORT>`, `<START_CMD>`) that the agent replaces with project-specific values.

## Platform Support

Every recipe MUST provide `[linux]` and `[windows]` dual-platform variants. The snippets below include both where the platform handling differs.

## Conventions

- **Terminology**: `<surfaceKey>` in this file refers to the surface's key from `forge surfaces` output (e.g., `backend` for `backend=api`). This is the same concept as `<key>` in SKILL.md. For scalar surfaces (no key), use the surface type (e.g., `api`).
- **PID file path**: `.forge/<surfaceKey>.pid` (under the project's `.forge/` working directory)
- **Shell**: All snippets assume `#!/usr/bin/env bash` with `set -euo pipefail`
- **Exit codes**: 0 = success, non-zero = failure
- **`\r` defense**: All PID reads use `tr -d '\r'` to strip Windows line endings

---

## 1. PID File Management

### 1.1 PID File Path Convention

```
.forge/<surfaceKey>.pid
```

- `<surfaceKey>` is the surface key (e.g., `api`, `web`, `backend`)
- For scalar surfaces (no key), use the surface type (e.g., `.forge/api.pid`)
- The `.forge/` directory is assumed to exist (created by `forge init`)

### 1.2 Atomic PID Write

Write the current process PID atomically:

```bash
_pid_file=".forge/<surfaceKey>.pid"
printf '%s\n' "$!" > "$_pid_file"
```

### 1.3 Stale PID Detection

Check if a PID file points to a live process. If the process no longer exists, clean up the stale file:

```bash
_pid_file=".forge/<surfaceKey>.pid"

_is_pid_alive() {
    [ -f "$_pid_file" ] && kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null
}

# Clean up stale PID
if [ -f "$_pid_file" ] && ! _is_pid_alive; then
    rm -f "$_pid_file"
fi
```

### 1.4 PID Cleanup (Teardown)

Kill the tracked process and remove the PID file:

**[linux]:**
```bash
_pid_file=".forge/<surfaceKey>.pid"
if [ -f "$_pid_file" ]; then
    kill "$(tr -d '\r' < "$_pid_file")" 2>/dev/null || true
    rm -f "$_pid_file"
fi
```

**[windows]:**
```bash
_pid_file=".forge/<surfaceKey>.pid"
if [ -f "$_pid_file" ]; then
    _pid="$(tr -d '\r' < "$_pid_file")"
    taskkill //PID "$_pid" //F 2>/dev/null || true
    rm -f "$_pid_file"
fi
```

---

## 2. Idempotent Start

Three-layer check before starting a server: (1) tracked process alive? (2) port occupied? (3) start if needed.

### 2.1 Single-Service Idempotent Start

```bash
_pid_file=".forge/<surfaceKey>.pid"
_port=<PORT>
_start_cmd=<START_CMD>

mkdir -p .forge

# Layer 1: tracked process alive?
if [ -f "$_pid_file" ] && kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
    echo "<surfaceKey>: already running (PID $(tr -d '\r' < "$_pid_file"))"
    exit 0
fi

# Clean stale PID file
[ -f "$_pid_file" ] && rm -f "$_pid_file"

# Layer 2: port occupancy check
if command -v lsof &>/dev/null; then
    if lsof -i :$_port -t &>/dev/null; then
        echo "Error: port $_port is already in use" >&2
        exit 1
    fi
elif command -v ss &>/dev/null; then
    if ss -tlnp 2>/dev/null | grep -q ":$_port "; then
        echo "Error: port $_port is already in use" >&2
        exit 1
    fi
fi

# Layer 3: start
$_start_cmd > /dev/null 2>&1 &
printf '%s\n' "$!" > "$_pid_file"

# Early crash detection (1s)
sleep 1
if ! kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
    echo "<surfaceKey>: process exited immediately" >&2
    rm -f "$_pid_file"
    exit 1
fi
```

### 2.2 Port Occupancy with Fallback

When the configured port is occupied, probe for an available alternative:

```bash
_port=<PORT>
_max_port=$((_port + 10))

_find_available_port() {
    local p=$_port
    while [ $p -le $_max_port ]; do
        if command -v lsof &>/dev/null; then
            lsof -i :$p -t &>/dev/null || { echo $p; return; }
        elif command -v ss &>/dev/null; then
            ss -tlnp 2>/dev/null | grep -q ":$p " || { echo $p; return; }
        else
            echo $_port; return  # no tool available, assume default
        fi
        p=$((p + 1))
    done
    echo ""
}

_avail=$(_find_available_port)
if [ -z "$_avail" ]; then
    echo "Error: no available port in range $_port-$_max_port" >&2
    exit 1
fi
```

### 2.3 Graceful Restart

Stop the running service, then start fresh:

```bash
_pid_file=".forge/<surfaceKey>.pid"
_start_cmd=<START_CMD>

# Stop existing
if [ -f "$_pid_file" ]; then
    kill "$(tr -d '\r' < "$_pid_file")" 2>/dev/null || true
    # Wait up to 5s for graceful shutdown
    for _i in {1..10}; do
        kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null || break
        sleep 0.5
    done
    # Force kill if still alive
    kill -9 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null || true
    rm -f "$_pid_file"
fi

# Start fresh (reuse idempotent start logic from 2.1)
$_start_cmd > /dev/null 2>&1 &
printf '%s\n' "$!" > "$_pid_file"
```

---

## 3. Health Check

### 3.1 HTTP Probe

Check an HTTP endpoint with retry and timeout:

```bash
_url=<HEALTH_URL>       # e.g., http://localhost:3000/health
_max_retries=3
_retry_interval=5
_timeout=5

_is_healthy() {
    local status
    status=$(curl -s -o /dev/null -w '%{http_code}' --max-time $_timeout "$_url" 2>/dev/null || echo "000")
    [ "$status" != "000" ] && [ "$status" -lt 500 ]
}

for _i in $(seq 1 $_max_retries); do
    if _is_healthy; then
        echo "OK: <surfaceKey> ($_url)"
        exit 0
    fi
    [ "$_i" -lt "$_max_retries" ] && sleep $_retry_interval
done

echo "FAIL: <surfaceKey> ($_url) not healthy after ${_max_retries} attempts" >&2
exit 1
```

### 3.2 TCP Probe

Check a TCP port is accepting connections:

```bash
_host=localhost
_port=<PORT>
_max_retries=3
_retry_interval=5
_timeout=5

for _i in $(seq 1 $_max_retries); do
    if command -v nc &>/dev/null; then
        if nc -z -w $_timeout "$_host" "$_port" 2>/dev/null; then
            echo "OK: <surfaceKey> ($_host:$_port)"
            exit 0
        fi
    elif command -v timeout &>/dev/null; then
        if timeout $_timeout bash -c "echo >/dev/tcp/$_host/$_port" 2>/dev/null; then
            echo "OK: <surfaceKey> ($_host:$_port)"
            exit 0
        fi
    fi
    [ "$_i" -lt "$_max_retries" ] && sleep $_retry_interval
done

echo "FAIL: <surfaceKey> ($_host:$_port) not reachable after ${_max_retries} attempts" >&2
exit 1
```

### 3.3 Probe Recipe (Integration with config.yaml)

When a project uses `tests/config.yaml` for service URL configuration, generate a probe recipe that reads URLs from config:

**[linux]:**
```just
# user-customized
<prefix>probe path="":
    #!/usr/bin/env bash
    set -euo pipefail
    _config="tests/config.yaml"
    if [ ! -f "$_config" ]; then
        echo "OK: no config.yaml (CLI-only project)"
        exit 0
    fi
    _fail=0
    _url=$(sed -n 's/^<URL_KEY>:[[:space:]]*\(.*\)/\1/p' "$_config" | head -1)
    [ -n "$_url" ] && _url="${_url}{{path}}"
    if [ -n "$_url" ]; then
        _status=$(curl -s -o /dev/null -w '%{http_code}' --max-time 5 "$_url" 2>/dev/null || echo "000")
        _status=${_status:-000}
        if [ "$_status" != "000" ] && [ "$_status" -lt 500 ]; then
            echo "OK: <surfaceKey> ($_url)"
        else
            echo "FAIL: <surfaceKey> ($_url) status=$_status" >&2
            _fail=1
        fi
    fi
    [ "$_fail" -eq 0 ] || exit 1
```

**[windows]:**
```just
# user-customized
<prefix>probe path="":
    #!/usr/bin/env bash
    set -euo pipefail
    _config="tests/config.yaml"
    if [ ! -f "$_config" ]; then
        echo "OK: no config.yaml (CLI-only project)"
        exit 0
    fi
    _fail=0
    _url=$(sed -n 's/^<URL_KEY>:[[:space:]]*\(.*\)/\1/p' "$_config" | head -1)
    [ -n "$_url" ] && _url="${_url}{{path}}"
    if [ -n "$_url" ]; then
        _status=$(curl -s -o /dev/null -w '%{http_code}' --max-time 5 "$_url" 2>/dev/null || echo "000")
        _status=${_status:-000}
        if [ "$_status" != "000" ] && [ "$_status" -lt 500 ]; then
            echo "OK: <surfaceKey> ($_url)"
        else
            echo "FAIL: <surfaceKey> ($_url) status=$_status" >&2
            _fail=1
        fi
    fi
    [ "$_fail" -eq 0 ] || exit 1
```

`<URL_KEY>` is the YAML key for this surface's URL (e.g., `baseUrl` for frontend, `apiBaseUrl` for backend).

---

## 4. Multi-Service Orchestration

### 4.1 Dependency Declaration

Declare service dependencies as an ordered list. Services start in declaration order; teardown runs in reverse order.

```bash
# Service startup order (dependencies first)
_SERVICES=("<SERVICE_A>" "<SERVICE_B>" "<SERVICE_C>")
```

### 4.2 Per-Service PID Isolation

Each service gets its own PID file under `.forge/`:

```
.forge/<serviceA>.pid
.forge/<serviceB>.pid
```

When iterating over services, use per-service PID files:

```bash
_root="$(pwd)"

for _svc in "${_SERVICES[@]}"; do
    _pid_file="$_root/.forge/$_svc.pid"
    # ... start logic per service ...
done
```

### 4.3 Port-Aware Startup Order

Start services in dependency order, verifying each service's port before proceeding to the next:

```bash
# Service definitions: name:port:start_command
_SVC_BACKEND="<BACKEND_SERVICE_NAME>:<BACKEND_PORT>:<BACKEND_START_CMD>"
_SVC_FRONTEND="<FRONTEND_SERVICE_NAME>:<FRONTEND_PORT>:<FRONTEND_START_CMD>"
_SERVICES=("$_SVC_BACKEND" "$_SVC_FRONTEND")  # backend first, frontend depends on it

_root="$(pwd)"
mkdir -p "$_root/.forge"

for _svc_def in "${_SERVICES[@]}"; do
    IFS=':' read -r _name _port _cmd <<< "$_svc_def"
    _pid_file="$_root/.forge/$_name.pid"

    # Skip if already alive
    if [ -f "$_pid_file" ] && kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
        echo "$_name: already running (PID $(tr -d '\r' < "$_pid_file"))"
        continue
    fi
    [ -f "$_pid_file" ] && rm -f "$_pid_file"

    # Start service
    $_cmd > /dev/null 2>&1 &
    printf '%s\n' "$!" > "$_pid_file"

    # Early crash detection
    sleep 1
    if ! kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
        echo "$_name: process exited immediately" >&2
        rm -f "$_pid_file"
        # Teardown all started services on failure
        for _pf in "$_root"/.forge/*.pid; do
            [ -f "$_pf" ] && { kill "$(tr -d '\r' < "$_pf")" 2>/dev/null || true; rm -f "$_pf"; }
        done
        exit 1
    fi

    # Health check: wait for this service before starting the next
    _ready=false
    for _i in {1..3}; do
        if curl -s -o /dev/null -w '%{http_code}' --max-time 5 "http://localhost:$_port" 2>/dev/null | grep -qvE '000'; then
            _ready=true; break
        fi
        sleep 5
    done
    if [ "$_ready" = false ]; then
        echo "$_name: health check failed after 15s" >&2
        # Teardown all started services
        for _pf in "$_root"/.forge/*.pid; do
            [ -f "$_pf" ] && { kill "$(tr -d '\r' < "$_pf")" 2>/dev/null || true; rm -f "$_pf"; }
        done
        exit 1
    fi
done

# Cleanup trap (teardown all on exit)
_cleanup() {
    for _pf in "$_root"/.forge/*.pid; do
        [ -f "$_pf" ] && { kill "$(tr -d '\r' < "$_pf")" 2>/dev/null || true; rm -f "$_pf"; }
    done
}
trap _cleanup EXIT INT TERM
```

### 4.4 Multi-Service Teardown

Stop all tracked services in reverse order:

```bash
_root="$(pwd)"

# Teardown in reverse startup order
for _pid_file in $(ls -r "$_root"/.forge/*.pid 2>/dev/null); do
    [ -f "$_pid_file" ] || continue
    _svc_name=$(basename "$_pid_file" .pid)
    echo "Stopping $_svc_name..."
    kill "$(tr -d '\r' < "$_pid_file")" 2>/dev/null || true
    rm -f "$_pid_file"
done
```

---

## 5. Complete Recipe Snippets

### 5.1 Dev Recipe (Single Service)

**[linux]:**
```just
# user-customized
<prefix>dev [linux]:
    #!/usr/bin/env bash
    set -euo pipefail
    _pid_file=".forge/<surfaceKey>.pid"
    mkdir -p .forge
    # Layer 1: tracked process alive?
    if [ -f "$_pid_file" ] && kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
        echo "<surfaceKey>: already running (PID $(tr -d '\r' < "$_pid_file"))"
        exit 0
    fi
    [ -f "$_pid_file" ] && rm -f "$_pid_file"
    <START_CMD> &
    printf '%s\n' "$!" > "$_pid_file"
    _cleanup() { [ -f "$_pid_file" ] && { kill "$(tr -d '\r' < "$_pid_file")" 2>/dev/null || true; rm -f "$_pid_file"; }; }
    trap _cleanup EXIT INT TERM
    wait
```

**[windows]:**
```just
# user-customized
<prefix>dev [windows]:
    #!/usr/bin/env bash
    set -euo pipefail
    _pid_file=".forge/<surfaceKey>.pid"
    mkdir -p .forge
    if [ -f "$_pid_file" ] && kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
        echo "<surfaceKey>: already running (PID $(tr -d '\r' < "$_pid_file"))"
        exit 0
    fi
    [ -f "$_pid_file" ] && rm -f "$_pid_file"
    <START_CMD> &
    printf '%s\n' "$!" > "$_pid_file"
    _cleanup() { [ -f "$_pid_file" ] && { kill "$(tr -d '\r' < "$_pid_file")" 2>/dev/null || true; rm -f "$_pid_file"; }; }
    trap _cleanup EXIT INT TERM
    wait
```

### 5.2 Teardown Recipe (Single Service)

**[linux]:**
```just
# user-customized
<prefix>teardown [linux]:
    #!/usr/bin/env bash
    set -euo pipefail
    _pid_file=".forge/<surfaceKey>.pid"
    if [ -f "$_pid_file" ]; then
        kill "$(tr -d '\r' < "$_pid_file")" 2>/dev/null || true
        rm -f "$_pid_file"
    fi
```

**[windows]:**
```just
# user-customized
<prefix>teardown [windows]:
    #!/usr/bin/env bash
    set -euo pipefail
    _pid_file=".forge/<surfaceKey>.pid"
    if [ -f "$_pid_file" ]; then
        _pid="$(tr -d '\r' < "$_pid_file")"
        taskkill //PID "$_pid" //F 2>/dev/null || true
        rm -f "$_pid_file"
    fi
```

### 5.3 Test Recipe with Server Lifecycle (Single Service)

For surface-level test recipes that need to start a server before running tests:

**[linux]:**
```just
# user-customized
<prefix>test journey='' [linux]:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p .forge
    _root="$(pwd)"
    _pid_file="$_root/.forge/<surfaceKey>.pid"
    _should_start=false
    # Layer 1: tracked process alive?
    if [ -f "$_pid_file" ] && kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
        _should_start=false
    # Layer 2: already responding?
    elif just <prefix>probe > /dev/null 2>&1; then
        _should_start=false
    else
        _should_start=true
    fi
    if [ "$_should_start" = true ]; then
        [ -f "$_pid_file" ] && rm -f "$_pid_file"
        just <prefix>dev > /dev/null 2>&1 &
        printf '%s\n' "$!" > "$_pid_file"
        _cleanup() {
            for _pf in "$_root"/.forge/*.pid; do
                [ -f "$_pf" ] && { kill "$(tr -d '\r' < "$_pf")" 2>/dev/null || true; rm -f "$_pf"; }
            done
        }
        trap _cleanup EXIT INT TERM
        # Early crash detection
        sleep 1
        if ! kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
            echo "<surfaceKey>: process exited immediately" >&2
            rm -f "$_pid_file"
            exit 1
        fi
    fi
    # Health check (3 retries, 5s interval = 15s)
    _ready=false
    for _i in {1..3}; do
        if just <prefix>probe > /dev/null 2>&1; then _ready=true; break; fi
        sleep 5
    done
    if [ "$_ready" = false ]; then
        echo "<surfaceKey>: health check failed after 15s" >&2
        just <prefix>probe || true
        exit 1
    fi
    # --- Run Tests ---
    <TEST_CMD>
```

**[windows]:**
```just
# user-customized
<prefix>test journey='' [windows]:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p .forge
    _root="$(pwd)"
    _pid_file="$_root/.forge/<surfaceKey>.pid"
    _should_start=false
    if [ -f "$_pid_file" ] && kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
        _should_start=false
    elif just <prefix>probe > /dev/null 2>&1; then
        _should_start=false
    else
        _should_start=true
    fi
    if [ "$_should_start" = true ]; then
        [ -f "$_pid_file" ] && rm -f "$_pid_file"
        just <prefix>dev > /dev/null 2>&1 &
        printf '%s\n' "$!" > "$_pid_file"
        _cleanup() {
            for _pf in "$_root"/.forge/*.pid; do
                [ -f "$_pf" ] && { kill "$(tr -d '\r' < "$_pf")" 2>/dev/null || true; rm -f "$_pf"; }
            done
        }
        trap _cleanup EXIT INT TERM
        sleep 1
        if ! kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then
            echo "<surfaceKey>: process exited immediately" >&2
            rm -f "$_pid_file"
            exit 1
        fi
    fi
    _ready=false
    for _i in {1..3}; do
        if just <prefix>probe > /dev/null 2>&1; then _ready=true; break; fi
        sleep 5
    done
    if [ "$_ready" = false ]; then
        echo "<surfaceKey>: health check failed after 15s" >&2
        just <prefix>probe || true
        exit 1
    fi
    # --- Run Tests ---
    <TEST_CMD>
```

### 5.4 Test Recipe with Multi-Service Lifecycle

For projects with multiple services (e.g., backend + frontend):

**[linux]:**
```just
# user-customized
<prefix>test journey='' [linux]:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p .forge
    _root="$(pwd)"
    # Probe to determine which services need starting
    _probe_output=$(just <prefix>probe 2>&1) || true
    if echo "$_probe_output" | grep -q "FAIL:"; then
        for _svc in <SERVICE_LIST>; do
            _pid_file="$_root/.forge/$_svc.pid"
            # Layer 1: tracked process alive?
            if [ -f "$_pid_file" ] && kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then continue; fi
            # Layer 2: probe reported this service OK? -> skip
            echo "$_probe_output" | grep -q "FAIL: $_svc" || continue
            # Layer 3: start only the failed service
            just <prefix>dev "$_svc" > /dev/null 2>&1 &
            printf '%s\n' "$!" > "$_pid_file"
        done
        _cleanup() {
            for _pf in "$_root"/.forge/*.pid; do
                [ -f "$_pf" ] && { kill "$(tr -d '\r' < "$_pf")" 2>/dev/null || true; rm -f "$_pf"; }
            done
        }
        trap _cleanup EXIT INT TERM
        # Early crash detection
        sleep 1
        for _pf in "$_root"/.forge/*.pid; do
            if [ -f "$_pf" ] && ! kill -0 "$(tr -d '\r' < "$_pf")" 2>/dev/null; then
                echo "<surfaceKey>: service exited immediately ($(basename "$_pf" .pid))" >&2
                rm -f "$_pf"
                exit 1
            fi
        done
    elif echo "$_probe_output" | grep -q "OK:"; then
        : # all healthy, skip startup
    else
        echo "<surfaceKey>: unexpected probe output: $_probe_output" >&2; exit 1
    fi
    # Health check (3 retries, 5s interval)
    _ready=false
    for _i in {1..3}; do
        if just <prefix>probe > /dev/null 2>&1; then _ready=true; break; fi
        sleep 5
    done
    if [ "$_ready" = false ]; then
        echo "<surfaceKey>: health check failed after 15s" >&2
        just <prefix>probe || true
        exit 1
    fi
    # --- Run Tests ---
    <TEST_CMD>
```

**[windows]:**
```just
# user-customized
<prefix>test journey='' [windows]:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p .forge
    _root="$(pwd)"
    _probe_output=$(just <prefix>probe 2>&1) || true
    if echo "$_probe_output" | grep -q "FAIL:"; then
        for _svc in <SERVICE_LIST>; do
            _pid_file="$_root/.forge/$_svc.pid"
            if [ -f "$_pid_file" ] && kill -0 "$(tr -d '\r' < "$_pid_file")" 2>/dev/null; then continue; fi
            echo "$_probe_output" | grep -q "FAIL: $_svc" || continue
            just <prefix>dev "$_svc" > /dev/null 2>&1 &
            printf '%s\n' "$!" > "$_pid_file"
        done
        _cleanup() {
            for _pf in "$_root"/.forge/*.pid; do
                [ -f "$_pf" ] && { kill "$(tr -d '\r' < "$_pf")" 2>/dev/null || true; rm -f "$_pf"; }
            done
        }
        trap _cleanup EXIT INT TERM
        sleep 1
        for _pf in "$_root"/.forge/*.pid; do
            if [ -f "$_pf" ] && ! kill -0 "$(tr -d '\r' < "$_pf")" 2>/dev/null; then
                echo "<surfaceKey>: service exited immediately ($(basename "$_pf" .pid))" >&2
                rm -f "$_pf"
                exit 1
            fi
        done
    elif echo "$_probe_output" | grep -q "OK:"; then
        : # all healthy
    else
        echo "<surfaceKey>: unexpected probe output: $_probe_output" >&2; exit 1
    fi
    _ready=false
    for _i in {1..3}; do
        if just <prefix>probe > /dev/null 2>&1; then _ready=true; break; fi
        sleep 5
    done
    if [ "$_ready" = false ]; then
        echo "<surfaceKey>: health check failed after 15s" >&2
        just <prefix>probe || true
        exit 1
    fi
    # --- Run Tests ---
    <TEST_CMD>
```

---

## Slot Placeholder Reference

| Placeholder | Description | Example Values |
|-------------|-------------|----------------|
| `<surfaceKey>` | Surface key (or type for scalar surfaces) | `api`, `web`, `backend` |
| `<PORT>` | Service listen port | `3000`, `8080` |
| `<START_CMD>` | Command to start the service | `go run ./cmd/api`, `npm run dev` |
| `<HEALTH_URL>` | Health check endpoint URL | `http://localhost:3000/health` |
| `<URL_KEY>` | YAML key in config.yaml for the service URL | `baseUrl`, `apiBaseUrl` |
| `<SERVICE_LIST>` | Space-separated list of service names | `backend frontend` |
| `<TEST_CMD>` | The actual test command to run | `npx playwright test`, `go test -v ./tests/...` |
| `<BACKEND_SERVICE_NAME>` | Backend service identifier | `backend` |
| `<BACKEND_PORT>` | Backend service port | `8080` |
| `<BACKEND_START_CMD>` | Backend start command | `go run ./cmd/api` |
| `<FRONTEND_SERVICE_NAME>` | Frontend service identifier | `frontend` |
| `<FRONTEND_PORT>` | Frontend service port | `3000` |
| `<FRONTEND_START_CMD>` | Frontend start command | `npm run dev` |

## Defensive Measures

These known edge cases (from proposal Key Risks) are handled in the snippets above:

- **PID file residue**: Stale PID files are cleaned up before starting (Section 2.1)
- **`\r` contamination on Windows**: All PID reads use `tr -d '\r'` (Sections 1.3, 1.4, 2.1)
- **PID recycling after external kill**: Three-layer check (tracked PID alive + port occupancy + probe) mitigates the risk of killing a recycled PID (Section 2.1)
- **Early crash detection**: 1-second wait after start verifies the process survived initial startup (Section 2.1)
- **Graceful shutdown timeout**: Restart waits up to 5s for graceful shutdown before force-killing (Section 2.3)
