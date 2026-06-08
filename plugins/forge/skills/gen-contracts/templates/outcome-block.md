## Outcome "{{OUTCOME_NAME}}"
- Preconditions: "{{PRECONDITIONS}}"
  fixture_spec:
    entities:
      - entity_type: ""           # Required
        min_count: 1              # Required
        # relationship_type: ""   # Optional
        # parent_entity: ""       # Optional
        # field_constraints:      # Optional
        #   - field: ""
        #     value: ""
    # state_requirements:         # Optional
    #   - description: ""
    #     prerequisite_entity: ""
- Input: {{INPUT}}
- Output: {{OUTPUT}}
- State: {{STATE}}
- Side-effect: {{SIDE_EFFECT}}
{{STEP_INVARIANTS}}
