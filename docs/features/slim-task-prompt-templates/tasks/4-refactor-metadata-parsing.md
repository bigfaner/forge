---
id: "4"
title: "Refactor Go metadata parsing for grouped frontmatter"
priority: "P0"
estimated_time: "2h"
complexity: "high"
dependencies: [2, 3]
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.enhancement"
mainSession: false
---

# 4: Refactor Go metadata parsing for grouped frontmatter

## Description
Extend the Go metadata parsing infrastructure to support semantic grouping of template variables. The current flat `variables` list in metadata frontmatter becomes a structured format with `identity`, `context`, `conditional`, and `variables` groups.

Changes:
1. Extend `TemplateMetadata` struct with `Identity`, `Context`, `Conditional` map fields and `AllFields()` method
2. Update `parseMetadataFrontmatter()` to parse grouped YAML (recommend `gopkg.in/yaml.v3` for robust parsing)
3. Extend `validateMetadataVariables()` to validate fields across all groups
4. Maintain backward compatibility — old format (variables only) continues to parse correctly
5. Add comprehensive unit tests

## Reference Files
- forge-cli/pkg/prompt/metadata.go: TemplateMetadata struct + parseMetadataFrontmatter + validateMetadataVariables (source: proposal.md#In-Scope-Frontmatter-重构)
- forge-cli/pkg/prompt/metadata_test.go: Add grouped parsing, validation, and backward compatibility tests (source: proposal.md#In-Scope-Frontmatter-重构)
- forge-cli/pkg/prompt/prompt.go: Embed FS template loading — reference only, no changes to Synthesize() main logic (source: proposal.md#Constraints-&-Dependencies)

## Acceptance Criteria
- [ ] SC-FM-2: parseMetadataFrontmatter backward compatible — old format (flat variables list, no groups) parses correctly, Variables field matches old parser output
- [ ] SC-FM-3: validateMetadataVariables validates all grouped fields (Identity/Context/Conditional keys + Variables) against corresponding Go struct via reflect
- [ ] AllFields() method returns union of Identity + Context + Conditional keys + Variables list for backward compatibility
- [ ] All unit tests pass: grouped parsing, grouped validation, backward compatibility, edge cases

## Hard Rules
- **parseMetadataFrontmatter backward compatibility**: Templates without frontmatter or with old format must continue to parse without errors
- **Field name PascalCase alignment**: Group field names must match Go struct field names exactly (reflect.FieldByName is case-sensitive)
- **rendered frontmatter separation**: This task only modifies metadata frontmatter parsing (first `---` block). Rendered frontmatter (second `---` block in task/record templates) is parsed by `task/frontmatter.go` and must not be affected

## Implementation Notes

### Test Impact
- Affected test suite(s): forge-cli/pkg/prompt/
- Expected fixture changes: test template files with new grouped frontmatter format
- Risk level: medium

- Recommend using `gopkg.in/yaml.v3` for YAML parsing instead of extending hand-written line-level parser — eliminates indent state tracking bugs
- TemplateMetadata struct per proposal:
  ```go
  type TemplateMetadata struct {
      Type        string
      Category    string
      Identity    map[string]bool
      Context     map[string]bool
      Conditional map[string]bool
      Variables   []string
  }
  ```
- Validation struct mapping: prompt → `promptTemplateData`, task → `TemplateData`, record → `RecordTemplateData`
