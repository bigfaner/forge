# Self-Check Rules

Verification checklist for PRD quality before presenting to the user.

**Intent gate**: The checks below branch on the PRD's `intent` field (read from `docs/proposals/<slug>/proposal.md` frontmatter). Valid values: `new-feature`, `enhancement`, `refactor`, `cleanup`, `fix`, `doc`. Default is `new-feature` if missing or unrecognized.

## new-feature intent

| Check | What to verify |
|-------|----------------|
| Background completeness | Reason + target users + stakeholders all present and specific |
| Goals quantified | At least one numeric target (% , count, time) |
| Flow diagram | Mermaid flowchart with decision points (diamond nodes) and at least one error/exception branch |
| Functional specs | prd-spec.md references prd-ui-functions.md; prd-ui-functions.md tables filled — no placeholder rows |
| User stories | One story per user role, each with Given/When/Then AC |
| Scope consistency | In-scope items match what's described in Functional Specs and user stories |
| No vague language | No "better", "faster", "improved" without quantification |
| Placement completeness | Every UI Function has a Placement section with Mode and target |
| Placement consistency | existing-page routes exist in sitemap.json (only when project has `web` surface — check via `forge surfaces --json`) |
| Sitemap availability | If project has `web` surface but sitemap.json not found, warn: "Sitemap unavailable — existing-page routes cannot be validated. Run /gen-web-sitemap." If no `web` surface, skip this check entirely |
| Page Composition valid | Page Composition table lists all pages with correct UI Function references |
| db-schema filled | db-schema frontmatter is "yes" or "no" (not empty) |

## enhancement intent

When `intent` is `enhancement`, the PRD uses Simplified format (Background + Goals + Test Pipeline). Apply the checks below:

| Check | What to verify |
|-------|----------------|
| Background completeness | Reason + what is being improved, clearly stated |
| Goals quantified | At least one numeric target or measurable improvement criterion |
| Test Pipeline | prd-spec.md includes a Test Pipeline section ensuring enhancement has test coverage |
| No vague language | No "better", "faster", "improved" without quantification |
| Scope consistency | In-scope items match the described enhancement boundaries |
| db-schema filled | db-schema frontmatter is "yes" or "no" (not empty) |
| Override signals | If any override signal was triggered, verify `<!-- Override: ... -->` comment is present |

**Skipped for enhancement** (these artifacts are not generated):
- User stories (Given/When/Then) — existing user base, no new user flows
- Flow diagram (Mermaid) — enhancement does not introduce new flows
- Functional specs (prd-ui-functions.md reference)
- Placement completeness / consistency
- Page Composition validity

## refactor / cleanup / fix intent

When `intent` is `refactor`, `cleanup`, or `fix`, the PRD uses spec-only format — no user stories, no UI functions, no flow diagrams (unless the refactoring changes an external flow). Apply only the checks below:

| Check | What to verify |
|-------|----------------|
| Change Scope present | `prd-spec.md` contains a "Change Scope" (变更范围) section listing concrete modules, files, or packages |
| Constraints present | `prd-spec.md` contains a "Constraints" (约束条件) section with behavioral invariants that must be preserved |
| Verification Criteria present | `prd-spec.md` contains a "Verification Criteria" (验证标准) section with regression acceptance criteria |
| Background completeness | Reason + target users + stakeholders all present and specific |
| No vague language | No "better", "faster", "improved" without quantification |
| Scope consistency | Change Scope items match the described refactoring boundaries |
| db-schema filled | db-schema frontmatter is "yes" or "no" (not empty) |
| Override signals | If any override signal was triggered, verify `<!-- Override: ... -->` comment is present |

**Skipped for refactor/cleanup/fix** (these artifacts are not generated):
- User stories (Given/When/Then)
- Flow diagram (Mermaid)
- Functional specs (prd-ui-functions.md reference)
- Placement completeness / consistency
- Page Composition validity

## doc intent

When `intent` is `doc`, the PRD uses Minimal format (title + goals + scope only). Apply only the checks below:

| Check | What to verify |
|-------|----------------|
| Title present | One-sentence description of the documentation change target and purpose |
| Goals present | Lists specific documentation files to update/create and expected changes |
| Scope present | Clear boundaries of what documentation is in-scope and out-of-scope |
| Scope consistency | Listed files match the described documentation change boundaries |
| db-schema filled | db-schema frontmatter is "no" (doc changes should not involve DB schema) |

**Skipped for doc** (these artifacts are not generated):
- User stories (Given/When/Then)
- Flow diagram (Mermaid)
- Functional specs (prd-ui-functions.md reference)
- Placement completeness / consistency
- Page Composition validity
- Override signal checks (doc intent has no overridable pipeline steps — signals are no-op)
