---
created: 2026-05-16
author: faner
status: Draft
---

# Proposal: Project Knowledge Discovery via Frontmatter Domain Tags

## Problem

Skills and prompt templates hardcode how to find and load `docs/conventions/` and `docs/business-rules/` files, using either explicit keyword→filename mapping tables (fix-bug.md line 51) or vague "read relevant files" instructions (implementation.md, fix.md). When `/consolidate-specs` creates new convention/business-rule files, existing skills cannot discover them without manual updates to their mapping tables.

### Evidence

- `fix-bug.md:51` has a hardcoded keyword→filename table: `"API"/"endpoint"/"route" → conventions/api.md`, `"error"/"status code" → conventions/error-handling.md`, etc.
- `implementation.md:12` and `fix.md:20` say "Read relevant project knowledge files" with no guidance on how to determine relevance
- `consolidate-specs-alignment.md` duplicates a similar mapping table
- `gen-test-cases-per-type-dispatch` proposal adds per-type `conventions` frontmatter fields — solving this for one skill only, not system-wide
- New conventions added via `/consolidate-specs` (e.g., a future `docs/conventions/api.md`) are invisible to `fix-bug.md` and prompt templates until someone manually updates the mapping tables

### Urgency

The `gen-test-cases-per-type-dispatch` feature is about to add convention loading to its skills. Without a shared discovery mechanism, each new skill that needs conventions will invent its own mapping approach, creating more duplication. Solving this now prevents N×M hardcoding (N skills × M convention files).

## Proposed Solution

Add a `domains` field to the YAML frontmatter of every `docs/conventions/` and `docs/business-rules/` file. Each file self-describes the domains it covers. Skills and prompt templates use a shared discovery instruction: glob the directory, read frontmatter, match domains against task context, load matched files.

`/consolidate-specs` auto-manages the frontmatter when creating or updating files — no separate sync step needed.

### Frontmatter Schema

```yaml
---
title: "Error Handling Conventions"  # existing field, unchanged
domains: [error, status, response, stderr]  # NEW: keywords this file covers
---
```

### Shared Discovery Instruction

A reusable instruction snippet that replaces all hardcoded mapping tables:

```
Glob docs/conventions/*.md and docs/business-rules/*.md.
For each file, read its YAML frontmatter and extract the `domains` list.
Match the file's domains against keywords from the task description, affected files, and error context.
Load files whose domains overlap with the task context.
If a file has no `domains` frontmatter, load it (backward compatibility).
If no files match, skip — no matching convention files for this task.
```

### Innovation Highlights

This follows the **self-describing artifact** pattern common in content management systems. Each document declares its own metadata, eliminating the need for external registries or mapping tables. The approach is straightforward but effective — it replaces N×M hardcoded mappings with N skill reads of M self-describing files.

The key insight is that `consolidate-specs` already owns the lifecycle of these files, making it the natural place to manage the metadata. No new sync points or maintenance burden.

## Requirements Analysis

### Key Scenarios

- **Task execution**: Agent receives T-impl-3 with scope "backend" and keywords "error handling, CLI output". Discovery loads `error-handling.md` (domains match) and `error-reporting.md` (domains match), skips `profile-system.md` (no match)
- **Bug fix**: User runs `/fix-bug "TypeError in task claim"`. Discovery loads `task-lifecycle.md` (domains: task, claim, status) and `error-handling.md` (domains: error)
- **New convention**: `/consolidate-specs` creates `docs/conventions/api.md` with `domains: [api, endpoint, route, rest]`. Next task execution automatically discovers it — no skill updates needed
- **No matching files**: Task about CSS layout in a frontend-only project with no CSS conventions file. Discovery finds no match, skips loading, task proceeds normally
- **Legacy file without frontmatter**: An older file missing the `domains` field is loaded unconditionally (backward compatibility)

### Non-Functional Requirements

- **Token efficiency**: Discovery reads only frontmatter (first ~5 lines) before deciding to load the full file. No full-file reads for irrelevant documents
- **Backward compatibility**: Files without `domains` frontmatter are loaded unconditionally, preserving current behavior
- **Extensibility**: Adding new convention files requires zero changes to any skill or prompt template

### Constraints & Dependencies

- All existing convention/business-rule files must be updated with `domains` frontmatter
- `consolidate-specs` SKILL.md must be updated to write `domains` when creating/updating files
- Prompt templates (`implementation.md`, `fix.md`) are plain text — the discovery instruction must be embeddable as-is

