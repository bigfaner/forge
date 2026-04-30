# --- forge standard recipes ---

project-type:
    @echo "backend"

compile:
    @echo "Error: compile recipe not implemented — edit justfile to add build command" >&2 && exit 1

build:
    @echo "Error: build recipe not implemented — edit justfile to add build command" >&2 && exit 1

run:
    @echo "Error: run recipe not implemented — edit justfile to add run command" >&2 && exit 1

dev:
    @echo "Error: dev recipe not implemented — edit justfile to add dev command" >&2 && exit 1

test:
    @echo "Error: test recipe not implemented — edit justfile to add test command" >&2 && exit 1

[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ "{{feature}}" != "" ]; then
        scripts_dir="tests/e2e/{{feature}}"
        fail=0
        for spec in "$scripts_dir"/*.spec.ts; do
            [ -f "$spec" ] && npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    else
        [ ! -d tests/e2e/node_modules ] && npm install --prefix tests/e2e
        fail=0
        for spec in $(find tests/e2e -mindepth 2 -name '*.spec.ts'); do
            npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    fi

lint:
    @echo "Error: lint recipe not implemented — edit justfile to add lint command" >&2 && exit 1

fmt:
    @echo "Error: fmt recipe not implemented — edit justfile to add fmt command" >&2 && exit 1

check:
    @echo "Error: check recipe not implemented — edit justfile to add check command" >&2 && exit 1

clean:
    @echo "Error: clean recipe not implemented — edit justfile to add clean command" >&2 && exit 1

install:
    @echo "Error: install recipe not implemented — edit justfile to add install command" >&2 && exit 1

ci: install compile build test lint

e2e-setup:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ ! -f tests/e2e/package.json ]; then
        echo "Error: tests/e2e/package.json not found" >&2
        exit 1
    fi
    if [ ! -d tests/e2e/node_modules ]; then
        npm install --prefix tests/e2e
    fi
    npx --prefix tests/e2e playwright install chromium
    echo "OK: e2e dependencies ready"

[arg("feature", long)]
e2e-verify feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -z "{{feature}}" ]; then
        echo "Usage: just e2e-verify --feature <slug>" >&2
        exit 1
    fi
    if [ ! -d "tests/e2e/{{feature}}" ]; then
        echo "Error: tests/e2e/{{feature}}/ not found" >&2
        exit 1
    fi
    matches=$(grep -rn '// VERIFY:' "tests/e2e/{{feature}}/" --include='*.spec.ts' || true)
    if [ -n "$matches" ]; then
        count=$(echo "$matches" | wc -l | tr -d ' ')
        echo "Error: $count unresolved // VERIFY: marker(s) in tests/e2e/{{feature}}/" >&2
        echo "$matches" >&2
        exit 1
    fi
    echo "OK: no unresolved // VERIFY: markers in tests/e2e/{{feature}}/"

# --- end forge standard recipes ---
