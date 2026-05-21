# Contract Validation Rules

After generating all Contracts for a Journey, validate each one:

| Check | Rule |
|-------|------|
| Mandatory dimensions | Each Outcome MUST have non-empty: Preconditions, Input, Output, State |
| Semantic descriptor purity | No dimension value may contain regex syntax |
| Outcome name uniqueness | Outcome names within a Step MUST be unique |
| Preconditions mutual exclusivity | Different Outcomes' Preconditions MUST be distinguishable |
| Journey Invariants | Every Contract file MUST have a `## Journey Invariants` section with at least 1 entry |
| Side-effect default | When Side-effect is omitted or empty, it defaults to `none` |
| Outcome count checkpoint | Steps with > 5 Outcomes trigger a review warning |
| Unclassified validation points | Any validation point that cannot be mapped to a dimension MUST go to Invariants with `dimension: unclassified` annotation |

**Validation failure handling**:
- If mandatory dimensions are empty: fix the Contract (add content from Journey + Fact Table)
- If semantic descriptors contain regex: rewrite as natural language
- If Preconditions are not mutually exclusive: differentiate or merge Outcomes
- If Journey Invariants are missing: generate from workflow analysis

<HARD-RULE>
- Semantic descriptors MUST NOT contain regex syntax.
- Outcome Preconditions MUST be mutually exclusive.
- Steps with > 5 Outcomes trigger an LLM review checkpoint.
- Validation points that cannot be classified into existing dimensions MUST go to Invariants with `dimension: unclassified` annotation.
</HARD-RULE>

# Error Handling

| Situation | Action |
|-----------|--------|
| Journey manifest missing | Abort with prompt to run `/gen-journeys` |
| Journey file not found | Abort with error listing the missing file path |
| Language detection fails | Load Convention files from `docs/conventions/` by `domains` frontmatter (match `testing`, `go`, `typescript`, etc.), extract from `Framework` section. Fallback: scan source/test files. If still unresolved, ask user |
| Interface detection fails | Ask user to configure `interfaces` in config.yaml |
| Source files not found for Fact Table | Mark as `UNKNOWN`, do not fabricate values |
| State verification level ambiguous | Default to `partial`, annotate with comment |
| Mandatory dimension empty after generation | Fix using Journey + Fact Table context, retry once |
| Semantic descriptor contains regex | Rewrite as natural language |
| Preconditions not mutually exclusive | Differentiate or merge Outcomes |
| Journey Invariants missing | Generate from workflow analysis |
