---
created: 2026-06-09
author: faner
status: Draft
---

# 提案：init-justfile 精简 — CLI scaffold 替代 prompt 层模板

## 成功标准

1. 生成的 justfile 对全部 5 种 surface type（cli/tui/api/web/mobile）语法正确，`just --list` 零错误
2. 所有下游 consumer（`run-tests` skill、quality gate 机制）在新 recipe 命名模型下功能不变
3. `forge justfile scaffold` Go 代码有单元测试覆盖（每个 surface type 至少 1 个 test case + 聚合模式 test case）
4. prompt 层（SKILL.md + 保留的 rules）总行数 < 300 行
5. 用户执行 `/init-justfile` 后得到的 justfile 与旧版行为等价（相同的 recipe、相同的占位符填值结果）

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
| `--key` | 否 | Surface key（named surface 必传，scalar 省略）。若对 scalar surface 传入 `--key`，CLI 报错退出 |

**输出**：带 `<<PLACEHOLDER>>` 占位符的 just recipe 代码，输出到 stdout。CLI 不需要知道任何项目细节。**占位符语法**：使用 `<<...>>` 而非 `{{...}}` 避免与 Go template `{{...}}` 和 justfile 变量语法 `{{var}}` 冲突。

**每个 surface type 生成的 recipes**：

| Surface Type | Lifecycle Recipes | Quality Recipes | Aggregate |
|--------------|------------------|----------------|-----------|
| `cli` | `test`, `teardown` | `compile`, `fmt`, `lint`, `unit-test` | 无 |
| `tui` | `test`, `teardown` | `compile`, `fmt`, `lint`, `unit-test` | 无 |
| `api` | `dev`, `probe`, `test`, `teardown`, `<key>` | `compile`, `fmt`, `lint`, `unit-test` | `<key>` (dev→probe→test→teardown) |
| `web` | `dev`, `probe`, `test`, `teardown`, `<key>` | `compile`, `fmt`, `lint`, `unit-test` | `<key>` (dev→probe→test→teardown) |
| `mobile` | `test-setup`, `dev`, `probe`, `test`, `teardown`, `<key>` | `compile`, `fmt`, `lint`, `unit-test` | `<key>` (test-setup→dev→probe→test→teardown) |

Named surface（有 key）：所有 recipe 名加 `<key>-` 前缀（如 `backend-dev`）。Scalar surface（无 key）：无前缀（如 `dev`）。

所有 recipe 均包含 `[linux]` 和 `[windows]` 双平台变体。**标记策略**：所有 lifecycle recipes（dev/probe/test/teardown/`<key>`）和 quality recipes（compile/fmt/lint/unit-test）均标记 `# user-customized`，聚合 recipes（install/ci/clean）不标记。标记 `# user-customized` 的 recipe 在后续 `/init-justfile` 调用时被保护机制保留，不会被覆盖。

**占位符清单**：

| 占位符 | 说明 | Agent 解析来源 |
|--------|------|---------------|
| `<<START_CMD>>` | 启动服务的命令 | `package.json` scripts / `go run` + entry point / Convention |
| `<<PORT>>` | 服务监听端口 | Convention / 常见默认值（Node 3000, Go 8080） |
| `<<HEALTH_URL>>` | 健康检查 URL | `http://localhost:<<PORT>>/healthz`（默认） |
| `<<COMPILE_CMD>>` | 编译命令 | `go build ./...` / `npm run build` / Convention |
| `<<UNIT_TEST_CMD>>` | 单元测试命令 | `go test -v ./...` / `npm test` / Convention |
| `<<LINT_CMD>>` | 静态分析命令 | `golangci-lint run` / `npm run lint` / Convention |
| `<<FMT_CMD>>` | 格式化命令 | `gofmt -w .` / `npm run fmt` / Convention |
| `<<BUILD_CMD>>` | 构建命令 | `go build -o bin/app` / `npm run build` / Convention |
| `<<CLEAN_CMD>>` | 清理命令 | `go clean` / `rm -rf dist` / Convention |
| `<<INSTALL_CMD>>` | 安装依赖命令 | `go mod download` / `npm install` / Convention |
| `<<TEST_CMD>>` | Surface-level 测试命令 | Convention test runner + file pattern |
| `<<URL_KEY>>` | 服务标识键名（用于 PID 文件命名） | Surface key（与 `--key` 参数一致） |
| `<<SERVICE_LIST>>` | 多服务编排时的服务启动依赖列表 | Convention multi-service 定义 |

**聚合 recipe 生成**（`--aggregate` 模式）：

CLI 读取 `forge surfaces` 获取全部 surface key，输出：
- `install` = 所有 `<key>-install` 的聚合
- `ci` = 所有 `<key>-lint` + `<key>-compile` + `<key>-unit-test` 的聚合（不含 surface-level test recipe，因为 surface test 需要运行时环境，属于 `test-setup` + `<key>-test` 编排范畴，不应混入 CI 流水线）
- `clean` = 所有 `<key>-clean` 的聚合

