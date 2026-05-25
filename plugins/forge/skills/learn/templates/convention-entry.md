# Convention / Business-Rule Entry Template

Used by `/learn` when writing to `docs/conventions/` or `docs/business-rules/`.

## Entry Format

Append to the target file using the project-global ID encoding:

### Convention Entry (docs/conventions/{{TOPIC}}.md)

```markdown
### TECH-{{TOPIC}}-{{NNN}}: {{SPEC_TITLE}}

**Requirement**: {{CONCISE_REQUIREMENT}}
**Scope**: [CROSS]
**Source**: /learn entry {{DATE}}

{{IMPLEMENTATION_DETAILS}}
```

### Business-Rule Entry (docs/business-rules/{{DOMAIN}}.md)

```markdown
### BIZ-{{DOMAIN}}-{{NNN}}: {{RULE_TITLE}}

**Rule**: {{CONCISE_RULE_STATEMENT}}
**Context**: {{WHY_THIS_RULE_EXISTS}}
**Scope**: [CROSS]
**Source**: /learn entry {{DATE}}

{{ADDITIONAL_DETAILS}}
```

## Project-Global ID Encoding

- **Prefix**: derived from target filename (e.g., `business-rules/auth.md` -> `auth`, `conventions/api.md` -> `api`)
- **Sequence**: file-internal -- find max existing NNN in the target file + 1
- **Format**: `BIZ-{{DOMAIN}}-{{NNN}}` for business rules, `TECH-{{TOPIC}}-{{NNN}}` for tech specs
- **Examples**: `BIZ-auth-001`, `TECH-api-003`

## New File Frontmatter

When creating a new file, include YAML frontmatter:

```yaml
---
title: "{{DESCRIPTIVE_TITLE}}"
domains: [{{KEYWORD_1}}, {{KEYWORD_2}}, ..., {{KEYWORD_N}}]
---
```

- `title`: Human-readable title derived from the domain/topic name
- `domains`: 3-7 specific keywords derived from the entry content

## Domain Derivation

Domains are derived from spec content, not invented:

1. Extract tokens from project-global IDs in the file
2. Extract recurring domain-specific nouns from rule titles and requirement statements
3. Deduplicate, lowercase, keep only specific terms (not generic words like "rule", "spec")
4. Each file gets 3-7 specific keywords

For existing files lacking a `domains` field, derive and add it during the append.
