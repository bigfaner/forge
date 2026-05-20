# Architecture Decisions

| Date | Feature | Decision | Rationale | Source |
|------|---------|----------|-----------|--------|
| 2026-04-30 | justfile-standard-vocabulary | Defer `task scope` command until scope resolution fails in practice | Current wiring (SCOPE in dispatch prompt + guide.md protocol) may be sufficient; avoid premature complexity | manual |
| 2026-05-19 | simplify-breakdown-tasks-prompt | Decompose skills with conditionally-gated rules into skeleton + `rules/` directory with condition-rule matrix | Reduces token cost (23KB→7.7KB for base case), improves execution stability by eliminating LLM tag-evaluation errors. Chose LLM-instructed conditional loading over CLI-driven assembly to preserve skill self-containment. Qualifies `docs/conventions/skill-self-containment.md`: skeleton is complete for base workflows; rule files are additive. | run-tasks |
| 2026-05-20 | system-type-exclusion | Reverse exclusion (SystemTypes blacklist) over forward whitelist (BusinessTypes) for task type validation | System types (13) form a stable closed set owned by forge pipeline; business types grow with user needs. Blacklist avoids updating for every new business type. Dual-identity types (doc.consolidate, doc.drift) excluded from blacklist to allow manual creation in legacy projects. | run-tasks |
