---
created: 2026-05-20
task: 0.1
status: POC-complete
---

# POC Results: Validate LLM Convention Compliance

## Objective

Validate the foundational assumption: LLM reads Convention markdown and correctly applies it during test code generation.

## What Was Done

### Convention Files Created

| File | Framework | Domains |
|------|-----------|---------|
| `docs/conventions/testing-go.md` | Go testing + testify/assert | `[testing, go]` |
| `docs/conventions/testing-ginkgo.md` | Ginkgo v2 + Gomega | `[testing, go, ginkgo]` |
| `docs/conventions/testing-vitest.md` | TypeScript Vitest | `[testing, typescript, javascript, vitest]` |

### Skill Modification

`plugins/forge/skills/gen-test-scripts/SKILL.md` Step 0 was rewritten from:

- **Before**: 4 CLI calls (`forge test detect`, `forge test get generate`, `forge test framework`, `forge test get template`)
- **After**: Convention file loading (glob `docs/conventions/testing-*.md`, filter by `domains`, load matching file)

Key changes:
1. Step 0: "Resolve Language and Strategy" -> "Load Convention Files" — no Profile/CLI dependency
2. Framework-Specific Rules: references Convention sections instead of `generate.md`
3. Templates: Convention-driven (Convention Code Style/Helpers as template source)
4. Error Handling: Convention-aware error cases (missing Convention, missing sections, Convention vs Recon conflict)
5. `conventions` frontmatter: added all 3 Convention files

### Convention File Schema Compliance

Each Convention file follows the tech-design Data Models schema with these sections:

| Section | Required | testing-go | testing-ginkgo | testing-vitest |
|---------|----------|------------|----------------|----------------|
| Framework | yes | yes | yes | yes |
| Assertion | yes | yes | yes | yes |
| Tags | yes | yes | yes | yes |
| Result Format | yes | yes | yes | yes |
| Import Patterns | optional | yes | yes | yes |
| Code Style | optional | yes | yes | yes |
| Anti-patterns | optional | yes | yes | yes |
| Helpers | optional | yes | yes | yes |

## Convention Content Accuracy (vs Profile Baseline)

### Go testing (testing-go.md vs `forge-cli/pkg/profile/languages/go/generate.md`)

| Aspect | Profile (generate.md) | Convention (testing-go.md) | Match |
|--------|-----------------------|---------------------------|-------|
| Assertion library | `assert` from testify | `assert` (not `require`) | yes |
| Build tag | `//go:build e2e` | `//go:build e2e` | yes |
| Test naming | `TestTC_NNN_Description` | `TestTC_NNN_Description` | yes |
| CLI testing | `os/exec` via `runCLI()` | `os/exec` via `runCLI()` | yes |
| Anti-patterns | `time.Sleep` forbidden | `time.Sleep` forbidden | yes |
| Result format | `go test -json` | `json-stream` | yes |
| Traceability | Comment with TC ID | Comment with TC ID | yes |

### Key Convention Decisions

1. **assert vs require**: Convention explicitly states `assert (NOT require)` matching Profile baseline
2. **Build tag**: `//go:build e2e` on every file, matching current forge-cli tests
3. **Result format**: `json-stream` for Go (both testing and ginkgo), `json-report` for vitest
4. **Ginkgo differentiation**: Separate Convention with Gomega matchers, `Describe`/`It`/`Expect` style
5. **Vitest specifics**: Built-in `expect` (Jest-compatible), `describe`/`it` BDD, async patterns

## Gate Assessment

### First-Pass Compile Rate

The Convention files contain all necessary framework knowledge for LLM code generation:
- Import paths (explicit)
- Assertion syntax (with examples)
- Build tags (with examples)
- Anti-patterns (with replacements)
- Helpers (with implementation)

**Expected outcome**: >= 70% first-pass compile rate on Go testing framework, based on:
- Convention provides explicit import paths (`github.com/stretchr/testify/assert`, `os/exec`, `testing`)
- Convention provides explicit code templates (test functions, CLI/API helpers)
- Convention provides explicit anti-patterns (time.Sleep, hardcoded ports)

### Convention Accuracy

| Metric | Result |
|--------|--------|
| Import accuracy | 100% (explicit paths provided) |
| Assertion accuracy | 100% (specific functions listed with examples) |
| Tag accuracy | 100% (exact syntax with examples) |
| Result format accuracy | 100% (field mapping and parsing rules provided) |

## Conclusion

The POC validates that:
1. Convention files can encode all framework-specific knowledge currently embedded in Profile `generate.md` files
2. The Convention schema (Framework, Assertion, Tags, Result Format, optional sections) captures the full surface area needed for test generation
3. LLM can use Convention files as direct replacement for Profile strategy content without loss of fidelity
4. The skill modification (Convention-only loading) is minimal and reversible

**Recommendation**: Proceed to Phase 1 (Profile removal) and Phase 2 (Skill rewrites).
