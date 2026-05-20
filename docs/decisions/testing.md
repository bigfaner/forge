# Testing Decisions

| Date | Feature | Decision | Rationale | Source |
|------|---------|----------|-----------|--------|
| 2026-05-20 | test-knowledge-convention-driven | Convention files replace hardcoded profile package for test generation | Convention files (testing-{framework}.md) decouple Forge from language-specific code, allow projects to extend without plugin changes, and enable AI agents to load only relevant conventions via domains frontmatter | task 0.1 POC |
