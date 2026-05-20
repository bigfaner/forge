# Spec Classification and ID Rules

## CROSS/LOCAL Classification

For each extracted rule or spec:
- `[CROSS]`: Referenced by 2+ features, or expresses a domain invariant (not feature behavior), or establishes a naming/error-handling convention
- `[LOCAL]`: Only meaningful within this feature's scope

## Project-Global ID Encoding

Each entry gets a project-global ID (not the feature-local BIZ-NNN/TECH-NNN):

- **Prefix** derived from target filename: `business-rules/auth.md` -> prefix `auth`, `conventions/api.md` -> prefix `api`
- **Sequence**: file-internal -- find max existing NNN in the target file + 1
- **Format**: `BIZ-<domain>-<NNN>` for business rules, `TECH-<topic>-<NNN>` for tech specs
- **Examples**: `BIZ-auth-001`, `TECH-api-003`
- **Source traceability**: `Source: feature/<slug> BIZ-001` (links back to feature-local preview ID)

## Preview ID Numbering

Feature-local IDs in preview files use sequential 3-digit numbering starting at 001, independent per file:
- `biz-specs.md`: BIZ-001, BIZ-002, ...
- `tech-specs.md`: TECH-001, TECH-002, ...

## New File Frontmatter

When creating a new project-level spec file, include this frontmatter:

```yaml
---
title: "<Descriptive Title>"
domains: [<keyword1>, <keyword2>, ..., <keywordN>]
---
```

- `title`: Human-readable title derived from the domain/topic name (existing behavior, unchanged)
- `domains`: 3-7 specific keywords derived from the spec content being written into the file, per the Domain Derivation Rules in `rules/domain-frontmatter.md`

For **existing files** that lack a `domains` field, derive and add it during integration (do not modify existing `title`).
