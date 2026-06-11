---
id: "4"
title: "Wire BodyContext through BuildIndex with proposal/PRD data extraction"
priority: "P1"
estimated_time: "1-2h"
dependencies: ["2", "3"]
scope: "backend"
breaking: true
type: "coding.feature"
mainSession: false
---

# 4: Wire BodyContext through BuildIndex with proposal/PRD data extraction

## Description

Update `BuildIndex()` in `build.go` to populate `BodyContext` with planning-time data extracted from the proposal or PRD, and pass it through to `GenerateTestTaskMD()`. This connects the data source (proposal/PRD) to the rendering pipeline (renderBody) so that template placeholders are filled with real feature context.

## Reference Files
- `docs/proposals/auto-task-main-instructions/proposal.md` — "Changes to BuildIndex()" section
- `forge-cli/pkg/task/build.go` — BuildIndex() function (lines 268, 301 call GenerateTestTaskMD)
- `forge-cli/pkg/task/autogen.go` — GenerateTestTaskMD() and renderBody()

## Acceptance Criteria

- [ ] `BuildIndex()` extracts `Scope` from proposal/PRD "## Scope > ### In Scope" section
- [ ] `BuildIndex()` extracts `SuccessCriteria` from proposal/PRD "## Success Criteria" section
- [ ] `BuildIndex()` extracts `AcceptanceCriteria` from PRD "## Acceptance Criteria" section (breakdown mode only)
- [ ] `FeatureSlug` and `Mode` are populated from existing BuildIndex data
- [ ] `Interfaces` populated from `.forge/config.yaml` (already available in BuildIndex via `forgeconfig.ReadInterfaces`)
- [ ] Both `GenerateTestTaskMD()` call sites in `build.go` pass populated `BodyContext`
- [ ] Existing tests pass (backward compatible — empty BodyContext produces same output as before)

## Hard Rules

- MUST NOT do directory scanning, spec file listing, or code path derivation in BodyContext population
- MUST only extract data that exists at planning time (proposal/PRD content, config)
- Runtime-only data (test results, generated files) MUST NOT be in BodyContext
- MUST handle missing proposal/PRD gracefully — empty BodyContext fields are valid

## Implementation Notes

### Data extraction approach:

1. **FeatureSlug + Mode**: Already available in BuildIndex (`opts.FeatureSlug`, `detectMode()`)
2. **Scope**: Parse proposal/PRD file, extract bullet list under "## Scope > ### In Scope"
3. **SuccessCriteria**: Parse proposal/PRD file, extract checked/unchecked items under "## Success Criteria"
4. **AcceptanceCriteria**: Parse PRD only, extract items under "## Acceptance Criteria" (empty in quick mode)
5. **Interfaces**: Already available via `forgeconfig.ReadInterfaces()` in BuildIndex

### Proposal/PRD file path resolution:

```go
// Already set by setFeatureMetadata()
proposalPath := filepath.Join(opts.ProjectRoot, index.Proposal) // quick mode
prdPath := filepath.Join(opts.ProjectRoot, "docs", "features", opts.FeatureSlug, "prd", "prd-spec.md") // breakdown mode
```

### BuildIndex call sites to update (build.go):

Line 268: `evalContent, err := GenerateTestTaskMD(evalTask, opts.FeatureSlug)`
Line 301: `content, genErr := GenerateTestTaskMD(td, opts.FeatureSlug)`

Both need to change to pass BodyContext instead of string.
