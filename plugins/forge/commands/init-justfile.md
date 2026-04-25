---
name: init-justfile
description: Scaffold a Justfile with standard forge targets for the current project.
---

# /init-justfile

生成包含 forge 标准 target 的 Justfile，作为测试/构建命令的抽象层。

## Prerequisites

**安装 just（>= 1.50.0）**

| 平台 | 命令 |
|------|------|
| macOS / Linux | `brew install just` |
| Windows (Scoop) | `scoop install just` |
| Windows (winget) | `winget install --id Casey.Just --exact` |
| Cargo (通用) | `cargo install just` |

```bash
just --version  # 需要 >= 1.50.0（支持 [arg] 命名选项语法）
```

若版本低于 1.50.0：`cargo install just`

## 标准 Target 契约

| Target | 必须 | 用途 |
|--------|------|------|
| `test` | 是 | 单元 + 集成测试 |
| `test-e2e` | 否 | E2E 测试 |
| `build` | 否 | 编译/打包 |
| `lint` | 否 | 静态分析 |

`test-e2e` 调用方式：

| 调用 | 说明 |
|------|------|
| `just test-e2e` | 回归测试（`tests/e2e/`） |
| `just test-e2e --feature <slug>` | 指定 feature 的测试脚本 |

## 工作流

### Step 1: 检测项目类型

```bash
ls go.mod package.json pyproject.toml setup.py Cargo.toml 2>/dev/null
```

| 文件 | 项目类型 |
|------|---------|
| `go.mod` | Go |
| `package.json` | Node.js |
| `pyproject.toml` / `setup.py` | Python |
| `Cargo.toml` | Rust |
| 其他 | Generic |

### Step 2: 检查已有文件

```bash
ls justfile Justfile Makefile 2>/dev/null
```

- 若 `justfile` 或 `Justfile` 已存在 → 询问是否覆盖，若否则中止
- 若 `Makefile` 已存在 → 读取内容，将已有 target 迁移到 Justfile

### Step 3: 生成 Justfile

写入 `justfile`（小写）。所有模板共用以下 `test-e2e`，仅 `test`/`build`/`lint` 按语言不同。

**test-e2e（所有语言通用）：**

```just
# Run e2e tests: "just test-e2e" (regression) or "just test-e2e --feature <slug>" (feature tests)
[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    if [ "{{feature}}" != "" ]; then
        scripts_dir="docs/features/{{feature}}/testing/scripts"
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

**语言特定 recipes：**

Go:
```just
test:
    go test -race ./...
build:
    go build ./...
lint:
    golangci-lint run ./...
```

Rust:
```just
test:
    cargo test
build:
    cargo build --release
lint:
    cargo clippy -- -D warnings
```

Node.js:
```just
test:
    npm test
build:
    npm run build
lint:
    npm run lint
```

Python:
```just
test:
    pytest
build:
    python -m build
lint:
    ruff check .
```

Generic:
```just
test:
    echo "TODO: implement test recipe"
build:
    echo "TODO: implement build recipe"
lint:
    echo "TODO: implement lint recipe"
```

### Step 4: 输出确认

```
Created justfile with standard forge targets (Go project)

Targets:
  just test                       → go test -race ./...
  just test-e2e                   → regression tests in tests/e2e/
  just test-e2e --feature <slug>  → feature tests in docs/features/<slug>/testing/scripts/
  just build                      → go build ./...
  just lint                       → golangci-lint run ./...

Edit justfile to customize commands for your project.
task all-completed will now use `just test` automatically.
```

## 注意事项

- **just >= 1.50.0**：`[arg("feature", long)]` 生成 `--feature <value>` 命名选项，调用时必须带值
- 调用方（CI、`task all-completed`）负责传入 slug：`just test-e2e --feature <slug>`
- 若从 Makefile 迁移，保留原有命令逻辑，仅调整格式
