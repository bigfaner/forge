claude:
    claude --dangerously-skip-permissions

claude-c:
    claude --dangerously-skip-permissions -c

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

# Requires just >= 1.50.0. [arg("feature", long)] declares a long-form CLI flag:
# --feature <value>. "long" maps the argument to --<name> form; omitting it makes
# the argument positional.
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