---
name: init-justfile
description: Scaffold a Justfile with standard zcode targets for the current project.
---

# /init-justfile

生成包含 zcode 标准 target 的 Justfile，作为测试/构建命令的抽象层。

## Prerequisites

**安装 just（>= 1.46.0）**

| 平台 | 命令 |
|------|------|
| macOS / Linux | `brew install just` |
| Windows (Scoop) | `scoop install just` |
| Windows (winget) | `winget install --id Casey.Just --exact` |
| Windows (Chocolatey) | `choco install just` |
| Cargo (通用) | `cargo install just` |

验证版本：

```bash
just --version
# 需要 >= 1.46.0（支持 [arg] flag 语法）
```

若版本低于 1.46.0，通过 Cargo 安装最新版：`cargo install just`

## 标准 Target 契约

| Target | 必须 | 用途 |
|--------|------|------|
| `test` | 是 | 单元 + 集成测试 |
| `test-e2e` | 否 | E2E 测试（`--feature` flag 切换模式） |
| `build` | 否 | 编译/打包 |
| `lint` | 否 | 静态分析 |

### test-e2e 参数

| 调用方式 | 说明 |
|----------|------|
| `just test-e2e` | 已毕业的回归测试（`tests/e2e/`，默认） |
| `just test-e2e --feature` | 当前 feature 的测试脚本 |

`test-e2e` 内部通过 `task feature` 获取当前 feature slug，无需在 justfile 中硬编码。

### 测试脚本生命周期

```
docs/features/<slug>/testing/scripts/  →  tests/e2e/<target>/
       ↑ 开发阶段                         ↑ 毕业后（回归阶段）
       just test-e2e --feature            just test-e2e
```

- **开发阶段**：`/gen-test-scripts` 生成脚本到 `docs/features/<slug>/testing/scripts/`，`just test-e2e --feature` 运行当前 feature 的测试
- **毕业**：`task all-completed` 在 e2e 测试首次通过后，按 target 将 spec 文件复制到 `tests/e2e/<target>/`，然后依次运行 `just test`（单元测试）和 `just test-e2e`（全量回归）
- **回归阶段**：`just test-e2e` 运行 `tests/e2e/` 下所有已毕业的 spec 文件

`task all-completed` 会按以下优先级检测测试命令：
1. `index.json` 中的 `testCommand`（显式配置）
2. `justfile`/`Justfile` 含 `test` recipe → `just test`
3. `Makefile` 含 `test` target → `make test`
4. 语言特定检测（go.mod / package.json / pytest.ini）

## 工作流

### Step 1: 检测项目类型

按优先级检查以下文件：

| 文件 | 项目类型 |
|------|---------|
| `go.mod` | Go |
| `package.json` | Node.js |
| `pyproject.toml` / `setup.py` | Python |
| `Cargo.toml` | Rust |
| 其他 | Generic |

```bash
ls go.mod package.json pyproject.toml setup.py Cargo.toml 2>/dev/null
```

### Step 2: 检查已有文件

```bash
ls justfile Justfile Makefile 2>/dev/null
```

- 若 `justfile` 或 `Justfile` 已存在 → 告知用户，询问是否覆盖，若否则中止
- 若 `Makefile` 已存在 → 读取内容，将已有 target 迁移到 Justfile（保留注释和命令）

### Step 3: 生成 Justfile

根据检测到的项目类型，写入 `justfile`（小写，just 的推荐命名）。

所有语言模板共享相同的 `test-e2e` 实现（`--feature` flag 切换模式），仅 `test`、`build`、`lint` 按语言不同。

**test-e2e 通用实现**：

```just
# Run e2e tests: "just test-e2e" (regression) or "just test-e2e --feature" (current feature)
[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    if [ "{{feature}}" != "" ]; then
        slug=$(task feature 2>/dev/null)
        if [ -z "$slug" ]; then
            echo "No active feature. Run: task feature <slug>" >&2; exit 1
        fi
        scripts_dir="docs/features/$slug/testing/scripts"
        fail=0
        for spec in "$scripts_dir"/*.spec.ts; do
            [ -f "$spec" ] && npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    else
        [ ! -d tests/e2e/node_modules ] && npm install --prefix tests/e2e
        fail=0
        for spec in $(find tests/e2e -mindepth 2 -name '*.spec.ts'); do
            npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    fi
```

**Go 项目：**

```just
# Run e2e tests: "just test-e2e" (regression) or "just test-e2e --feature" (current feature)
[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    if [ "{{feature}}" != "" ]; then
        slug=$(task feature 2>/dev/null)
        if [ -z "$slug" ]; then
            echo "No active feature. Run: task feature <slug>" >&2; exit 1
        fi
        scripts_dir="docs/features/$slug/testing/scripts"
        fail=0
        for spec in "$scripts_dir"/*.spec.ts; do
            [ -f "$spec" ] && npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    else
        [ ! -d tests/e2e/node_modules ] && npm install --prefix tests/e2e
        fail=0
        for spec in $(find tests/e2e -mindepth 2 -name '*.spec.ts'); do
            npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    fi

# Run unit and integration tests
test:
    go test -race ./...

# Build the project
build:
    go build ./...

# Run linters
lint:
    golangci-lint run ./...
```

