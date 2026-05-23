# File Signal Detection Reference

This document defines the complete detection signals for identifying test frameworks from project file signals. Detection is purely file-based (no code execution).

## Detection Algorithm

```
1. Language Detection (marker files) -> detected_languages
2. Framework Detection (dependency + file patterns) -> detected_frameworks
3. Cross-validation (eliminate false positives) -> validated_frameworks
```

## Language Detection

Scan the project root for marker files:

```bash
ls go.mod package.json Cargo.toml pom.xml build.gradle pyproject.toml setup.py build.sbt *.csproj 2>/dev/null
```

**Marker-to-language mapping:**

| Marker File       | Language       | Convention scope |
| ----------------- | -------------- | ---------------- |
| `go.mod`          | Go             | `go`             |
| `package.json`    | JavaScript/TS  | `javascript`     |
| `Cargo.toml`      | Rust           | `rust`           |
| `pom.xml`         | Java           | `java`           |
| `build.gradle`    | Java/Groovy    | `java`           |
| `pyproject.toml`  | Python         | `python`         |
| `setup.py`        | Python         | `python`         |
| `build.sbt`       | Scala          | `scala`          |
| `*.csproj`        | C# / .NET      | `dotnet`         |

**Classification algorithm:**

1. Check for each marker file's existence in the project root.
2. Collect all detected languages into `detected_languages`.
3. If `--scope` was provided: use that scope directly, skip language detection.
4. If `detected_languages` is empty: output error "No known project markers detected. Expected one of: go.mod, package.json, Cargo.toml, pom.xml, pyproject.toml, etc." and ask user to specify `--scope`.
5. If `detected_languages` has exactly one entry: use it as `target_scope`.
6. If `detected_languages` has multiple entries: list all detected languages and ask the user to select which one(s) to generate Conventions for (one Convention per language).

## Framework Detection Signals

Each framework requires THREE signals to confirm detection. A framework is "detected" when all three signal types are present.

### Go testing (standard library)

| Signal Type    | Detection Method                                                                                              | Required Signal                                          |
| -------------- | ------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| Marker file    | `go.mod` exists in project root                                                                               | `go.mod` present                                         |
| Test files     | Glob `**/*_test.go`                                                                                           | At least 1 file matching `*_test.go`                     |
| Import pattern | Grep `testing.T` in `*_test.go` files                                                                         | At least 1 file containing `*testing.T` in function signature |

**Confirmation**: `go.mod` + `*_test.go` + `testing.T` => Go testing detected.

### Ginkgo

| Signal Type    | Detection Method                                                                                              | Required Signal                                          |
| -------------- | ------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| Marker file    | `go.mod` exists in project root                                                                               | `go.mod` present                                         |
| Dependency     | Grep `ginkgo` in `go.mod` (look for `github.com/onsi/ginkgo` or `github.com/onsi/ginkgo/v2`)                 | ginkgo dependency listed in go.mod                       |
| Import pattern | Grep `. "github.com/onsi/ginkgo"` or `"github.com/onsi/ginkgo/v2"` in `*_test.go` files                      | At least 1 file importing ginkgo package                 |

**Confirmation**: `go.mod` + ginkgo in go.mod + ginkgo import in test file => Ginkgo detected.

### Vitest

| Signal Type    | Detection Method                                                                                              | Required Signal                                          |
| -------------- | ------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| Marker file    | `package.json` exists in project root                                                                         | `package.json` present                                   |
| Dependency     | Check `package.json` devDependencies or dependencies for `vitest`                                             | `vitest` listed as dependency                            |
| Test files     | Glob `**/*.test.{ts,js,tsx,jsx}` or `**/*.spec.{ts,js,tsx,jsx}`                                               | At least 1 file matching the pattern                     |

**Confirmation**: `package.json` + vitest dependency + `*.test.ts` (or spec) => Vitest detected.

### Jest

| Signal Type    | Detection Method                                                                                              | Required Signal                                          |
| -------------- | ------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| Marker file    | `package.json` exists in project root                                                                         | `package.json` present                                   |
| Dependency     | Check `package.json` devDependencies or dependencies for `jest`                                               | `jest` listed as dependency                              |
| Test files     | Glob `**/*.test.{ts,js,tsx,jsx}` or `**/*.spec.{ts,js,tsx,jsx}`                                               | At least 1 file matching the pattern                     |

**Confirmation**: `package.json` + jest dependency + test file => Jest detected.

