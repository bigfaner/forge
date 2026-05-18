# Step 0: Resolve Language

1. **Detect language**: Run `forge testing detect` to auto-detect the project's test language(s) from file signals.
2. **On failure** (no language detected): ask the user to add `languages` to `.forge/config.yaml` (e.g., `languages: [go]`).

<HARD-RULE>
Do NOT silently default to any language. If `forge testing detect` returns no result and the user cannot configure `languages`, abort the skill.
</HARD-RULE>
