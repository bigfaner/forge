<!-- Snippet to append/update in manifest.md after /tech-design completes -->

## Documents (updated)

Add rows:
| Tech Design | design/tech-design.md | {{TECH_DESIGN_SUMMARY}} |
| API Handbook | design/api-handbook.md | {{API_HANDBOOK_SUMMARY}} |
| ER Diagram | design/er-diagram.md | {{ER_DIAGRAM_SUMMARY}} |  <!-- only when db-schema="yes" -->
| SQL Schema | design/schema.sql | {{SCHEMA_SQL_SUMMARY}} |    <!-- only when db-schema="yes" -->

## Traceability (updated)

Add entries linking PRD sections to design sections:
| PRD Section | Design Section | UI Component | Tasks |
|-------------|----------------|--------------|-------|
| "{{PRD_SECTION}}" (prd-spec §N) | "{{DESIGN_SECTION}}" (tech-design §N) | — | <!-- task IDs added by /breakdown-tasks --> |

## Frontmatter

Update `status` to `design` if /ui-design already completed or if UI is not applicable.
