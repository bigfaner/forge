---
created: 2026-05-17
author: faner
status: Draft
---

# Proposal: Knowledge Accumulation Unified Entry + Vocabulary

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

**Two-part solution: unified entry point + suggestive vocabulary.**

### Part 1: `/learn` — Unified Knowledge Skill

A new skill that serves as the single recommended entry point for all knowledge accumulation. The user describes what they learned/decided in free-form text; the skill identifies knowledge type(s), suggests classification from a built-in vocabulary, and writes to the appropriate directory.

```
/learn "发现 map 非线程安全导致 race condition，决定用 sync.Map"
  → Agent identifies: lesson + decision
  → Suggests domain: architecture (from vocabulary)
  → Writes to: docs/lessons/gotcha-race-condition.md + docs/decisions/architecture.md
  → User reviews and confirms

/learn "所有 API 必须有 /api 前缀"
  → Agent identifies: convention
  → Suggests domain: interface (from vocabulary)
  → Appends to: docs/conventions/api.md
  → User reviews and confirms
```

**Scope of `/learn`:**
- Single-record operations: one decision, one lesson, one convention entry, one business-rule entry
- Multi-type capture: a single input that produces entries in multiple directories
- Convention/business-rule appending to existing files (or creating new files)
- Bulk extraction from feature docs → delegates to `/consolidate-specs`

**Old skills:** `/record-decision` and `/learn-lesson` remain functional but are demoted to "low-level API" — still available for power users, but `/learn` is the recommended entry point. `/consolidate-specs` remains unchanged (it handles complex bulk extraction with drift detection).

### Part 2: Built-in Vocabulary

A reference file (`plugins/forge/references/shared/vocabulary.md`) containing suggested categories for knowledge classification. The vocabulary is presented as recommendations during the `/learn` flow — users can accept suggestions or type any custom value.

```yaml
types:
  - decision       # "Why we chose X"
  - lesson         # "What went wrong and how to fix it"
  - convention     # "Technical standard: always do X"
  - business-rule  # "Business constraint: X must satisfy Y"

domains:
  - architecture        # System structure, layering
  - interface           # API contracts, data shapes
  - data-model          # Schema, indexing, soft-delete
  - dependencies        # Library choices, version constraints
  - error-handling      # Error types, status codes, propagation
  - testing             # Test patterns, coverage, mocking
  - security            # Auth, permissions, data protection
  - local-dev-deployment # Dev environment, tooling, deployment
```

The domain vocabulary matches the existing 8-category system used by decisions and lessons, ensuring backward compatibility. Users can enter any custom domain (e.g., "concurrency", "performance") — the vocabulary is suggestive, not enforced.

### Part 3: Trigger Points

Context-aware suggestions at 3 natural workflow completion points, driving users to `/learn`:

| Trigger Point | Detection | Suggestion |
|---------------|-----------|------------|
| `run-tasks` completes all tasks | Architectural decisions, novel patterns, business rules in task outcomes | Suggest `/learn` with pre-filled summary |
| `fix-bug` completes a fix | Non-obvious root cause or notable debugging pattern | Suggest `/learn` with root cause summary |
| `write-prd` / `tech-design` completes | New business rules, architecture decisions | Suggest `/learn` with decision/rule summary |

`/quick` is covered by `run-tasks` since quick calls run-tasks internally.

### Innovation Highlights

The design borrows from **unified search bars** (Google, VS Code Command Palette) — instead of requiring users to know the specific command, a single entry point interprets intent and routes to the right destination. The vocabulary acts like **autocomplete suggestions** — it narrows the space without constraining it. This is the opposite of the current design where users must navigate a menu of specialized commands.

## Requirements Analysis

### Key Scenarios

- **Multi-type knowledge**: User runs `/learn "race condition from non-thread-safe map, decided to use sync.Map"`. Agent identifies both lesson and decision. Suggests writing to both `docs/lessons/` and `docs/decisions/`. User confirms both entries.
- **Single convention**: User runs `/learn "all API endpoints must have /api prefix"`. Agent identifies as convention, suggests `docs/conventions/api.md`. Appends rule with project-global ID.
- **Custom domain**: User runs `/learn "websocket connections must have heartbeat every 30s"`. Agent suggests domain `interface` from vocabulary. User types `websocket` instead. Accepted without error.
- **Trigger from run-tasks**: After completing tasks that created a new authentication module, run-tasks detects architectural knowledge produced. Suggests `/learn` with summary: "Created new auth module with JWT-based session management."
- **Trigger from fix-bug**: After fixing a deadlock caused by lock ordering, fix-bug detects notable root cause. Suggests `/learn` with summary: "Deadlock from inconsistent lock acquisition order."
- **Nothing notable**: After routine config changes, run-tasks detects nothing worth capturing. Silent — no suggestion.
- **Bulk extraction**: User has a completed feature with PRD + design docs. Wants to extract all business rules at once. `/learn` detects this is a bulk operation and suggests `/consolidate-specs` instead.

### Non-Functional Requirements

- **Vocabulary non-enforcement**: Any value accepted for type and domain, even if not in vocabulary. Zero friction for custom entries.
- **Backward compatibility**: Existing files in `docs/decisions/`, `docs/lessons/`, `docs/conventions/`, `docs/business-rules/` are not modified or migrated. Old skills continue to work.
- **Zero noise from triggers**: Trigger points are silent when no notable knowledge is detected. False-positive rate < 30%.

### Constraints & Dependencies

- Knowledge directory formats (decision rows, lesson files, convention entries) must remain compatible with existing skills
- `/learn` reuses existing format specifications from `decision-logging.md` and `learn-lesson` templates
- No code changes — all changes are prompt-level (SKILL.md, command files, reference files)

