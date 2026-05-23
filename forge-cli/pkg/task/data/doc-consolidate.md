Extract and consolidate business rules and tech specs from the {{FEATURE_SLUG}} feature.

## Feature Context
- Scope: {{SCOPE}}

## Discovery Strategy
1. Scan docs/features/{{FEATURE_SLUG}}/ for all feature documents (PRD, design, task records)
2. Scan docs/proposals/{{FEATURE_SLUG}}/ for proposal
3. Extract rules and specs from discovered documents
4. Compare against existing specs in docs/business-rules/ and docs/conventions/

Run in non-interactive mode: auto-integrate all CROSS items. Commit with [auto-specs] tag.