**多服务编排模式**：当项目存在多个需要运行时依赖的 surface（如 api + web），`--aggregate` 模式额外生成 `test-setup` 聚合 recipe，按依赖顺序编排各 surface 的启动（先 api 后 web）和 teardown（逆序）。此 recipe 仅在 `forge surfaces` 返回多个 service-type surface（api/web/mobile）时生成，不适用于纯 cli/tui 组合。

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

**Convention Cold Start Fallback**：当 Convention 文件不存在时（新项目或未配置），agent 按以下策略回退占位符值：优先从项目文件推断（`package.json` scripts、`Makefile`、go module），推断失败则使用语言默认值（Node: 3000/Go: 8080），无法推断的占位符保留 `<<PLACEHOLDER>>` 原样并在最终报告中列出。此策略以 5-10 行摘要保留在 SKILL.md 中。

删除的内容：

| 删除项 | 行数 | 理由 |
|--------|------|------|
| Step 1d Load Server Lifecycle Patterns | ~15 | CLI 处理 |
| Step 3b Surface recipe 生成细节 | ~70 | CLI scaffold 替代，缩为 ~10 行 |
| Step 3b Slot Placeholder Resolution 表 | ~15 | 简化为引用 CLI 输出的占位符 |
| Surface rule 加载逻辑（Step 0） | ~20 | 删除 rule 文件 |
| Phase 1 Consistency Verification | ~20 | 简化为轻量级 recipe 完整性断言（`just --list` 验证），替代完整的一致性检查 |
| Surface Gate Targets 段 | ~20 | 取消独立 gate recipe 概念 |
| EXTREMELY-IMPORTANT 重复项 | ~10 | 删除 CLI 已覆盖的条目 |
| Notes 重复 | ~25 | 去重后保留要点 |
| INLINE test-type-model | ~15 | 保留（唯一残留的声明性知识） |

**Gate recipe 概念取消**：不再区分"language-level recipe"和"gate recipe"。CLI 为每个 surface 生成完整的 recipe 集（含 compile / fmt / lint / unit-test），quality gate 直接调用 `<key>-compile` 等。多 surface 项目无全局 recipe。

### Consumer Impact

取消 gate recipe 对下游 consumer 的影响：

| Consumer | 变更前 | 变更后 |
|----------|--------|--------|
| `run-tests` skill | 解析 `gate-compile` / `gate-lint` 等 gate recipe 名 | 直接调用 `<key>-compile` / `<key>-lint` 等 surface recipe，移除 fallback 链（`<key>-compile` 不存在 → fallback `compile`） |
| Quality gate 机制 | 查找全局 `compile` / `lint` recipe | 按 surface key 查找 `<key>-compile` / `<key>-lint`，单 surface scalar 项目无前缀 |
| `forge quality-gate` Go binary | 硬编码调用 `just compile` / `just unit-test` | 需更新为 `just <key>-compile` / `just <key>-unit-test`（见行动项） |
| `ci` 聚合 recipe | 不存在 | 新增，聚合所有 surface 的 lint + compile + unit-test |

### 保留

| 文件 | 说明 |
|------|------|
| `rules/self-correction.md` (34 行) | 错误模式表，仍由 agent 在 Phase 3 使用 |

### Agent 新流程

**用户可见行为**：开发者执行 `/init-justfile` 时，体验与旧版一致——最终得到一个完整的 justfile，包含所有 surface 的 recipe。区别在于生成速度更快（CLI 秒级输出骨架 vs. agent 逐行生成），且 customization workflow 不变（`# user-customized` 标记的 recipe 被保留）。

**错误场景覆盖**：
- `forge surfaces` 返回未知 surface type → agent 跳过该 surface 并在报告中警告
- CLI scaffold stdout 解析失败（非预期输出格式）→ agent 报错并中止，提示运行 `forge justfile scaffold` 手动排查
- 已有 justfile 中存在 boundary marker 外的手动 recipe → 保护机制保留这些 recipe，仅替换 marker 内区域

