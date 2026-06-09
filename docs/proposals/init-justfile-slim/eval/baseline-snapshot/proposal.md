---
created: 2026-06-09
author: faner
status: Draft
---

# 提案：init-justfile 精简 — CLI scaffold 替代 prompt 层模板

## 问题

`init-justfile` skill 当前 1645 行（8 个文件），是 Forge 最大的 skill。臃肿的根本原因：**将机械性的 bash 代码模板（PID 管理、健康检查、teardown）维护在 LLM prompt 中，依赖 agent 理解和复用。** 这带来三个问题：

1. **Token 浪费**：`server-lifecycle.md`（745 行，占 45%）是纯 bash 代码模板，让 LLM 逐行理解这些代码是低效的——它应该只关心"用什么命令启动服务、端口是多少"
2. **维护分散**：surface type 的差异化逻辑分布在 5 个 rule 文件（318 行）中，每个文件重复相同的表格结构（orchestration sequence、recipe invocation contract、journey filter）
3. **职责错位**：recipe 的 bash 骨架（PID 文件路径、idempotent start、健康检查重试）是确定性逻辑，适合程序生成而非 LLM 推理

### 证据

- `server-lifecycle.md` 中 745 行 bash 代码，agent 每次调用 `/init-justfile` 都需要加载并理解
- 5 个 surface rule 文件之间高度同构：api.md / web.md / mobile.md 共享 dev→probe→test→teardown 序列，cli.md / tui.md 共享 test→teardown 序列，但各自独立维护
- SKILL.md 的 EXTREMELY-IMPORTANT 块和 Notes 段与 body 内容大量重复（如 "CLI/TUI 不生成 dev/probe" 出现 3 次）
- Phase 1 consistency check 对 LLM 生成结果做防御性校验，但如果 producer 是可信的程序，这一层不再必要

## 提议方案

### 核心思路

将 **recipe 代码生成** 从 prompt 层下沉到 Forge CLI 命令，agent 职责从"理解 bash 模板并生成代码"变为"调用命令 + 检测语言 + 填占位符"。

### 新增：`forge justfile scaffold` CLI 命令

| 接口 | 说明 |
|------|------|
| `forge justfile scaffold --type <type> --key <key>` | 为单个 surface 生成完整 recipe 集（lifecycle + quality），输出到 stdout |
| `forge justfile scaffold --aggregate` | 生成跨 surface 聚合 recipe（install / ci / clean），读取 `forge surfaces` 获取全部 surface |

**输入参数**：

| 参数 | 必选 | 说明 |
|------|------|------|
| `--type` | 是 | Surface type：`cli` / `tui` / `api` / `web` / `mobile` |
| `--key` | 否 | Surface key（named surface 必传，scalar 省略） |

**输出**：带 `{{PLACEHOLDER}}` 占位符的 just recipe 代码，输出到 stdout。CLI 不需要知道任何项目细节。

**每个 surface type 生成的 recipes**：

| Surface Type | Lifecycle Recipes | Quality Recipes | Aggregate |
|--------------|------------------|----------------|-----------|
| `cli` | `test`, `teardown` | `compile`, `fmt`, `lint`, `unit-test` | 无 |
| `tui` | `test`, `teardown` | `compile`, `fmt`, `lint`, `unit-test` | 无 |
| `api` | `dev`, `probe`, `test`, `teardown`, `<key>` | `compile`, `fmt`, `lint`, `unit-test` | `<key>` (dev→probe→test→teardown) |
| `web` | `dev`, `probe`, `test`, `teardown`, `<key>` | `compile`, `fmt`, `lint`, `unit-test` | `<key>` (dev→probe→test→teardown) |
| `mobile` | `test-setup`, `dev`, `probe`, `test`, `teardown`, `<key>` | `compile`, `fmt`, `lint`, `unit-test` | `<key>` (test-setup→dev→probe→test→teardown) |

Named surface（有 key）：所有 recipe 名加 `<key>-` 前缀（如 `backend-dev`）。Scalar surface（无 key）：无前缀（如 `dev`）。

所有 recipe 均包含 `[linux]` 和 `[windows]` 双平台变体。Lifecycle recipes 标记 `# user-customized`。

**占位符清单**：

| 占位符 | 说明 | Agent 解析来源 |
|--------|------|---------------|
| `{{START_CMD}}` | 启动服务的命令 | `package.json` scripts / `go run` + entry point / Convention |
| `{{PORT}}` | 服务监听端口 | Convention / 常见默认值（Node 3000, Go 8080） |
| `{{HEALTH_URL}}` | 健康检查 URL | `http://localhost:{{PORT}}/healthz`（默认） |
| `{{COMPILE_CMD}}` | 编译命令 | `go build ./...` / `npm run build` / Convention |
| `{{UNIT_TEST_CMD}}` | 单元测试命令 | `go test -v ./...` / `npm test` / Convention |
| `{{LINT_CMD}}` | 静态分析命令 | `golangci-lint run` / `npm run lint` / Convention |
| `{{FMT_CMD}}` | 格式化命令 | `gofmt -w .` / `npm run fmt` / Convention |
| `{{BUILD_CMD}}` | 构建命令 | `go build -o bin/app` / `npm run build` / Convention |
| `{{CLEAN_CMD}}` | 清理命令 | `go clean` / `rm -rf dist` / Convention |
| `{{INSTALL_CMD}}` | 安装依赖命令 | `go mod download` / `npm install` / Convention |
| `{{TEST_CMD}}` | Surface-level 测试命令 | Convention test runner + file pattern |

**聚合 recipe 生成**（`--aggregate` 模式）：

