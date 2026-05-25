# Vocabulary Index Generation Rules

<!-- AUTO-GENERATED -- do not edit manually. Regenerated on every /consolidate-specs run. -->

Scan all four knowledge directories and produce a vocabulary index for use by `/learn` and auto-extract triggers. This step runs unconditionally -- even when knowledge directories are sparse or empty.

### Scan Targets

| Directory | What to extract | Source field |
|-----------|----------------|--------------|
| `docs/decisions/*.md` | Type names from decision row table, domain keywords from Decision column text | Table rows |
| `docs/lessons/*.md` | Tags from YAML frontmatter `tags:` field | `tags` frontmatter |
| `docs/conventions/*.md` | Domains from YAML frontmatter `domains:` field | `domains` frontmatter |
| `docs/business-rules/*.md` | Domains from YAML frontmatter `domains:` field | `domains` frontmatter |

### Base Vocabulary

The base 8-category vocabulary is always included, even when no knowledge files exist:

1. **architecture** -- system structure, layering, modules
2. **interface** -- API contracts, data shapes, serialization
3. **data-model** -- schema, indexing, soft-delete, data ownership
4. **dependencies** -- libraries, versions, packages
5. **error-handling** -- error types, status codes, error propagation
6. **testing** -- test patterns, coverage, mocking
7. **security** -- auth, permissions, data protection
8. **local-dev-deployment** -- dev environment, tooling, deployment

### Aggregation

1. **Types**: Collect unique knowledge types found: `decision`, `lesson`, `convention`, `business-rule`. Report which directories are non-empty vs empty.

2. **Domains**: Aggregate unique domain keywords from all scanned sources (tags from lessons + domains from conventions/business-rules + type-derived keywords from decisions). Merge with the base 8 categories. Deduplicate and normalize (lowercase, sorted).

3. **Counts**: For each type and domain keyword, count how many entries exist across all directories.

### Output Format

Use the output template from `templates/vocabulary-index.md`.

### Idempotency

The vocabulary file is fully regenerated on every `/consolidate-specs` run -- no incremental updates. Previous content is replaced entirely.

### Usage by Other Skills

- `/learn` reads `docs/.vocabulary.md` at runtime to suggest type and domain classifications for user input
- Auto-extract triggers (in `run-tasks`, `fix-bug`, `write-prd`, `tech-design`) read the vocabulary to classify extracted knowledge
- Both `/learn` and triggers accept values outside the vocabulary -- it is suggestive, not restrictive
