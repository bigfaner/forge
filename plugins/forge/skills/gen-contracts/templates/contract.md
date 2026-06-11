---
journey: "{{JOURNEY_NAME}}"
step: "{{STEP_NUMBER}}"
step-action: "{{STEP_ACTION}}"
generated: "{{DATE}}"
sources:
  - docs/features/{{FEATURE_SLUG}}/testing/{{JOURNEY_NAME}}/journey.md

# Technical Anchors — filled by gen-contracts from handbook/page-map/screen-map.
# Omitted when handbook does not exist; pipeline degrades gracefully.
anchors:
  # API surface (source: api-handbook)
  api:
    endpoint: ""        # required, format: /path/:param
    method: ""          # required, HTTP verb (GET, POST, PUT, DELETE, PATCH)
    content_type: ""    # optional, e.g. application/json
    auth_required:      # optional, boolean

  # CLI surface (source: cli-handbook)
  cli:
    command: ""         # required, e.g. forge surfaces list
    subcommand: ""      # optional, nested subcommand
    flags: []           # optional, list of flag names
    aliases: []         # optional, list of command aliases

  # TUI surface (source: cli-handbook)
  tui:
    command: ""         # required, entry command for interactive terminal
    interactive_prompt: ""  # optional
    keybindings: []     # optional, list of keybinding identifiers

  # Web surface (source: page-map)
  web:
    page: ""            # required, page name
    route: ""           # optional, URL path
    requires_auth:      # optional, boolean
    layout: ""          # optional, layout template name

  # Mobile surface (source: screen-map)
  mobile:
    screen: ""          # required, screen name
    navigation_path: [] # optional, navigation hierarchy
    deeplink: ""        # optional, deep link URL scheme
    platform: ""        # optional, e.g. ios, android

last_anchor_sync: ""    # ISO-8601 timestamp of last anchor fill from handbook
---

# Contract: {{JOURNEY_NAME}} / Step {{STEP_NUMBER}}: {{STEP_ACTION}}

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

{{OUTCOME_BLOCKS}}

## Journey Invariants

{{JOURNEY_INVARIANTS}}

## Fixture Specification

This Contract requires the following pre-existing data state. See `rules/fixture-spec.md` for schema details.

```yaml
fixture_spec:
  entities:
    - entity_type: ""           # Required: entity type name (e.g., "Project", "Milestone")
      min_count: 1              # Required: minimum count >= 1
      # relationship_type: ""   # Optional: "belongs_to" | "has_many" | "has_one"
      # parent_entity: ""       # Optional: parent entity type name
      # field_constraints:      # Optional: field value constraints
      #   - field: ""
      #     value: ""
  # state_requirements:         # Optional: system-level state
  #   - description: ""
  #     prerequisite_entity: ""
```
