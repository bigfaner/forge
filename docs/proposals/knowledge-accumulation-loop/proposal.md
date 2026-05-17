---
created: 2026-05-17
author: faner
status: Draft
---

# Proposal: Knowledge Accumulation Unified Entry + Auto-Extract

## Problem

Forge's knowledge accumulation relies on three separate skills (`/record-decision`, `/learn-lesson`, `/consolidate-specs`) that capture different facets of the same events. Users must (1) decide which skill to invoke, (2) invoke each one separately, and (3) manually ensure consistency across directories. The result: knowledge directories stay empty because the cognitive overhead exceeds the perceived benefit.

### Evidence

- Three skills share the same 8-category vocabulary (architecture, interface, data-model, etc.) but classify independently — decisions use type files, lessons use tags, conventions use domains frontmatter
- `/consolidate-specs` already has explicit overlap detection with decisions and lessons (Step 5-6), proving the system acknowledges knowledge fragmentation across directories
- A single debugging session ("race condition due to non-thread-safe map, decided to use sync.Map") produces a lesson AND a decision — requiring two separate skill invocations
- 40+ lesson files exist in `docs/lessons/` but `docs/decisions/` has sparse entries, suggesting users gravitate toward one skill and skip the others
- Knowledge discovery (`domains` frontmatter) only covers `docs/conventions/` and `docs/business-rules/`, not `docs/decisions/` or `docs/lessons/`

### Urgency

Knowledge discovery was implemented but has nothing to discover. Each new feature that completes without knowledge accumulation widens the gap between what the system can discover and what actually exists. The fragmentation gets worse with each feature, not better.

## Proposed Solution

**Three-part solution: unified manual entry + auto-extract triggers + vocabulary归纳.**

### Part 1: `/learn` — Unified Manual Entry

A new skill that absorbs `/record-decision` and `/learn-lesson` into a single on-demand entry point. Used for knowledge that bypasses the pipeline — ad-hoc debugging insights, spontaneous realizations, mid-task discoveries.

**Two input modes:**

```
# Mode 1: Interactive — agent asks, then classifies
/learn
  → Agent: "What did you learn or decide?"
  → User describes...
  → Agent classifies, writes, reports

# Mode 2: Direct input — skip the first question
/learn "race condition from non-thread-safe map, decided to use sync.Map"
  → Agent identifies: lesson + decision
  → Writes to: docs/lessons/gotcha-race-condition.md + docs/decisions/architecture.md
  → Report includes entries for review
```

**Write-first, review-after:** Agent classifies and writes entries immediately, then includes all written entries in the final report for user review. This avoids interrupting the task flow — the user sees what was written and can correct anything after the fact.

**Scope of `/learn`:**
- Single-record operations: one decision, one lesson, one convention entry, one business-rule entry
- Multi-type capture: a single input that produces entries in multiple directories
- Convention/business-rule appending to existing files (or creating new files)
- Bulk extraction from feature docs → delegates to `/consolidate-specs`

**Removed skills:** `/record-decision` and `/learn-lesson` are deleted. Their functionality is fully absorbed by `/learn`.

### Part 2: Auto-Extract Triggers

Instead of suggesting `/learn`, triggers automatically extract knowledge from the current feature's artifacts and present for confirmation.

**Flow:**

```
Feature completes (run-tasks / fix-bug / write-prd / tech-design)
  → Scan feature's PRD + tech-design + task outcomes
  → Identify new knowledge (decisions, lessons, conventions, business rules)
  → Extract & summarize
  → Report to user for review
  → Write to knowledge dirs on confirmation
```

**Trigger points:**

| Trigger Point | What to scan | Knowledge types to look for |
|---------------|-------------|---------------------------|
| `run-tasks` completes all tasks | Task outcomes, code changes, manifest | Architectural decisions, novel patterns, gotchas, business rules |
| `fix-bug` completes a fix | Root cause analysis, fix approach | Non-obvious root causes, debugging patterns |
| `write-prd` completes | PRD content | New business rules, user-facing constraints |
| `tech-design` completes | Design document | Architecture decisions, dependency choices, data model decisions |

