# Deviation Categories

Use these categories to classify each finding in the forensic analysis.

| Category | Description | Example |
|----------|-------------|---------|
| `instruction-gap` | Skill definition missing a critical rule | No instruction to handle MAIN_SESSION flag |
| `context-starvation` | Agent lacked necessary information | Agent didn't see the record.json content |
| `trust-without-verify` | Agent trusted its own output | Marked AC as met without running the artifact |
| `wrong-priority` | Agent followed wrong priority | Chose "efficiency" over "safety" |
| `scope-creep` | Agent exceeded its defined scope | Task executor claimed multiple tasks |
| `pipeline-gap` | No enforcement between stages | Dispatcher checked file existence, not content |
