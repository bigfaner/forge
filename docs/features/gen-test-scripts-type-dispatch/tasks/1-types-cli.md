---
id: "1"
title: "Create types/cli.md instruction file"
priority: "P1"
estimated_time: "1-2h"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Create types/cli.md instruction file

## Description

Extract CLI-specific generation logic from the monolithic `gen-test-scripts` SKILL.md into a dedicated `types/cli.md` type instruction file. This file guides the agent when generating CLI test scripts — reconnaissance strategy, Fact Table requirements, generation patterns, verification methods, and antipattern guards.

Modeled after `gen-test-cases/types/cli.md` structure.

## Reference Files
- `docs/proposals/gen-test-scripts-type-dispatch/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-cases/types/cli.md` — Reference structure (gen-test-cases CLI type file)
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — Source of CLI-specific content to extract (lines 249-258 reconnaissance, 280-290 Fact Table, 384-392 verification, 443-465 antipattern guards)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-test-scripts/types/cli.md` | CLI type instruction file with conventions frontmatter |

### Modify
| File | Changes |
|------|---------|
| _(none — extraction happens in task 6)_ | |

## Acceptance Criteria

- [ ] `plugins/forge/skills/gen-test-scripts/types/cli.md` exists
- [ ] Frontmatter declares `type: cli` and `conventions: [testing-cli.md]`
- [ ] Contains a dedicated **Reconnaissance Strategy** section with CLI-specific search patterns (grep cobra.Command, flag parsing, CLI entry points, command registration)
- [ ] Contains a **Fact Table Required Keys** section listing minimum keys for CLI type (at least one CLI command name entry)
- [ ] Contains a **Verification Method** section describing how to confirm the project exposes a CLI interface (grep bin in package.json, ls cmd/, grep cobra.Command)
- [ ] Contains a **Generation Patterns** section describing how CLI test cases translate to executable scripts (process execution, stdout/stderr assertions, exit code checks, argument/flag testing)
- [ ] Contains a **CLI Antipattern Guards** section beyond the generic 6 (recursive test invocation, static file text grep, interactive prompts)
- [ ] At least 3 section headings are unique to this file (not shared with other type files)
- [ ] Content is grounded in the current SKILL.md CLI-specific branches, not invented from scratch

## Hard Rules

- Follow the same structural pattern as `gen-test-cases/types/cli.md` for consistency
- Reconnaissance patterns must cite actual grep commands or search strategies, not vague guidance
- Fact Table keys must be concrete (e.g., `CLI_COMMAND_*`, not "CLI-related entries")

## Implementation Notes

- The current SKILL.md has CLI reconnaissance at the "Required reads" table (line 257: "CLI entry points") — extract and expand this into a full reconnaissance strategy
- The Fact Table Completeness Gate (lines 280-290) has CLI-specific requirements — extract these
- Step 4 verification (lines 384-392) has CLI-specific verification method — extract this
- Antipattern guards 1 (recursive test invocation) and 6 (static file text grep) are CLI-relevant — include in the type-specific antipattern section
- CLI vs TUI disambiguation rules from gen-test-cases `types/cli.md` are useful reference
