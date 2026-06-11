---
status: "completed"
started: "2026-06-08 22:14"
completed: "2026-06-08 22:17"
time_spent: "~3m"
---

# Task Record: 2 Rewrite SKILL.md for agent-driven justfile generation

## Summary
Rewrote SKILL.md to replace template-driven generation with agent-driven generation: removed --type parameter and project type detection, made surfaces a prerequisite with forge init prompt, replaced Step 0/3 template flow with marker-file language detection + Convention loading + agent-driven recipe synthesis, referenced rules/server-lifecycle.md for server lifecycle patterns, and added post-generation consistency verification (Step 4a) that validates generated recipes against surface rule Recipe Invocation Contracts.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/SKILL.md

### Key Decisions
无

## Document Metrics
527 lines (was 491), 6 AC items all met

## Referenced Documents
- docs/proposals/agent-driven-justfile-generation/proposal.md
- plugins/forge/skills/init-justfile/rules/server-lifecycle.md
- docs/conventions/forge-distribution.md
- plugins/forge/skills/init-justfile/rules/surfaces/api.md
- plugins/forge/skills/init-justfile/rules/surfaces/web.md

## Review Status
final

## Acceptance Criteria
- [x] --type parameter removed from frontmatter, Parameters table, and all references
- [x] Project type detection step (Step 1a) removed along with rules/project-detection.md, FRONTEND_DIR/BACKEND_DIR/BACKEND_ENTRY/FRONTEND_RUN_SCRIPT references
- [x] Surfaces as prerequisite: Outcome B prompts user to run forge init instead of silently skipping
- [x] Step 3 agent-driven: marker file language detection (Step 0a), Convention loading (Step 0b), agent synthesizes recipes from detected language + Convention + surface rules + agent knowledge
- [x] References rules/server-lifecycle.md via Step 0c HARD-RULE, replaces template-based server lifecycle code, Step 0 template loading flow removed
- [x] Post-generation consistency verification (Step 4a): validates recipe names/signatures/platform variants against surface rule Recipe Invocation Contract for init-justfile/run-tests dual-consumer consistency

## Notes
Structural changes: Process Flow renumbered from 0/1/1s/2/3/4/5 to 0/1s/2/3/4/5 (Step 1 removed). Step 0 split into 0a (language detection), 0b (Convention loading), 0c (server lifecycle). Step 4 expanded from 2-phase to 3-phase verification (consistency + dry-run + execution). Backward compatibility preserved: boundary markers, recipe naming, dual-platform variants, user-customized markers, exit code semantics unchanged.
