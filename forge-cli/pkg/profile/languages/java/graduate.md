# JUnit 5 Graduate Strategy

Profile-specific graduation rules for the `graduate-tests` skill.

## Source File Discovery

| Item | Value |
|------|-------|
| File extension | `Test.java` or `E2E.java` |
| Source directory | `tests/e2e/features/<slug>/` (staging) |
| Target directory | `tests/e2e/<module>/` (regression) |

## Import Rewrite

Java uses package imports — no path rewriting needed. Graduated files retain their `package` and `import` statements unchanged.

## Validation

| Check | Command | Failure action |
|-------|---------|----------------|
| Pre-flight compilation | `just e2e-compile` | Abort before touching anything |
| Post-migration compilation | `just e2e-compile` | Rollback via migration manifest |
| Test discovery | `just e2e-discover` | Rollback via migration manifest |

## Compilation

Maven Surefire auto-discovers test classes matching `*Test.java`, `*E2E.java` patterns. No manual test registration required.

## Merge Procedure

When a target file already exists at the graduation destination:

1. Read both source and target Java files
2. Backup target file (only if no backup exists — prevents overwriting original on re-run)
3. Combine imports, deduplicate
4. Match test methods by name — deduplicate by method name (identical names: keep source version; same TC ID prefix but different names: keep both)
5. Preserve class structure (package, class declaration, annotations)
6. Append new methods that don't exist in target
7. Write the merged file

### Class Merge Rules

- Keep the target's `package` declaration
- Merge `import` statements (deduplicate)
- Keep `@Tag` and other class-level annotations from both (deduplicate)
- Merge `@BeforeAll`/`@AfterAll` static setup blocks: combine logic, deduplicate
- Merge test methods: deduplicate by method name

## Graduation Marker

Written only after validation passes (atomic — no marker = not graduated):

```yaml
schema_version: 1
status: completed
timestamp: <UTC ISO timestamp>
source: tests/e2e/features/<slug>/
targets:
  - tests/e2e/<module>/<java-file>
modules:
  - <module-name>
testCount: <N>
```
