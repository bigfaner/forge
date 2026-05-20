# Self-Check Rules

Verification checklist for PRD quality before presenting to the user.

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
