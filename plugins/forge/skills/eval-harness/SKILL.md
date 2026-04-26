---
name: eval-harness
description: Harness health evaluation centered on three core principles - Design Environment, Clarify Intent, Build Feedback Loops. Based on OpenAI's harness engineering practices.
---

# Eval Harness

Evaluate project harness health based on OpenAI's three core Harness Engineering principles.

## When to Use

**Trigger:**
- User asks to "evaluate harness" or "check harness health"
- User provides `/eval-harness` command
- Periodic health check (recommended: weekly)

**Skip:**
- User wants to implement improvements directly (use `/improve-harness`)

## Core Framework

```
┌─────────────────────────────────────────────────────────────┐
│                    HARNESS ENGINEERING                       │
├─────────────────────────────────────────────────────────────┤
│   Design Environment → Clarify Intent → Build Feedback Loops│
│         │                      │                    │        │
│         ▼                      ▼                    ▼        │
│   Structure, boundaries,  Plans, specs,         Validation, │
│   tools, abstractions     principles,           correction, │
│                            constraints           observability│
└─────────────────────────────────────────────────────────────┘
```

## Workflow

```
1. Check Environment → 2. Check Intent → 3. Check Feedback → 4. Generate Report
```

### Step 1: Design Environment

| Dimension | Checks | Command |
|-----------|--------|---------|
| Progressive disclosure | CLAUDE.md < 100 lines, has paths frontmatter | `wc -l CLAUDE.md` |
| Architecture boundaries | Dependency direction defined, boundary check scripts | `ls scripts/lint*.sh 2>/dev/null` |
| Tool abstractions | skills, agents exist | `ls .claude/skills/ .claude/agents/ 2>/dev/null` |
| Agent readability | Architecture visible in code | Project-specific command |

### Step 2: Clarify Intent

| Dimension | Checks | Command |
|-----------|--------|---------|
| Golden principles | Principles have enforcement mechanisms | Check CLAUDE.md or project rules |
| Plan artifacts | Tasks have schema validation | `ls docs/features/*/tasks/index.json 2>/dev/null` |
| Invariants | Lint, format automated | Check Makefile or CI config |

### Step 3: Build Feedback Loops

| Dimension | Checks | Command |
|-----------|--------|---------|
| Immediate feedback | hooks configured | `cat .claude/settings.json 2>/dev/null` |
| Auto-fix | error-fixer agent | `ls .claude/agents/error-fixer.md 2>/dev/null` |
| Debt GC | code simplification skill | `grep -r simplify .claude/ 2>/dev/null` |
| Observability | Structured test results | Project-specific command |

### Step 4: Generate Report

Output to `docs/harness-reports/YYYY-MM-DD.md` using `templates/report.md` template.

## Language Detection

Auto-detect project language and use corresponding check commands:

| Language | Detection Marker | Format Check | Lint | Test |
|----------|-----------------|--------------|------|------|
| Go | `go.mod` | `gofmt` | `go vet`, `golangci-lint` | `go test` |
| Node.js | `package.json` | `prettier` | `eslint` | `npm test` |
| Python | `pyproject.toml` | `black` | `ruff`, `pylint` | `pytest` |
| Rust | `Cargo.toml` | `rustfmt` | `clippy` | `cargo test` |
| Java | `pom.xml` | `google-java-format` | `checkstyle` | `mvn test` |

## Grading Rules

### Sub-dimension

| Grade | Condition |
|-------|-----------|
| A | All checks pass with automation |
| B | Most pass, minor manual intervention |
| C | Defined but enforcement incomplete |
| F | Critical checks missing |

### Overall

| Grade | Condition |
|-------|-----------|
| A | All 3 dimensions A/B, at least 2 A's |
| B | No F, max 2 C's |
| C | 1 F or 3+ C's |
| D | 2 F's |
| F | 3+ F's or core capability missing |

## Output

After evaluation:

1. Create report: `docs/harness-reports/YYYY-MM-DD.md`
2. Update link in `docs/HARNESS-EVALUATION.md`
3. List priority improvements (P0/P1/P2)

## Related

- `/improve-harness` - Implement improvements from report
- `docs/HARNESS-EVALUATION.md` - Current evaluation summary
- `docs/harness-reports/` - Historical reports
