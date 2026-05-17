# Convention / Business-Rule Entry Template

Used by `/learn` when writing to `docs/conventions/` or `docs/business-rules/`.

## Entry Format

Append to the target file using the project-global ID encoding:

### Convention Entry (docs/conventions/<topic>.md)

```markdown
### TECH-<topic>-<NNN>: <Spec Title>

**Requirement**: <concise requirement>
**Scope**: [CROSS]
**Source**: /learn entry <YYYY-MM-DD>

<Implementation details, examples>
```

### Business-Rule Entry (docs/business-rules/<domain>.md)

```markdown
### BIZ-<domain>-<NNN>: <Rule Title>

**Rule**: <concise rule statement>
**Context**: <why this rule exists>
**Scope**: [CROSS]
**Source**: /learn entry <YYYY-MM-DD>

<Additional details, examples, or edge cases>
```

## Project-Global ID Encoding

- **Prefix**: derived from target filename (e.g., `business-rules/auth.md` -> `auth`, `conventions/api.md` -> `api`)
- **Sequence**: file-internal -- find max existing NNN in the target file + 1
- **Format**: `BIZ-<domain>-<NNN>` for business rules, `TECH-<topic>-<NNN>` for tech specs
- **Examples**: `BIZ-auth-001`, `TECH-api-003`

## New File Frontmatter

When creating a new file, include YAML frontmatter:

```yaml
---
title: "<Descriptive Title>"
domains: [<keyword1>, <keyword2>, ..., <keywordN>]
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
