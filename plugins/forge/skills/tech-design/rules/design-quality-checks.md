# Design Quality Check Rules

Rules for verifying design completeness before approval, during Step 5 of the tech-design flow.

## 5.1 PRD Coverage Verification

After drafting each section, verify every PRD acceptance criterion is addressed:

1. For each AC from `prd-user-stories.md`, identify which interface, model, or component handles it
2. If an AC has no corresponding design element, add one
3. Document the mapping in the "PRD Coverage Map" section of the template

## 5.2 Breakdown-Readiness Check

Before seeking approval, verify the design can be directly decomposed into implementation tasks:

| Check | Requirement |
|-------|-------------|
| Components enumerable | Can you list and count all components/modules by name? |
| Interfaces → tasks | Does each interface map to at least one implementation task? |
| Models → tasks | Does each data model map to at least one schema/migration task? |
| PRD AC coverage | Are all acceptance criteria from user stories addressed? |
| Cross-layer consistency | If feature spans layers, does the Data Map cover every field that crosses boundaries? |

If any check fails, add the missing detail before presenting to the user.

## 5.3 Cross-Layer Data Map

If the feature touches more than one architectural layer (database, API, UI, CLI, etc.):
- Complete the "Cross-Layer Data Map" table in the template
- Every field that appears in multiple layers must have a row showing its type/shape at each layer
- This becomes the Ground Truth for type decisions during task execution

If the feature is single-layer (e.g., only affects CLI output formatting):
- Write "Single-layer feature. Cross-Layer Data Map not applicable." in the section

## 5.4 Integration Specs

For each UI Function with `placement: existing-page:<route>`, generate an Integration Spec in the tech design document. Read the UI Design's Placement section for context.

The Integration Spec declares what file to modify and where:
- Do NOT specify implementation details (import statements, prop interfaces)
- Do specify: target file path, insertion point description, data source

This spec is consumed by breakdown-tasks to generate separate integration tasks.

If no UI Function has `placement: existing-page`, write "No existing-page integrations — not applicable."

## 5.5 DB Schema Branch (conditional)

**When `db-schema: "yes"`**:
1. Generate `design/er-diagram.md` using `templates/er-diagram.md` — Mermaid erDiagram + entity detail tables + index design + relationship descriptions
2. Generate `design/schema.sql` using `templates/schema.sql` — CREATE TABLE / ALTER TABLE with inline COMMENT syntax
3. Replace Data Models section in tech-design.md with cross-reference summary + Field Quick Reference table

**When `db-schema: "no"`**:
Data Models stays inline. After drafting, scan content for keywords: `TABLE`, `REFERENCES`, `FOREIGN KEY`, `CREATE TABLE`, `ALTER TABLE`, `migration`, `schema`. If found, prompt: "PRD marked db-schema 'no' but design references database tables. Generate er-diagram.md and schema.sql?" — Yes → proceed with db-schema "yes" path. No → keep inline.
