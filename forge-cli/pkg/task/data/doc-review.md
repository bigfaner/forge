Review documentation quality for the {{FEATURE_SLUG}} feature ({{MODE}} mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

{{DOC_TASK_AC}}

## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/{{FEATURE_SLUG}}/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/{{FEATURE_SLUG}}/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
