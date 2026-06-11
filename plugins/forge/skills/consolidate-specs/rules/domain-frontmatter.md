# Domain Frontmatter Rules

Convention and business-rule files carry a `domains` field in their YAML frontmatter that enables lightweight discovery by consumers (prompt templates, commands, agents).

```yaml
---
title: "Error Handling Conventions"  # existing field, unchanged
domains: [error, status, response, stderr]  # keywords this file covers
---
```

## Domain Derivation Rules

Domains are **derived programmatically from spec content** -- never invented by the agent. The derivation algorithm:

1. **Spec ID keywords**: Extract tokens from project-global IDs in the file (e.g., `BIZ-auth-001` contributes `auth`, `TECH-api-003` contributes `api`)
2. **Source keywords**: Extract recurring domain-specific nouns from rule titles, requirement statements, and source references (e.g., a rule about "token validation" contributes `token`, `validation`)
3. **Deduplicate and normalize**: Lowercase, remove duplicates, keep only specific terms (not generic words like "rule", "spec", "requirement")
4. **Cardinality**: Each file gets **3-7 specific keywords**

## Domain Overlap Detection

When multiple files in the same directory (`docs/conventions/` or `docs/business-rules/`) have `domains` fields, compute keyword overlap:

- **Overlap ratio** = `|intersection(domains_A, domains_B)| / min(|domains_A|, |domains_B|)`
- **Threshold**: If overlap ratio > 50%, flag as a potential duplicate/merge candidate during the user confirmation step (Step 6)
- **Action**: Display the warning; the user decides whether to merge or keep separate