**Rust 项目：**

```just
# Run e2e tests: "just test-e2e" (regression) or "just test-e2e --feature" (current feature)
[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    if [ "{{feature}}" != "" ]; then
        slug=$(task feature 2>/dev/null)
        if [ -z "$slug" ]; then
            echo "No active feature. Run: task feature <slug>" >&2; exit 1
        fi
        scripts_dir="docs/features/$slug/testing/scripts"
        fail=0
        for spec in "$scripts_dir"/*.spec.ts; do
            [ -f "$spec" ] && npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    else
        [ ! -d tests/e2e/node_modules ] && npm install --prefix tests/e2e
        fail=0
        for spec in $(find tests/e2e -mindepth 2 -name '*.spec.ts'); do
            npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    fi

# Run unit and integration tests
test:
    cargo test

# Build the project
build:
    cargo build --release

# Run linters
lint:
    cargo clippy -- -D warnings
```

**Node.js 项目：**

```just
# Run e2e tests: "just test-e2e" (regression) or "just test-e2e --feature" (current feature)
[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    if [ "{{feature}}" != "" ]; then
        slug=$(task feature 2>/dev/null)
        if [ -z "$slug" ]; then
            echo "No active feature. Run: task feature <slug>" >&2; exit 1
        fi
        scripts_dir="docs/features/$slug/testing/scripts"
        fail=0
        for spec in "$scripts_dir"/*.spec.ts; do
            [ -f "$spec" ] && npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    else
        [ ! -d tests/e2e/node_modules ] && npm install --prefix tests/e2e
        fail=0
        for spec in $(find tests/e2e -mindepth 2 -name '*.spec.ts'); do
            npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    fi

# Run unit and integration tests
test:
    npm test

# Build the project
build:
    npm run build

# Run linters
lint:
    npm run lint
```

**Python 项目：**

```just
# Run e2e tests: "just test-e2e" (regression) or "just test-e2e --feature" (current feature)
[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    if [ "{{feature}}" != "" ]; then
        slug=$(task feature 2>/dev/null)
        if [ -z "$slug" ]; then
            echo "No active feature. Run: task feature <slug>" >&2; exit 1
        fi
        scripts_dir="docs/features/$slug/testing/scripts"
        fail=0
        for spec in "$scripts_dir"/*.spec.ts; do
            [ -f "$spec" ] && npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    else
        [ ! -d tests/e2e/node_modules ] && npm install --prefix tests/e2e
        fail=0
        for spec in $(find tests/e2e -mindepth 2 -name '*.spec.ts'); do
            npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    fi

# Run unit and integration tests
test:
    pytest

# Build the project
build:
    python -m build

# Run linters
lint:
    ruff check .
```

**Generic（未识别项目类型）：**

```just
# Run e2e tests: "just test-e2e" (regression) or "just test-e2e --feature" (current feature)
[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    if [ "{{feature}}" != "" ]; then
        slug=$(task feature 2>/dev/null)
        if [ -z "$slug" ]; then
            echo "No active feature. Run: task feature <slug>" >&2; exit 1
        fi
        scripts_dir="docs/features/$slug/testing/scripts"
        fail=0
        for spec in "$scripts_dir"/*.spec.ts; do
            [ -f "$spec" ] && npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    else
        [ ! -d tests/e2e/node_modules ] && npm install --prefix tests/e2e
        fail=0
        for spec in $(find tests/e2e -mindepth 2 -name '*.spec.ts'); do
            npx tsx "$spec" || fail=$((fail+1))
        done
        [ "$fail" -eq 0 ]
    fi

# Run unit and integration tests
test:
    echo "TODO: implement test recipe"

# Build the project
build:
    echo "TODO: implement build recipe"

# Run linters
lint:
    echo "TODO: implement lint recipe"
```

### Step 4: 输出确认

生成完成后输出：

```
Created justfile with standard zcode targets (Go project)

Targets:
  just test              → go test -race ./...
  just test-e2e          → graduated regression tests in tests/e2e/
  just test-e2e --feature → current feature e2e (via task feature)
  just build             → go build ./...
  just lint              → golangci-lint run ./...

Edit justfile to customize commands for your project.
task all-completed will now use `just test` automatically.
```

## 注意事项

- **just 版本要求 >= 1.46.0**（`test-e2e` 使用 `[arg("feature", long)]` flag 语法）
- `just test-e2e`（默认）运行毕业后迁移到 `tests/e2e/` 的所有 spec 文件，用于日常回归
- `just test-e2e --feature` 运行当前 feature 开发阶段的测试脚本，由 `task all-completed` 调用
- 两种模式都逐个运行 spec 文件，确保即使部分失败也能收集完整结果
- `test-e2e` 内部通过 `task feature` 命令动态获取当前 feature slug，无需在 justfile 中硬编码
- 若从 Makefile 迁移，保留原有命令逻辑，仅调整格式（Makefile tab → just 4-space indent）
- 生成的 Justfile 是起点，用户应根据实际项目调整命令

## 相关命令

| 命令 | 用途 |
|------|------|
| `/init-zcode` | 安装 task-cli |
| `task all-completed` | 检查任务完成并运行测试 |
