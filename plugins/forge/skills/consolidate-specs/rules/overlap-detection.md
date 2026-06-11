# Overlap Detection Rules

Before presenting to user (Step 6), scan for related existing entries. This applies to BOTH business rules (biz-specs) and technical specs (tech-specs):

1. **Decisions**: match by filename -> for each extracted entry's domain (biz or tech), map it to the corresponding decision type file using the table below. Then check `docs/decisions/<type>.md` for rows where the Decision column text matches the entry's topic.
2. **Lessons**: match by tags -> for each extracted entry, infer which of the 8 tag vocabulary items best matches its domain (e.g., auth rules -> `security`, error patterns -> `error-handling`). Then grep `tags:` frontmatter in `docs/lessons/*.md` for exact tag value matches from the 8-item vocabulary.

**Domain-to-decision-file mapping** (from decision-logging protocol):

| Spec domain keywords | Decision file |
|---------------------|---------------|
| system structure, layering, modules, architecture | `architecture.md` |
| API contracts, data shapes, serialization, interface | `interface.md` |
| schema, indexing, soft-delete, data model, data ownership | `data-model.md` |
| libraries, versions, packages, dependencies | `dependencies.md` |
| error types, status codes, error propagation, error handling | `error-handling.md` |
| test patterns, coverage, mocking, testing | `testing.md` |
| auth, permissions, data protection, security, access control | `security.md` |
| dev environment, tooling, deployment, local setup | `local-dev-deployment.md` |
| naming, conventions, coding standards | `architecture.md` |
| validation, state transitions, calculation rules | closest match or `architecture.md` |
| performance, latency, caching, throughput | `architecture.md` |

If a domain does not clearly map to any file, skip the decisions overlap check for that entry.

Collect matches as "Related existing entries" for display in Step 6.