## Alternatives & Industry Benchmarking

### Industry Solutions

- **Unified search/command palette** (VS Code, Raycast): Single entry point that interprets intent and routes to specific actions. Users don't need to know command names. This is the UX model for `/learn`.
- **Tag-based knowledge management** (Notion, Obsidian): Knowledge entries tagged with flexible categories. Tags are suggestive (common tags shown) but not enforced. This is the vocabulary model.
- **Git commit hooks** (husky, pre-commit): Trigger points that run code at natural workflow boundaries. The trigger points in this proposal serve the same purpose — but they're prompt-based suggestions, not programmatic gates.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | Three separate skills, fragmented knowledge, low accumulation rate | Rejected |
| Merge all directories | Knowledge graph model | Single location for all knowledge | Loses semantic distinctions (decision vs lesson vs rule); breaks existing format consumers | Rejected: too disruptive |
| Wrapper skill over old skills | Adapter pattern | Minimal change to existing skills | Still depends on three separate skill implementations; wrapper adds indirection | Rejected: doesn't simplify the model |
| **Unified entry + suggestive vocabulary** | VS Code Command Palette + Obsidian tags | Single entry point; flexible classification; backward compatible; natural UX | Requires new skill implementation; vocabulary maintenance | **Selected: combines unified entry UX with flexible classification, minimal disruption** |

## Feasibility Assessment

### Technical Feasibility

All changes are prompt-level. `/learn` is a new skill that reads the vocabulary reference and writes to existing directory formats. Trigger points are additions to existing SKILL.md/command files. No code changes required.

### Resource & Timeline

| Deliverable | Type | Complexity |
|-------------|------|-----------|
| `/learn` SKILL.md + templates | New skill | Medium (multi-format output routing) |
| `references/shared/vocabulary.md` | New reference | Low (static file) |
| `hooks/guide.md` update | Edit existing | Low (add `/learn` section, demote old skills) |
| `commands/run-tasks.md` trigger | Edit existing | Low (add knowledge review step) |
| `commands/fix-bug.md` trigger | Edit existing | Low (add knowledge review step) |
| `skills/write-prd/SKILL.md` trigger | Edit existing | Low (add knowledge review step) |
| `skills/tech-design/SKILL.md` trigger | Edit existing | Low (add knowledge review step) |

### Dependency Readiness

- Existing format specs (`decision-logging.md`, `learn-lesson` template) are available for reuse
- Knowledge directory structures exist and are stable
- No external dependencies

## Scope

### In Scope

1. Create `/learn` skill (`plugins/forge/skills/learn/SKILL.md`) — unified knowledge accumulation entry point
2. Create `/learn` skill template (`plugins/forge/skills/learn/templates/`) — if needed for multi-format output
3. Create vocabulary reference (`plugins/forge/references/shared/vocabulary.md`) — built-in suggestive vocabulary
4. Update `plugins/forge/hooks/guide.md` — add `/learn` section, demote old skills to "advanced"
5. Add knowledge review trigger to `plugins/forge/commands/run-tasks.md`
6. Add knowledge review trigger to `plugins/forge/commands/fix-bug.md`
7. Add knowledge review trigger to `plugins/forge/skills/write-prd/SKILL.md`
8. Add knowledge review trigger to `plugins/forge/skills/tech-design/SKILL.md`

### Out of Scope

- Changes to `/record-decision`, `/learn-lesson`, or `/consolidate-specs` skills — they remain functional as-is
- Directory structure changes — all four knowledge directories keep their current formats
- Migration of existing knowledge files
- Changes to knowledge discovery mechanism (domains frontmatter)
- Programmatic enforcement (hooks) of directory conventions

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `/learn` misidentifies knowledge type | M | M | Present classification to user for confirmation before writing. User can always correct. The existing trigger vocabulary (decision, lesson, convention, business-rule) has clear semantic boundaries. |
| Trigger points produce too many false suggestions | M | M | Keep heuristics conservative: only suggest for genuinely non-obvious knowledge. If acceptance rate < 30%, tighten detection criteria. |
| Users still prefer old skills | L | L | Old skills remain functional. `/learn` is recommended but not forced. Gradual adoption is fine. |
| Vocabulary becomes stale as project evolves | L | L | Vocabulary is suggestive — users naturally use custom terms when built-in terms don't fit. The vocabulary file can be updated in future Forge releases. |

## Success Criteria

- [ ] `/learn` correctly identifies knowledge type(s) from free-form input in 4+ test scenarios (multi-type, single convention, custom domain, bulk delegation)
- [ ] `/learn` writes to all 4 knowledge directories (`docs/decisions/`, `docs/lessons/`, `docs/conventions/`, `docs/business-rules/`) using their existing formats
- [ ] `/learn` accepts custom vocabulary values (domains, types) not in the built-in vocabulary without error
- [ ] `/learn` suggests `/consolidate-specs` when the user's input describes a bulk extraction need
- [ ] Trigger at `run-tasks` completion suggests `/learn` when tasks produced architectural decisions or novel patterns
- [ ] Trigger at `fix-bug` completion suggests `/learn` when the root cause was non-obvious
- [ ] Triggers are silent when no notable knowledge was produced (routine tasks, trivial fixes)
- [ ] `guide.md` references `/learn` as the primary knowledge accumulation entry point
- [ ] Old skills (`/record-decision`, `/learn-lesson`) remain functional and documented as "advanced" alternatives
- [ ] All modified files pass `eval-forge` structural consistency check

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
