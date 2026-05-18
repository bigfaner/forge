# Intent Propagation

The proposal frontmatter may contain an `intent` field (e.g., `intent: cleanup`). When present, use it as the **default type** for all tasks in this feature:

1. Read `proposal.md` frontmatter → extract `intent` value
2. If `intent` is set and matches a valid type constant (`feature`, `enhancement`, `cleanup`, `refactor`) → use it as the default `type` for all business tasks
3. Individual task frontmatter `type` field **overrides** the proposal intent — use it when the task's primary output differs from the feature's dominant intent
4. If `intent` is empty or missing → fall back to per-task Type Assignment from the table above

The mapping is 1:1: proposal intent values use the same names as task type constants.
