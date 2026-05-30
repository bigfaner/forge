# Self-Check Rules

Verification checklist for PRD quality before presenting to the user.

**Intent gate**: The checks below branch on the PRD's `intent` field (read from `docs/proposals/<slug>/proposal.md` frontmatter). Default is `new-feature` if missing.

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
| Placement consistency | existing-page routes exist in sitemap.json (if sitemap available) |
| Sitemap availability | If sitemap.json not found, warn: "Sitemap unavailable — existing-page routes cannot be validated. Run /gen-sitemap." |
| Page Composition valid | Page Composition table lists all pages with correct UI Function references |
| db-schema filled | db-schema frontmatter is "yes" or "no" (not empty) |

## refactor / cleanup intent

When `intent` is `refactor` or `cleanup`, the PRD uses spec-only format — no user stories, no UI functions, no flow diagrams (unless the refactoring changes an external flow). Apply only the checks below:

| Check | What to verify |
|-------|----------------|
| Change Scope present | `prd-spec.md` contains a "Change Scope" (变更范围) section listing concrete modules, files, or packages |
| Constraints present | `prd-spec.md` contains a "Constraints" (约束条件) section with behavioral invariants that must be preserved |
| Verification Criteria present | `prd-spec.md` contains a "Verification Criteria" (验证标准) section with regression acceptance criteria |
| Background completeness | Reason + target users + stakeholders all present and specific |
| No vague language | No "better", "faster", "improved" without quantification |
| Scope consistency | Change Scope items match the described refactoring boundaries |
| db-schema filled | db-schema frontmatter is "yes" or "no" (not empty) |

**Skipped for refactor/cleanup** (these artifacts are not generated):
- User stories (Given/When/Then)
- Flow diagram (Mermaid)
- Functional specs (prd-ui-functions.md reference)
- Placement completeness / consistency
- Page Composition validity
