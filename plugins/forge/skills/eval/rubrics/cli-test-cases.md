---
scale: 1000
target: 900
iterations: 6
type: cli-test-cases
context:
  conventions: [cli, testing-isolation]
  business-rules: auto
---

# CLI Test Cases Evaluation Rubric

**Total: 1000 points**

## Required Sections

The cli-test-cases.md must contain these sections:

- [ ] Frontmatter with `feature`, `sources`, `generated`
- [ ] CLI Test Cases section with individual test cases
- [ ] Traceability table (TC ID → Source → Type → Target → Priority)

**Missing section**: 0 pts for every dimension that depends on the missing section.

## Dimensions

### 1. PRD Traceability (200 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| TC-to-AC mapping exists | 0-70 | Every TC has a `Source` field pointing to a specific PRD acceptance criterion. Not just "PRD" but "PRD AC-3.1" level specificity |
| Traceability table complete | 0-70 | Traceability table lists every TC with its PRD source, type, target, and priority. No TCs missing from the table |
| Reverse coverage | 0-60 | Every PRD acceptance criterion has at least one TC. No AC is orphaned — check against prd-user-stories.md and prd-spec.md |

### 2. Step Actionability (250 pts)

**Blocking threshold**: If this dimension scores < 200, downstream gen-test-scripts is blocked.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Steps are concrete actions | 0-90 | Each step describes a single, unambiguous CLI invocation. "`forge task list --status pending`" not "List pending tasks". "`forge task add --title 'test'`" not "Create a task" |
| Expected results are verifiable | 0-90 | Every expected result can be objectively verified: specific exit code, stdout content, stderr content. Not "should work" or "command succeeds" |
| Preconditions are explicit | 0-70 | TCs with dependencies (existing config, initialized project, specific state) declare them in Pre-conditions. No implicit assumptions |

### 3. Command Coverage Accuracy (150 pts)

This dimension also checks compliance with injected project conventions for CLI testing. The scorer should reference injected conventions to detect violations in command invocation patterns and output assertions.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Command coverage | 0-50 | All flags, subcommands, and argument combinations are tested. Every CLI command mentioned in the PRD has at least one TC. Flag combinations (short and long form) are covered. Required arguments vs optional arguments are distinguished in TCs |
| Output assertion specificity | 0-50 | Exit codes, stdout/stderr content, and error messages are explicitly asserted. Every TC specifies the expected exit code (0 for success, non-zero for errors). Stdout assertions use concrete text or regex patterns. Error message assertions specify the exact or pattern-matched error output |
| Convention compliance | 0-50 | Do test steps and assertions comply with project conventions for CLI testing? If injected conventions specify CLI output format (e.g., JSON output structure, table formatting), do TCs assert against that format? If conventions declare exit code ranges or error message patterns, are they tested? Deduct 10 pts per convention violation |

### 4. Completeness (200 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Boundary and edge cases | 0-70 | Where the PRD explicitly mentions error states, empty states, or boundary conditions, at least one TC covers each. Do not invent scenarios not present in the PRD |
| Integration scenarios | 0-70 | TCs cover cross-feature or cross-interface scenarios (e.g., CLI command triggers file system change, CLI output feeds into another tool) where applicable and mentioned in the PRD |
| CLI coverage breadth | 0-60 | All CLI commands and subcommands described in the PRD have corresponding test cases. Every command, flag, and argument pattern mentioned in the PRD is tested |

### 5. Structure & ID Integrity (100 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| TC IDs are sequential and unique | 0-40 | IDs follow the pattern (e.g., TC-001, TC-002...). No gaps, no duplicates, no re-used IDs |
| Classification is correct | 0-30 | Each TC is classified as CLI type. No UI, API, TUI, or Mobile TCs in this file |
| Summary table matches actual | 0-30 | Counts in the summary table match the actual number of TCs in the section |

### 6. Antipattern Prevention (100 pts)

This dimension evaluates whether test cases are designed to avoid common downstream antipatterns in `/gen-test-scripts`. Well-designed test cases prevent these issues upstream, making script generation more reliable.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Pre-conditions are concrete and creatable | 0-30 | Score only pre-conditions that exist in the document — missing pre-conditions are penalized under D2 (Step Actionability). For pre-conditions that ARE listed, every one must describe HOW to create the required state (e.g., "a temp directory with `forge init` already run"), not just assert it exists (e.g., "forge project exists"). If a pre-condition cannot be created using an isolated test fixture (temp directory, test container, mock server, etc.), the downstream script will generate a conditional skip without fixture. Deduct 10 pts per non-creatable pre-condition |
| Steps describe runtime behavior | 0-25 | No step describes reading source files (`.md`, `.go`, `.json`), checking documentation content, or verifying file existence. Every step interacts with the running CLI (command execution, flag parsing). Deduct 10 pts per static-file-check step |
| No duplicate scenarios | 0-20 | No two TCs test the same scenario with identical inputs and conditions. Duplicate TCs generate duplicate test functions that double CI time. Deduct 10 pts per duplicate pair |
| No meta-testing | 0-15 | No TC verifies test infrastructure ("all tests pass", "test suite compiles", "config is valid"). Every TC must test product behavior. Meta-tests cause recursive test invocation. Deduct 15 pts per meta-test TC |
| Every TC is implementable | 0-10 | No TC describes a scenario that requires unavailable infrastructure (e.g., "real production environment") without marking itself as non-implementable. Dead TCs generate unconditional skips — worse than no test at all. Deduct 5 pts per non-implementable TC without annotation |

## Deduction Rules

- **Score floor**: No criterion or dimension score may fall below 0. Clamp to 0 after applying all deductions.
- **Missing required section**: 0 pts for that dimension
- **Vague language without specificity**: -20 pts per instance ("run command" without specifying which, "check output" without expected content)
- **Cross-section inconsistency**: -30 pts per conflict (e.g., traceability table says TC-005 is API but it's listed under CLI)
- **Placeholder text ("TBD", "TODO", "N/A")**: -20 pts per instance
