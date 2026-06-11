# File Signal Detection Reference

Detect test frameworks from project file signals. Detection is purely file-based (no code execution).

## Detection Algorithm

```
1. Language Detection (marker files) -> detected_languages
2. Framework Detection (dependency + file patterns) -> detected_frameworks
3. Cross-validation (eliminate false positives) -> validated_frameworks
```

## Language Detection

Scan the project root for marker files (`go.mod`, `package.json`, `Cargo.toml`, `pom.xml`, `build.gradle`, `pyproject.toml`, `setup.py`, `build.sbt`, `*.csproj`). Collect all detected languages.

- Empty result → error, ask user to specify `--scope`.
- Single language → use as `target_scope`.
- Multiple languages → list all, ask user which to generate Conventions for.
- `--scope` provided → skip language detection, use scope directly.

## Framework Detection

Each framework requires **3 signals** to confirm: **marker file** + **dependency** + **file pattern**. Use your knowledge of language-specific test frameworks to identify the correct signals for the detected language(s).

**Confidence levels:**
- **high**: All 3 signals confirmed
- **medium**: 2 of 3 signals confirmed (usually missing test files for cold start)
- **low**: Only marker file confirmed (language detected, framework unknown)

## False Positive Principles

- Framework dependency overrides generic dependency (e.g., `test` script in `package.json` determines winner between Vitest/Jest).
- UI libraries (React, Vue) are NOT test frameworks.
- A project can have multiple test frameworks serving different purposes (unit vs e2e) — report both.
- When two frameworks from the same category conflict, report both and let user choose.

## Detection Result Format

```
Detection Result:
  Language: <language> (from <marker file>)
  Frameworks:
    - <framework name> (confidence: high/medium/low)
      Signals: <list of matched signals>
      Missing signals: <list of unmatched signals, if any>
```

## Cold Start Handling

When test files are missing (medium/low confidence): use dependency-only detection to propose framework with medium confidence. If ambiguous, present cold start candidate list from `rules/convention-structure.md` and ask user.