**Key behaviors:**
- Triggers are **silent when no notable knowledge is detected** — routine config changes, trivial fixes produce no output
- Extracted knowledge is **presented for user confirmation** before writing — the auto-extract is a draft, not a final action
- The extraction logic is a shared routine (prompt section) reused across all 4 trigger points
- `/consolidate-specs` vocabulary is used during extraction to suggest classifications

### Part 3: Auto-Generated Vocabulary via `/consolidate-specs`

No standalone vocabulary file. `/consolidate-specs` automatically归纳 and maintains a vocabulary index from existing knowledge files during its drift-detection pass. The vocabulary is regenerated each run, staying in sync with actual content.

Both `/learn` and the auto-extract triggers read this vocabulary at runtime to suggest classifications. No manual maintenance needed.

## Requirements Analysis

### Key Scenarios

- **Auto-extract after run-tasks**: Feature implementing auth module completes. Trigger scans task outcomes, identifies: "JWT-based session management" as architectural decision, "token expiry must be configurable" as business rule. Presents both for user review. User confirms. Written to `docs/decisions/` and `docs/business-rules/`.
- **Auto-extract after fix-bug**: Fix for deadlock from inconsistent lock ordering. Trigger extracts root cause as a lesson. Presents for review. User confirms. Written to `docs/lessons/`.
- **Auto-extract — nothing notable**: Routine config file changes. Trigger scans, finds no notable knowledge. Silent — no output.
- **Manual /learn**: Developer realizes mid-task that their ORM pattern has a gotcha. Runs `/learn "GORM hooks fire in creation order, not dependency order"`. Agent identifies as lesson, writes to `docs/lessons/`. Report shows entry for review.
- **Multi-type manual /learn**: `/learn "race condition from non-thread-safe map, decided to use sync.Map"`. Agent identifies both lesson and decision. Writes to both directories. Report shows both entries.
- **Custom domain**: User types a domain not in auto-vocabulary. Accepted without error.
- **Bulk extraction**: `/learn` detects bulk need, delegates to `/consolidate-specs`.

### Non-Functional Requirements

- **Vocabulary non-enforcement**: Any value accepted for type and domain. Zero friction for custom entries.
- **Backward compatibility**: Existing files in knowledge directories are not modified or migrated.
- **Zero noise from triggers**: False-positive rate < 30%. Silent when nothing notable.
- **Extraction consistency**: Auto-extract and `/learn` produce the same file formats.

### Constraints & Dependencies

- Knowledge directory formats must remain compatible with `/consolidate-specs`
- `/learn` reuses existing format specifications from the old `/record-decision` and `/learn-lesson` templates
- No code changes — all changes are prompt-level (SKILL.md, command files, reference files)
- `/consolidate-specs` needs a vocabulary generation step added to its drift-detection pass

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | Zero cost | Three separate skills, fragmented knowledge, low accumulation rate | Rejected |
| Wrapper skill over old skills | Minimal change | Still depends on three implementations; wrapper adds indirection | Rejected |
| Keep old skills + add /learn | Backward compatible | Four skills for same job; confusing for users | Rejected |
| **Unified /learn + auto-extract triggers + auto-vocabulary** | Single manual entry; auto-capture at pipeline boundaries; no manual vocab maintenance; removes old skills cleanly | More moving parts (triggers + extraction routine) | **Selected** |

## Feasibility Assessment

### Technical Feasibility

All changes are prompt-level. The shared extraction routine is a reusable prompt section. Trigger points append this routine to existing skill/command files. No code changes required.

### Resource & Timeline

