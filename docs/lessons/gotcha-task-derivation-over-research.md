---
created: 2026-05-15
tags: [architecture, testing]
---

# quick-tasks Instruction "Determine" Triggers Over-Research

## Problem

During quick-tasks, the agent read excessive documentation (gen-test-scripts SKILL.md 430 lines, breakdown-tasks SKILL.md, dispatch logic) before writing task files. This wasted context window and time.

## Root Cause

Causal chain (5 levels):

1. **Symptom**: Agent read too many downstream files during quick-tasks Step 2
2. **Direct cause**: Agent needed exact file paths to fill the task template's Create/Modify/Delete tables
3. **Trigger instruction**: quick-tasks SKILL.md Step 2 says "**Determine** affected file paths from the solution description" — the verb "determine" implies the agent must *discover* paths that aren't in the input
4. **Input lacks paths**: The proposal contains capability-level descriptions ("Add --type argument to gen-test-scripts"), not file-level paths. The "solution description" does not contain the paths the instruction tells the agent to "determine"
5. **Deepest cause**: **The instruction verb "determine" creates an implicit mandate to research beyond the input.** Compare with breakdown-tasks which says "**Inspect** the task's affected file paths (derived from the tech-design)" — "inspect" means paths already exist in the input, just read them. "Determine" means paths must be found. When the input doesn't contain them, the agent reads downstream skills to find them.

### Comparison: breakdown-tasks avoids this via instruction wording

| Factor | quick-tasks | breakdown-tasks |
|--------|-------------|-----------------|
| Instruction verb | "**Determine**" file paths | "**Inspect**" file paths |
| Implication | Paths must be discovered | Paths already exist in input |
| Input document | Proposal (no file paths) | Tech design (exact file paths listed) |
| Agent behavior | Over-research to discover paths | Copy paths from input |

## Solution

Change the instruction verb in quick-tasks SKILL.md Step 2 from "determine" to "infer", and add a degradation rule:

**Before**: "Determine affected file paths from the solution description"

**After**: "Infer affected file paths from the solution description. Use directory-level paths (e.g., `plugins/forge/skills/gen-test-scripts/`) when exact file paths are not specified in the proposal. Do NOT read referenced skill files solely to discover exact file paths."

## Reusable Pattern

**Principle**: Skill instruction verbs must match the input's information density. "Determine" implies discovery and mandates research. "Infer" implies reasoning from available information. "Inspect" implies reading what's already there.

**Detection rule**: When writing skill instructions that reference extracting information from input documents, audit the verb:
- **"Determine" / "Find" / "Discover"** — agent will read beyond the input to find the answer
- **"Infer" / "Derive"** — agent will reason from available information
- **"Inspect" / "Read" / "Extract"** — agent will use only what's in the input

Choose the verb that matches what the input actually contains.

## Example

```
# Proposal: "Add --type argument to gen-test-scripts skill"

# "Determine" path → reads gen-test-scripts/SKILL.md (430 lines) → discovers exact files
# "Infer" path  → infers plugins/forge/skills/gen-test-scripts/ from skill name → done
```

## Related Files

- plugins/forge/skills/quick-tasks/SKILL.md (Step 2, "Determine affected file paths" instruction)
- plugins/forge/skills/breakdown-tasks/SKILL.md (Step 4a, "inspect the task's affected file paths" for comparison)
- plugins/forge/skills/quick-tasks/templates/task.md (Affected Files section that triggers the need)