### pytest

| Signal Type    | Detection Method                                                                                              | Required Signal                                          |
| -------------- | ------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| Marker file    | `pyproject.toml` or `setup.py` exists in project root                                                         | At least one present                                     |
| Dependency     | Grep `pytest` in `pyproject.toml` (under `[project.dependencies]` or `[tool.pytest]`) or `setup.py`          | pytest listed as dependency                              |
| Test files     | Glob `**/test_*.py` or `**/*_test.py`                                                                         | At least 1 file matching the pattern                     |

**Confirmation**: `pyproject.toml`/`setup.py` + pytest dependency + `test_*.py` => pytest detected.

### JUnit

| Signal Type    | Detection Method                                                                                              | Required Signal                                          |
| -------------- | ------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| Marker file    | `pom.xml` or `build.gradle` exists in project root                                                            | At least one present                                     |
| Dependency     | Grep `junit` in `pom.xml` (under `<dependencies>`) or `build.gradle`                                          | JUnit dependency listed                                  |
| Test files     | Glob `**/*Test.java` or `**/*Tests.java`                                                                      | At least 1 file matching the pattern                     |

**Confirmation**: `pom.xml`/`build.gradle` + JUnit dependency + `*Test.java` => JUnit detected.

### Rust (cargo test)

| Signal Type    | Detection Method                                                                                              | Required Signal                                          |
| -------------- | ------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| Marker file    | `Cargo.toml` exists in project root                                                                           | `Cargo.toml` present                                     |
| Test attribute | Grep `#[cfg(test)]` or `#[test]` in Rust source files                                                         | At least 1 file containing test attribute                |
| Test directory | Check `tests/` directory exists at project root (for integration tests)                                       | `tests/` directory present OR `#[cfg(test)]` in src/     |

**Confirmation**: `Cargo.toml` + `#[cfg(test)]`/`#[test]` + tests/ (or test in src) => Rust testing detected.

## False Positive Exclusion Rules

To avoid misclassification, apply these exclusion rules after detection:

### Rule 1: Framework dependency overrides generic dependency

If both Vitest and Jest dependencies are found in `package.json`, check which one appears in `devDependencies` with the explicit test script:

```bash
grep -E '"test".*vitest|"test".*jest' package.json
```

The framework referenced by the `test` script wins.

### Rule 2: Test file patterns are language-specific

Do NOT use generic test patterns across languages:

- `package.json` containing React dependency does NOT imply a test framework. React is a UI library, not a test framework. Only match test frameworks explicitly listed in dependencies.
- A `package.json` with `vitest` AND `@playwright/test` should report BOTH as detected frameworks (they serve different purposes -- unit vs e2e).

### Rule 3: Ginkgo vs go testing precedence

If `go.mod` contains ginkgo dependency AND test files import ginkgo, classify as Ginkgo. If `go.mod` has no ginkgo dependency but test files use `testing.T`, classify as go testing. If both signals exist (some files use ginkgo, others use testing.T), report BOTH and let user choose.

### Rule 4: pytest vs unittest

If `pyproject.toml`/`setup.py` has both pytest and no external dependencies, check test file imports:

- Files importing `pytest` => pytest
- Files importing `unittest` only => unittest
- Both present => report both, let user choose

### Rule 5: Build tool ambiguity (Java)

If both `pom.xml` and `build.gradle` exist:

- Check each independently for JUnit dependency
- If JUnit found in both: classify by whichever has more test files matching the convention
- Report a single result (JUnit), not two detections

## Detection Result Format

After running detection, output a structured result:

```
Detection Result:
  Language: <language> (from <marker file>)
  Frameworks:
    - <framework name> (confidence: high/medium/low)
      Signals: <list of matched signals>
      Missing signals: <list of unmatched signals, if any>
```

**Confidence levels:**
- **high**: All 3 signal types confirmed (marker + dependency + file pattern)
- **medium**: 2 of 3 signals confirmed (usually missing test files for cold start)
- **low**: Only 1 signal confirmed (marker file only -- treat as "language detected, framework unknown")

## Cold Start Handling

When test files are missing (medium or low confidence), use dependency-only detection:

1. If the dependency signal identifies a unique framework (e.g., `vitest` in package.json), propose that framework with medium confidence.
2. If dependency is ambiguous or missing, present the cold start candidate list from `rules/convention-structure.md` and ask user to select.