```
Step 0: forge surfaces → 获取 surface 列表（key + type + form）
Step 1: 检测语言/工具 → 加载 Convention → 记录 slot 值
Step 2: 检查已有 justfile（boundary markers + user-customized 保护）
Step 3: 逐 surface 调用 forge justfile scaffold --type X --key Y
        → 拿到 recipe 骨架 → 填 <<PLACEHOLDER>>
        → 最后调用 forge justfile scaffold --aggregate → 拿到 install/ci/clean
        → boundary marker merge 组装 justfile
Step 4: 验证
        → 4a: 轻量级完整性断言（just --list 确认所有 recipe 可被 just 解析）
        → 4b: dry-run + actual execution 验证
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
| Prompt 层行数（SKILL.md + rules） | 1645 | ~284 | -83% |
| Prompt 层文件数 | 8 | 2 | -75% |
| Agent 每次加载的 prompt token | ~1645 行 | ~284 行 | -83% |
| 新增 CLI 代码 | 0 | ~600-1000 行 Go | 从 prompt 层迁移 |

**注**：-83% 为 prompt 层缩减。系统总代码量（prompt + Go）从 1645 行变为 ~884-1284 行，净减约 22-46%。核心收益不在总行数，而在将确定性 bash 模板从 LLM 推理域移入程序生成域，降低 token 成本和 agent 出错率。

## 风险与缓解

| 风险 | 可能性 | 影响 | 缓解 |
|------|--------|------|------|
| CLI scaffold 命令有 bug 生成错误代码 | M | H | Phase 2 dry-run + Phase 3 actual execution 验证；CLI 有单元测试；回滚开关（见下） |
| 新增 surface type 需要改 CLI 代码而非 rule 文件 | L | M | CLI 发版频率高（RC 阶段），且 surface type 枚举稳定（5 种） |
| 占位符列表不完整，某些 recipe 需要额外参数 | M | L | CLI 文档化完整占位符清单；agent 遇到未知占位符时保留原样并报告 |
| Agent 填错占位符（如 PORT 填了字符串） | M | M | SKILL.md 保留占位符值校验逻辑（类型/格式断言） |
| Go template 维护成本高于 markdown rule | L | L | Surface type 逻辑稳定后几乎不变；Go 有类型安全和编译检查 |
| 调试难度：Go 源码 vs. markdown rule 文件 | L | M | CLI 输出到 stdout，可直接肉眼检查；保留 `--debug` 标志输出渲染中间态 |

**回滚机制**：在 `config.yaml` 中新增 `forge.justfile.useScaffold` 开关（默认 `true`）。设为 `false` 时，agent 回退到旧版 prompt 层模板生成流程。旧版 rule 文件在 v3.0.0 发布前保留在代码库中（标记为 deprecated），不随新版本打包。

## 行动项

1. **实现 `forge justfile scaffold` 命令**（Go，~600-1000 行）：内置 5 种 surface type 的 recipe 模板、占位符机制、聚合生成
   - scaffold 生成器（~300 行）：surface type → recipe 模板映射、占位符注入、boundary marker 包装
   - 聚合生成器（~100 行）：读取 `forge surfaces` 输出，按依赖顺序生成 install/ci/clean
   - 参数校验（~50 行）：surface type 白名单、scalar/named 的 --key 校验
   - 单元测试（~200-500 行）：每个 surface type 1 个 test case + 聚合模式 + 边界场景
2. **重写 SKILL.md**：从 548 行精简到 ~250 行，删除 Step 1d / Phase 1 / Surface Gate Targets，简化 Step 3
3. **删除 6 个文件**：`server-lifecycle.md` + 5 个 surface rules
4. **更新 quality gate**：移除 fallback 链（`<key>-compile` 不存在时 fallback 到 `compile`），直接调 `<key>-compile`
5. **更新 `forge quality-gate` Go binary**：将硬编码的 `just compile` / `just unit-test` 改为按 surface key 拼接 recipe 名

## 范围外

以下内容明确 **不在本提案范围内**：

- **`forge surfaces` 命令**：不修改其输出格式或逻辑
- **Convention 系统**：不改变 Convention 文件结构或加载机制
- **其他 skill recipe 调用方**：`fix-bug.md`、`clean-code/SKILL.md` 等 skill 中对 just recipe 的调用不做修改
- **`forge quality-gate` Go binary 的 recipe 名解析逻辑重构**：本次仅做最小改动（硬编码替换），完整的动态解析机制留给后续 PR
- **justfile 语法版本升级**：不改变生成的 justfile 的语法兼容性要求

## 替代方案与行业基准

| 方案 | 描述 | 为何不采用 |
|------|------|-----------|
| **Do Nothing**（保持现状） | 维持 1645 行 prompt 层模板，依赖 agent 理解和复用 | Token 浪费严重，维护分散，职责错位未解决 |
| **Yeoman / Plop / Cookiecutter** | 使用通用脚手架工具，在项目初始化时生成 justfile | 这些工具面向整个项目骨架（目录结构+配置+代码），不适用于"运行时按需生成单一配置文件"的场景；引入额外 Node.js/Python 运行时依赖 |
| **Hygen** | 基于 template 文件的代码生成器 | 需要维护独立模板文件，且与 Forge CLI 集成需要额外胶水代码，不如直接在 Go 中实现轻量 |
| **本方案** | Forge CLI 内置 scaffold 子命令，agent 调用 CLI + 填值 | 选择理由：零外部依赖，与 Forge CLI 自然集成，surface type 枚举稳定适合硬编码，占位符机制将确定性逻辑与 LLM 推理干净分离 |