| Deliverable | Type | Complexity |
|-------------|------|-----------|
| `/learn` SKILL.md + templates | New skill (absorbs 2 old skills) | Medium (multi-format output routing) |
| Delete `/record-decision` skill | Remove | Low |
| Delete `/learn-lesson` skill | Remove | Low |
| Shared extraction routine (prompt section) | New shared reference | Medium (knowledge identification + extraction logic) |
| Update `/consolidate-specs` — add vocabulary generation | Edit existing | Medium (归纳 logic) |
| `hooks/guide.md` update | Edit existing | Low |
| `commands/run-tasks.md` trigger | Edit existing | Low (add extraction step) |
| `commands/fix-bug.md` trigger | Edit existing | Low |
| `skills/write-prd/SKILL.md` trigger | Edit existing | Low |
| `skills/tech-design/SKILL.md` trigger | Edit existing | Low |

### Dependency Readiness

- Existing format specs from old skills are available for reuse
- Knowledge directory structures exist and are stable
- No external dependencies

## Scope

### In Scope

1. Create `/learn` skill (`plugins/forge/skills/learn/SKILL.md`) — unified manual knowledge entry, absorbs `/record-decision` and `/learn-lesson`
2. Create `/learn` templates if needed (`plugins/forge/skills/learn/templates/`) — merged from old skills
3. Delete `/record-decision` skill (`plugins/forge/skills/record-decision/`)
4. Delete `/learn-lesson` skill (`plugins/forge/skills/learn-lesson/`)
5. Create shared extraction routine (`plugins/forge/references/shared/knowledge-extraction.md`) — reusable prompt section for auto-extract triggers
6. Update `/consolidate-specs` — add auto-vocabulary generation step
7. Update `plugins/forge/hooks/guide.md` — replace old skills with `/learn`, document auto-extract flow
8. Add auto-extract trigger to `plugins/forge/commands/run-tasks.md`
9. Add auto-extract trigger to `plugins/forge/commands/fix-bug.md`
10. Add auto-extract trigger to `plugins/forge/skills/write-prd/SKILL.md`
11. Add auto-extract trigger to `plugins/forge/skills/tech-design/SKILL.md`

### Out of Scope

- Changes to `/consolidate-specs` beyond vocabulary generation step
- Directory structure changes — all four knowledge directories keep their current formats
- Migration of existing knowledge files
- Changes to knowledge discovery mechanism (domains frontmatter)
- Programmatic enforcement (hooks) of directory conventions

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Auto-extract produces false positives (noise) | M | M | Conservative heuristics + user confirmation gate. If acceptance rate < 30%, tighten detection. |
| `/learn` misidentifies knowledge type | M | M | Write-first-then-report: user sees all entries in final report and can correct. |
| Shared extraction routine diverges across trigger points | L | M | Single shared reference file (`knowledge-extraction.md`) included by all triggers. |
| Auto-vocabulary becomes stale between `/consolidate-specs` runs | L | L | `/learn` works without vocabulary — falls back to unassisted classification. Vocabulary is suggestive, not required. |
| Breaking change for users who used old skills | L | L | `/learn` covers all old functionality. Guide.md documents the migration. Old formats preserved in `/learn`. |

## Success Criteria

- [ ] `/learn` correctly identifies knowledge type(s) from free-form input in 4+ test scenarios (multi-type, single convention, custom domain, bulk delegation)
- [ ] `/learn` writes to all 4 knowledge directories using their existing formats
- [ ] `/learn` accepts custom vocabulary values without error
- [ ] `/learn` suggests `/consolidate-specs` for bulk extraction needs
- [ ] `/learn` works in both interactive mode (no args) and direct-input mode (with args)
- [ ] `/learn` final report includes all written entries for user review
- [ ] Old skills (`/record-decision`, `/learn-lesson`) are fully removed
- [ ] `/consolidate-specs` generates vocabulary index from existing knowledge files
- [ ] Auto-extract trigger at `run-tasks` completion identifies and extracts notable knowledge from task outcomes
- [ ] Auto-extract trigger at `fix-bug` completion identifies non-obvious root causes
- [ ] Auto-extract trigger at `write-prd` completion identifies new business rules
- [ ] Auto-extract trigger at `tech-design` completion identifies architecture decisions
- [ ] Triggers are silent when no notable knowledge was produced
- [ ] `guide.md` references `/learn` as the manual entry point and documents auto-extract flow
- [ ] All modified files pass `eval-forge` structural consistency check

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
