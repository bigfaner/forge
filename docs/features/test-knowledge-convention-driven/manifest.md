---
feature: "test-knowledge-convention-driven"
status: tasks
---

# Feature: test-knowledge-convention-driven

<!-- Status flow: prd → design → tasks → in-progress → completed -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | Replace Profile system with user-editable Convention files; full Profile removal, config cleanup, skill rewrites, new test-guide command; compile gate with recovery; phase gates with go/no-go criteria |
| User Stories | prd/prd-user-stories.md | 6 stories with error scenario ACs covering non-default frameworks, Convention bootstrap, cold start, backward compatibility, multi-framework management, and compile gate recovery |
| Tech Design | design/tech-design.md | Three-layer Convention-driven architecture; pkg/forgeconfig extraction; 12 consumer files enumerated; section-level merge semantics; regression-first testing |
| Tasks | tasks/index.json | 29 tasks across 5 phases (POC, Profile removal, Skill rewrites, test-guide, Validation) |

## Traceability

| PRD Section | Design Section | UI Component | Placement | Tasks |
|-------------|----------------|--------------|-----------|-------|
| FS-1 Convention structure (prd-spec) | Data Models (tech-design) | — | — | 3.1, 3.2 |
| FS-2 Convention loading (prd-spec) | I-6 gen-test-scripts (tech-design) | — | — | 2.1 |
| FS-3 Code Reconnaissance (prd-spec) | I-6 gen-test-scripts (tech-design) | — | — | 2.1 |
| FS-4 Compile gate (prd-spec) | Error Handling (tech-design) | — | — | 2.1, 4.1 |
| FS-5 test-guide (prd-spec) | I-6 test-guide (tech-design) | — | — | 3.1 |
| FS-6 config.yaml cleanup (prd-spec) | I-1 forgeconfig (tech-design) | — | — | 1.1, 1.4 |
| FS-7 Profile removal (prd-spec) | Consumer File Enumeration (tech-design) | — | — | 1.1–1.5 |
| FS-8 Silent migration (prd-spec) | I-1 forgeconfig (tech-design) | — | — | 1.1 |
| FS-9 consolidate-specs (prd-spec) | Architecture (tech-design) | — | — | 3.2 |
| Story 1: Non-default framework (user-stories) | I-6 gen-test-scripts (tech-design) | — | — | 0.1, 2.1 |
| Story 2: Convention bootstrap (user-stories) | I-6 test-guide (tech-design) | — | — | 3.1, 3.2 |
| Story 3: Cold start (user-stories) | I-6 gen-test-scripts (tech-design) | — | — | 0.1, 2.1, 2.3 |
| Story 4: Backward compatibility (user-stories) | I-1 forgeconfig (tech-design) | — | — | 1.1, 1.5, 4.1, 4.2 |
| Story 5: Multi-framework (user-stories) | Convention merge semantics (tech-design) | — | — | 3.1, 3.2 |
| Story 6: Compile gate recovery (user-stories) | Error Handling (tech-design) | — | — | 2.1 |
