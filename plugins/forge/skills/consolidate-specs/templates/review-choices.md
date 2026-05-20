# review-choices.md Output Template

```markdown
---
feature: "<slug>"
reviewed: "<date>"
---

# Review Choices

## Approved for Integration

- BIZ-001 -> docs/business-rules/<domain>.md
- TECH-001 -> docs/conventions/<topic>.md

## Skipped

- (any items the user chose to skip)

## Related Existing Entries

- decisions/error-handling.md row "Adopt AIError struct" -> replaced by BIZ-auth-003
- lessons/gotcha-error-handling.md -> deleted, superseded by TECH-error-005
```