## Alternatives & Industry Benchmarking

### Industry Solutions

Content-addressable systems (e.g., Hugo frontmatter taxonomies, Jekyll categories) use self-describing metadata on documents. Plugin systems (VS Code extensions, ESLint configs) use manifest files with capability declarations. This proposal uses the simpler frontmatter approach since the file count is small (<20 files typically).

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero migration cost | Each new convention requires N skill updates; mappings drift | Rejected: doesn't scale |
| Central registry file | Plugin manifests | Single lookup point | Separate file to maintain; can drift from actual files | Rejected: single point of failure |
| Glob-and-load all | Simple scripts | Simplest implementation | Wastes context on irrelevant docs; grows with file count | Rejected: token waste |
| **Frontmatter domain tags** | Hugo/Jekyll taxonomies | Self-describing; consolidate-specs auto-manages; no extra files | Must update existing files once | **Selected: solves N×M problem with M self-descriptions** |

## Feasibility Assessment

### Technical Feasibility

All convention/business-rule files already have YAML frontmatter (the `title` field). Adding a `domains` field is a mechanical extension. The discovery instruction is plain text that can be embedded in prompt templates and SKILL.md files without code changes.

### Resource & Timeline

6 existing files need frontmatter updates (3 conventions + 3 business-rules). 3 prompt/skill files need the shared instruction. 1 skill (consolidate-specs) needs frontmatter management logic. All changes are small and well-bounded.

### Dependency Readiness

No external dependencies. All files are local and under version control.

## Scope

### In Scope

1. Define frontmatter `domains` schema for convention and business-rule files
2. Add `domains` frontmatter to all 6 existing files (3 conventions + 3 business-rules)
3. Create a shared discovery instruction snippet (reusable text block)
4. Update `forge-cli/pkg/prompt/data/implementation.md` — replace vague "read relevant" with discovery instruction
5. Update `forge-cli/pkg/prompt/data/fix.md` — same replacement
6. Update `plugins/forge/commands/fix-bug.md` — remove hardcoded keyword→filename mapping table, replace with discovery instruction
7. Update `plugins/forge/skills/consolidate-specs/SKILL.md` — write/maintain `domains` frontmatter when creating or updating files
8. Update `plugins/forge/hooks/guide.md` — document the discovery mechanism in the agent note about project knowledge

### Out of Scope

- Changes to `gen-test-cases-per-type-dispatch` proposal or its implementation (orthogonal concern — per-type instruction files declaring their own `conventions` dependencies is a separate pattern)
- Creating new convention or business-rule files
- Changes to eval skills or run-tasks
- Changes to `gen-test-scripts`, `gen-test-cases`, or `eval-test-cases` skills (they can adopt this pattern later)
- Automated domain extraction (consolidate-specs derives domains from spec content at write time, not a separate extraction step)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Domains tags too broad (every file matches every task) | M | M | consolidate-specs derives domains from spec ID and source keywords, not freeform tagging. Guidelines: 3-7 specific keywords per file |
| Domains tags too narrow (relevant files missed) | M | H | Backward compatibility: files without `domains` loaded unconditionally. Missing a match only means the agent doesn't read that file — same as current behavior when the hardcoded table doesn't cover the keyword |
| Agent misinterprets the matching instruction | M | M | Instruction is explicit: read frontmatter, match domains keywords, load on overlap. Same complexity as current "read relevant files" but with concrete criteria |
| consolidate-specs writes inaccurate domains | L | M | Domains are derived from the spec content itself (ID keywords, source keywords). Human reviews via existing consolidate-specs confirmation flow |

## Success Criteria

- [ ] All 6 existing convention/business-rule files have valid `domains` frontmatter with 3-7 specific keywords each
- [ ] `implementation.md`, `fix.md`, and `fix-bug.md` contain the shared discovery instruction and no hardcoded filename references
- [ ] `consolidate-specs` SKILL.md writes `domains` frontmatter when creating new files and updates it when drifting existing files
- [ ] For a task about "error handling", the discovery instruction loads `error-handling.md` and `error-reporting.md` but not `profile-system.md` or `quality-gate.md`
- [ ] A convention file without `domains` frontmatter is loaded unconditionally (backward compatibility)
- [ ] No keyword→filename mapping tables remain in any skill or prompt template

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
