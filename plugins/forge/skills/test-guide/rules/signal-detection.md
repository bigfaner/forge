# File Signal Detection Reference

## Language Detection

Scan the project root for marker files:

```bash
ls go.mod package.json Cargo.toml pom.xml build.gradle pyproject.toml setup.py build.sbt *.csproj 2>/dev/null
```

**Detection signal mapping:**

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

## Framework Detection

For each detected language, probe deeper for framework-specific signals:

| Language       | Probe Command                                                                    | Signal                     |
| -------------- | -------------------------------------------------------------------------------- | -------------------------- |
| Go             | `grep -l 'ginkgo' go.mod 2>/dev/null`                                            | Ginkgo if present          |
| Go             | `grep -l 'testify' go.mod 2>/dev/null`                                           | testify (default)          |
| JavaScript/TS  | `node -e "const d=JSON.parse(require('fs').readFileSync('package.json'));console.log(Object.keys(d.devDependencies||{}).concat(Object.keys(d.dependencies||{})).filter(p=>/vitest|jest|mocha|cypress|playwright/.test(p)).join(' '))" 2>/dev/null` | Framework from deps |
| Rust           | Default: `cargo test`                                                            | Standard                   |
| Python         | `grep -l 'pytest\|unittest\|nose' pyproject.toml setup.py 2>/dev/null`           | pytest/unittest            |
| Java           | `grep -l 'junit\|testng\|spock' pom.xml build.gradle 2>/dev/null`                | JUnit/TestNG/Spock         |

Record detected frameworks in `detected_frameworks`. This is a warm-start signal -- it narrows the candidate list but does NOT override Step 2's test file analysis.

## Complete File Signal Reference

| Signal File           | Language      | Framework candidates                            |
| --------------------- | ------------- | ----------------------------------------------- |
| `go.mod`              | Go            | go testing, Ginkgo                              |
| `package.json`        | JavaScript/TS | Vitest, Jest, Mocha, Cypress, Playwright        |
| `Cargo.toml`          | Rust          | cargo test (built-in)                           |
| `pom.xml`             | Java          | JUnit 4/5, TestNG, Spock                        |
| `build.gradle`        | Java/Groovy   | JUnit 4/5, TestNG, Spock, GroovyTestCase        |
| `pyproject.toml`      | Python        | pytest, unittest, nose2                         |
| `setup.py`            | Python        | pytest, unittest, nose2                         |
| `build.sbt`           | Scala         | ScalaTest, specs2, munit                        |
| `*.csproj`            | C# / .NET     | xUnit, NUnit, MSTest                            |
| `go.sum`              | Go            | (secondary -- use go.mod instead)               |
| `package-lock.json`   | JavaScript/TS | (secondary -- use package.json instead)         |
| `yarn.lock`           | JavaScript/TS | (secondary -- use package.json instead)         |
