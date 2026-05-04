# Design: consolidate-specs Knowledge Alignment

## Problem

consolidate-specs extracts business rules and technical specs from feature documents, but:
1. Target directories `docs/business-rules/` and `docs/conventions/` don't exist yet
2. Overlap with existing `docs/decisions/` and `docs/lessons/` has no detection or resolution mechanism
3. Extracted knowledge is passive files — agents never read them during task execution

## Decisions

| # | Decision Point | Choice |
|---|---------------|--------|
| 1 | Knowledge directories | Create `business-rules/` + `conventions/` as independent directories |
| 2 | Overlap with decisions/lessons | Write to target + immediately delete old entry |
| 3 | File format | One .md per domain, sections with BIZ-NNN/TECH-NNN |
| 4 | ID across features | Project-global renumbering during integrate |
| 5 | Global ID encoding | Prefix from target filename + sequence within file |
| 6 | Overlap detection | Filename matching for decisions, tag matching for lessons |
| 7 | Tags implementation | Only modify learn-lesson; decisions use filename matching |
| 8 | Cleanup timing | During integrate, immediately |
| 9 | Delete mechanism | Direct delete, no preservation |
| 10 | learn-lesson tags | User manually picks from fixed 8-category vocabulary |
| 11 | Knowledge activation | Embed read instructions in agent/command definitions |
| 12 | PHASE_SUMMARY vs project knowledge | Complementary, project knowledge first |

---

## 1. Directory Structure

```
docs/
  business-rules/       — Cross-feature business rules (by domain)
    auth.md              — Auth/permissions rules
    user.md              — User management rules
    ...
  conventions/          — Technical specs (coding standards, API conventions)
    api.md               — API design conventions
    error-handling.md    — Error handling patterns
    naming.md            — Naming conventions
    ...
  decisions/            — Technical decisions (/record-decision) — unchanged
  lessons/              — Lessons learned (/learn-lesson) — add tags frontmatter
```

---

## 2. File Format: business-rules/ and conventions/

One file per domain/topic. Entries grouped by feature section.

### business-rules/<domain>.md

```markdown
# Business Rules: <Domain>

## feature/<slug> (<date>)

### BIZ-<domain>-<NNN>: <Rule Title>

**Rule**: <concise rule statement>
**Context**: <why this rule exists>
**Scope**: [CROSS|LOCAL]
**Source**: <prd section reference>
```

### conventions/<topic>.md

```markdown
# Technical Conventions: <Topic>

## feature/<slug> (<date>)

### TECH-<topic>-<NNN>: <Spec Title>

**Requirement**: <concise requirement>
**Scope**: [CROSS|LOCAL]
**Source**: <design section reference>
```

### ID Encoding

- Prefix derived from target filename: `business-rules/auth.md` → prefix `auth`
- Sequence: file-internal, find max existing NNN + 1
- Examples: `BIZ-auth-001`, `TECH-api-003`
- Source traceability via `Source:` field: `Source: feature/auth-login BIZ-001`

---

## 3. Overlap Detection & Resolution

### During consolidate-specs Step 6 (review phase)

Before presenting to user, scan for related existing entries:

1. **Decisions**: match by filename → `docs/decisions/<topic>.md` (e.g., `auth` → `security.md`)
2. **Lessons**: match by tags frontmatter → grep `tags:` in `docs/lessons/*.md`

Display overlaps in review:
```
Related existing entries:
  - decisions/error-handling.md row "采用 AIError 结构体" (2026-04-30)
  - lessons/gotcha-error-handling.md tagged [error-handling, api]

Choose action for each:
  [skip] Keep both | [replace] Delete old, write new
```

### During consolidate-specs Step 7 (integrate)

For each user-approved replacement:

1. **Decisions**: delete the matching table row from `docs/decisions/<type>.md`
2. **Lessons**: delete `docs/lessons/<filename>.md`
3. Update `docs/decisions/manifest.md` counts
4. Write new entry to target file with project-global ID

### Review-choices.md update

Add "Related existing entries" section:
```markdown
## Related Existing Entries

- decisions/error-handling.md row "采用 AIError 结构体" → replaced by BIZ-auth-003
- lessons/gotcha-error-handling.md → deleted, superseded by TECH-error-005
```

---

## 4. learn-lesson Modifications

### New tags frontmatter

```markdown
---
tags: [error-handling, interface, security]
---

# gotcha-error-handling

## Problem
...
```

### Fixed tag vocabulary (aligned with decisions categories)

