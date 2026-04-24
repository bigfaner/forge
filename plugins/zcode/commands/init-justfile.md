---
name: init-justfile
description: Scaffold a Justfile with standard zcode targets for the current project.
---

# /init-justfile

生成包含 zcode 标准 target 的 Justfile，作为测试/构建命令的抽象层。

## 标准 Target 契约

| Target | 必须 | 用途 |
|--------|------|------|
| `test` | 是 | 单元 + 集成测试 |
| `test-e2e` | 否 | Feature e2e 测试 |
| `build` | 否 | 编译/打包 |
| `lint` | 否 | 静态分析 |

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

**Go 项目：**

```just
# Run unit and integration tests
test:
    go test -race ./...

# Run feature e2e tests (replace <slug> with actual feature slug, or run: task feature)
test-e2e:
    npm run test:all --prefix docs/features/<slug>/testing/scripts

# Build the project
build:
    go build ./...

# Run linters
lint:
    golangci-lint run ./...
```

**Rust 项目：**

```just
# Run unit and integration tests
test:
    cargo test

# Run feature e2e tests (replace <slug> with actual feature slug, or run: task feature)
test-e2e:
    npm run test:all --prefix docs/features/<slug>/testing/scripts

# Build the project
build:
    cargo build --release

# Run linters
lint:
    cargo clippy -- -D warnings
```

**Node.js 项目：**

```just
# Run unit and integration tests
test:
    npm test

# Run feature e2e tests
test-e2e:
    npm run test:e2e

# Build the project
build:
    npm run build

# Run linters
lint:
    npm run lint
```

**Python 项目：**

```just
# Run unit and integration tests
test:
    pytest

# Run feature e2e tests
test-e2e:
    pytest tests/e2e/

# Build the project
build:
    python -m build

# Run linters
lint:
    ruff check .
```

**Generic（未识别项目类型）：**

```just
# Run unit and integration tests
test:
    echo "TODO: implement test recipe"

# Run feature e2e tests
test-e2e:
    echo "TODO: implement test-e2e recipe"

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
✓ Created justfile with standard zcode targets (Go project)

Targets:
  just test       → go test -race ./...
  just test-e2e   → npm run test:all --prefix ...
  just build      → go build ./...
  just lint       → golangci-lint run ./...

Edit justfile to customize commands for your project.
task all-completed will now use `just test` automatically.
```

## 注意事项

- `test-e2e` 中的 `<slug>` 需替换为实际 feature slug（运行 `task feature` 查看当前 slug）
- 若从 Makefile 迁移，保留原有命令逻辑，仅调整格式（Makefile tab → just 4-space indent）
- 生成的 Justfile 是起点，用户应根据实际项目调整命令

## 相关命令

| 命令 | 用途 |
|------|------|
| `/init-zcode` | 安装 task-cli |
| `task all-completed` | 检查任务完成并运行测试 |
