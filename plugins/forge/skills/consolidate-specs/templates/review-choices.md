# review-choices.md Output Template

```markdown
---
feature: "{{SLUG}}"
reviewed: "{{DATE}}"
---

# Review Choices

## Approved for Integration

- BIZ-001 -> docs/business-rules/{{DOMAIN}}.md
- TECH-001 -> docs/conventions/{{TOPIC}}.md

## Skipped

- (any items the user chose to skip)

## Related Existing Entries

- decisions/error-handling.md row "Adopt AIError struct" -> replaced by BIZ-auth-003
- lessons/gotcha-error-handling.md -> deleted, superseded by TECH-error-005
```