| Tag | Domain |
|-----|--------|
| `architecture` | System structure, layering |
| `interface` | API contracts, data shapes |
| `data-model` | Schema, indexing, soft-delete |
| `dependencies` | Library choices, version constraints |
| `error-handling` | Error types, status codes, propagation |
| `testing` | Test patterns, coverage, mocking |
| `security` | Auth, permissions, data protection |
| `local-dev-deployment` | Dev environment, tooling, deployment |

### New interaction round in learn-lesson

After "Key Takeaway" step, add:
```
Select tags (comma-separated, from fixed vocabulary):
  architecture, interface, data-model, dependencies, error-handling, testing, security, local-dev-deployment
```

---

## 5. Knowledge Activation in Agents

### Reading order (during task execution)

```
1. Project knowledge (docs/conventions/, docs/business-rules/) — domain constraints
2. PHASE_SUMMARY — feature-local context from previous tasks
3. Task definition — current task requirements
```

### Domain inference

Agent infers which knowledge files to read from task context:

| Task signals | Read |
|-------------|------|
| title/scope mentions "auth", "login", "permission" | `docs/business-rules/auth.md` |
| title/scope mentions "API", "endpoint", "route" | `docs/conventions/api.md` |
| title/scope mentions "error", "status code" | `docs/conventions/error-handling.md` |
| No matching file exists | Skip this step |

### Files to modify

Each file's task-preparation section adds:
```
Before implementing, read relevant project knowledge files:
- Infer relevant domains from task title, scope, and feature slug
- Read matching files from docs/business-rules/ and docs/conventions/
- If no matching file exists, skip this step
```

Affected files:
- `plugins/forge/agents/task-executor.md` — Step 1
- `plugins/forge/commands/execute-task.md` — pre-implementation
- `plugins/forge/agents/error-fixer.md` — before fix
- `plugins/forge/commands/fix-bug.md` — before fix

---

## 6. consolidate-specs SKILL.md Changes

### Step 5 update (preview)

After generating preview files, before presenting to user:
- Scan `docs/decisions/<topic>.md` for rows matching entry domain
- Scan `docs/lessons/*.md` for files with matching `tags:` frontmatter
- Collect matches as "Related existing entries"

### Step 6 update (review)

Display related entries alongside extracted specs:
```
Extracted specs for <slug>:

Business Rules:
  [CROSS] BIZ-001: <rule> → docs/business-rules/<domain>.md

Technical Specs:
  [CROSS] TECH-001: <spec> → docs/conventions/<topic>.md

Related existing entries (may overlap):
  ⚠ decisions/error-handling.md: "采用 AIError 结构体" — appears to overlap with TECH-001
  ⚠ lessons/gotcha-error-handling.md [error-handling] — may overlap with TECH-002

For each overlap, choose: [skip] keep both | [replace] delete old + write new
```

### Step 7 update (integrate)

For each approved item:
1. Determine target file and compute project-global ID (read file, find max NNN + 1)
2. If user chose "replace": delete old entry from decisions/lessons
3. Append new entry section to target file
4. Update decisions/manifest.md if any rows were deleted

### Step 8 update (marker)

`.integrated` marker gains overlap resolution record:
```yaml
feature: "<slug>"
integrated: "<date>"
biz_count: <N>
tech_count: <M>
replaced:
  - decisions/error-handling.md row "采用 AIError..." → BIZ-auth-003
  - lessons/gotcha-error-handling.md → TECH-error-005
```

---

## 7. guide.md Updates

### Directory conventions table — add rows

```
| `/consolidate-specs` | `docs/business-rules/<domain>.md`, `docs/conventions/<topic>.md` | Skill | Cross-feature knowledge extraction |
```

### Project-Level Documents — add entries

```
docs/
  business-rules/       — Cross-feature business rules (by domain, e.g. auth.md)
  conventions/          — Technical specs (coding standards, API conventions, naming rules)
```

### Rules — add entry

```
- `business-rules/` and `conventions/` — Populated by /consolidate-specs. Agent reads during task execution.
```

---

## Implementation Checklist

### Phase 1: Foundation
- [x] Update `consolidate-specs/SKILL.md` Steps 5-8 (overlap detection, global ID, replace logic)
- [x] Update `learn-lesson/SKILL.md` (add tags interaction round + frontmatter)
- [x] Update `learn-lesson/templates/` (add tags to template)

### Phase 2: Activation
- [x] Update `task-executor.md` Step 1 (add project knowledge read)
- [x] Update `execute-task.md` (add project knowledge read)
- [x] Update `error-fixer.md` (add project knowledge read)
- [x] Update `fix-bug.md` (add project knowledge read)

### Phase 3: Documentation
- [x] Update `guide.md` directory conventions
- [x] Update `SKILLS.md` if descriptions need alignment