CLI 读取 `forge surfaces` 获取全部 surface key，输出：
- `install` = 所有 `<key>-install` 的聚合
- `ci` = 所有 `<key>-lint` + `<key>-compile` + `<key>-unit-test` 的聚合
- `clean` = 所有 `<key>-clean` 的聚合

### 删除：surface rule 文件 + server-lifecycle.md

CLI 命令内置所有 surface type 的差异化逻辑，以下文件删除：

| 文件 | 行数 | 去向 |
|------|------|------|
| `rules/server-lifecycle.md` | 745 | 迁移到 `forge justfile scaffold` Go 代码 |
| `rules/surfaces/api.md` | 67 | 迁移到 CLI 内部逻辑 |
| `rules/surfaces/cli.md` | 55 | 迁移到 CLI 内部逻辑 |
| `rules/surfaces/tui.md` | 55 | 迁移到 CLI 内部逻辑 |
| `rules/surfaces/web.md` | 69 | 迁移到 CLI 内部逻辑 |
| `rules/surfaces/mobile.md` | 72 | 迁移到 CLI 内部逻辑 |

### 精简：SKILL.md（548 行 → ~250 行）

保留的职责：流程编排、语言检测、Convention 加载、占位符填值、验证、输出确认。

删除的内容：

| 删除项 | 行数 | 理由 |
|--------|------|------|
| Step 1d Load Server Lifecycle Patterns | ~15 | CLI 处理 |
| Step 3b Surface recipe 生成细节 | ~70 | CLI scaffold 替代，缩为 ~10 行 |
| Step 3b Slot Placeholder Resolution 表 | ~15 | 简化为引用 CLI 输出的占位符 |
| Surface rule 加载逻辑（Step 0） | ~20 | 删除 rule 文件 |
| Phase 1 Consistency Verification | ~20 | CLI 是 trusted producer |
| Surface Gate Targets 段 | ~20 | 取消独立 gate recipe 概念 |
| EXTREMELY-IMPORTANT 重复项 | ~10 | 删除 CLI 已覆盖的条目 |
| Notes 重复 | ~25 | 去重后保留要点 |
| INLINE test-type-model | ~15 | 保留（唯一残留的声明性知识） |

**Gate recipe 概念取消**：不再区分"language-level recipe"和"gate recipe"。CLI 为每个 surface 生成完整的 recipe 集（含 compile / fmt / lint / unit-test），quality gate 直接调用 `<key>-compile` 等。多 surface 项目无全局 recipe。

### 保留

| 文件 | 说明 |
|------|------|
| `rules/self-correction.md` (34 行) | 错误模式表，仍由 agent 在 Phase 3 使用 |

### Agent 新流程

```
Step 0: forge surfaces → 获取 surface 列表（key + type + form）
Step 1: 检测语言/工具 → 加载 Convention → 记录 slot 值
Step 2: 检查已有 justfile（boundary markers + user-customized 保护）
Step 3: 逐 surface 调用 forge justfile scaffold --type X --key Y
        → 拿到 recipe 骨架 → 填 {{PLACEHOLDER}}
        → 最后调用 forge justfile scaffold --aggregate → 拿到 install/ci/clean
        → boundary marker merge 组装 justfile
Step 4: 验证（dry-run + actual execution，无 Phase 1 consistency check）
Step 5: 输出确认
```

### Recipe 命名统一模型

| 项目形态 | Recipe 命名 | 示例 |
|----------|------------|------|
| 单 surface scalar（`surfaces: cli`） | 无前缀 | `compile`, `test`, `teardown` |
| 单 surface named（`surfaces: {key: app, type: tui}`） | `<key>-` 前缀 | `app-compile`, `app-test`, `app-teardown` |
| 多 surface（`backend=api` + `frontend=web`） | `<key>-` 前缀 | `backend-dev`, `frontend-compile` |
| 多 surface 聚合 | 无前缀 | `install`, `ci`, `clean` |

### 向后兼容

硬切换。Forge v3.0.0 尚未发布，无存量用户。已有 justfile 中 `# user-customized` 标记的 recipe 会被保护机制保留，未标记的全局 recipe 在重新 `/init-justfile` 时被新结构替换。

## 减重效果

| 指标 | Before | After | 变化 |
|------|--------|-------|------|
| 总行数 | 1645 | ~284 | -83% |
| 文件数 | 8 | 2 | -75% |
| Agent 每次加载的 prompt token | ~1645 行 | ~284 行 | -83% |
| 新增 CLI 代码 | 0 | ~500 行 Go | prompt 层转移 |

## 风险与缓解

| 风险 | 影响 | 缓解 |
|------|------|------|
| CLI scaffold 命令有 bug 生成错误代码 | 所有使用 `/init-justfile` 的项目受影响 | Phase 2 dry-run + Phase 3 actual execution 验证；CLI 有单元测试 |
| 新增 surface type 需要改 CLI 代码而非 rule 文件 | 发布节奏依赖 forge CLI 发版 | CLI 发版频率高（RC 阶段），且 surface type 枚举稳定（5 种） |
| 占位符列表不完整，某些 recipe 需要额外参数 | Agent 无法正确填值 | CLI 文档化完整占位符清单；agent 遇到未知占位符时保留原样并报告 |

## 行动项

1. **实现 `forge justfile scaffold` 命令**（Go）：内置 5 种 surface type 的 recipe 模板、占位符机制、聚合生成
2. **重写 SKILL.md**：从 548 行精简到 ~250 行，删除 Step 1d / Phase 1 / Surface Gate Targets，简化 Step 3
3. **删除 6 个文件**：`server-lifecycle.md` + 5 个 surface rules
4. **更新 quality gate**：移除 fallback 链（`<key>-compile` 不存在时 fallback 到 `compile`），直接调 `<key>-compile`
