---
title: "Skill Self-Containment"
domains: [skill, SKILL.md, command, convention, self-contained]
---

# Skill Self-Containment

Each skill/command must be logically self-contained — a reader should understand the full workflow by reading that single file, without needing to cross-reference other skills or shared references.

Therefore, cross-skill instruction duplication is expected and acceptable. Do not refactor shared logic into reference files solely to reduce duplication. Duplication is the correct trade-off when it preserves self-containment.

## Examples

**Acceptable duplication**: If both `/execute-task` and `/fix-bug` need to describe the Quality Gate sequence (compile -> fmt -> lint -> test), each SKILL.md should include the full sequence inline. A reader of `/fix-bug` should not need to open `/execute-task` to understand the gate steps.

**Violation**: A SKILL.md that says "follow the Quality Gate protocol defined in the execute-task skill" without inlining the actual steps. This forces the agent to cross-reference another skill, breaking self-containment.

**Acceptable reference**: A SKILL.md may reference external files that are *data sources* (e.g., `rubrics/prd.md`, `templates/decision-entry.md`) — these are inputs to the skill's workflow, not instructions the reader must understand to execute the workflow. The SKILL.md should still describe what the file contains and how to use it.
