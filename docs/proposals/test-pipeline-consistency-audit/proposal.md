---
name: test-pipeline-consistency-audit
status: Draft
created: 2026-05-27
supersedes: surface-test-type-model (recipe 命名部分：`test-<surface-type>-<scope>` → `<surface-key>-test`)
---

# 测试 Pipeline 一致性审计与修复

## Problem

Forge 测试 pipeline 经历了多次架构演进（Profile → Convention、staging/graduation → tag-based promotion、统一 e2e → Surface 区分的 Functional/E2E），但 Go 代码层未同步更新，导致代码与文档之间存在系统性术语和路径断裂。当前 state：

- **Go 代码层**仍使用旧模型路径 `tests/e2e/`、旧术语 `E2E`/`graduation`/`staging`/`profile`
- **Skill 文档层**已迁移到新模型 `tests/<journey>/`、tag-based promotion、Convention 驱动
- 两层之间的不一致会误导 AI agent 生成错误的测试代码或路径

**Evidence**: `forge-cli/pkg/feature/constants.go` 定义 `E2ETestsBaseDir = "tests/e2e"`，而 `gen-test-scripts/SKILL.md` 明确声明 "Tests go directly to `tests/<journey>/`, NOT to `tests/e2e/features/`"。`ARCHITECTURE.md` Quick 模式图包含 gen-contracts + gen-scripts，但代码 `GetQuickTestTasks()` 跳过两者。Go build tag `//go:build e2e` 在 Forge 自身（CLI 项目）的测试中使用了与其 surface 类型不匹配的标签名。

**Urgency**: v3.0.0 是重构大版本，发布前统一术语、路径和 build tag 可避免技术债累积。当前不一致已在影响 `quality_gate.go` 的错误路径输出和 `mobile-test-setup` 的缺失调用。

## Solution

分两层系统性修复所有活跃代码和文档中的不一致，使 Go 代码层、Skill 文档层、ARCHITECTURE 文档达成术语、路径和 build tag 统一。

### 新目录结构

物理路径从 `tests/e2e/` 扁平化为 `tests/`，消除与 surface 类型无关的中间目录层：

```
旧结构:                          新结构:
tests/                           tests/
  e2e/                             config.yaml          (探针配置)
    config.yaml                    results/
    results/                         raw-output.txt
      raw-output.txt              <journey>/            (测试脚本，不变)
    features/     (staging，删除)
    .graduated/   (graduation，删除)
  <journey>/
```

常量映射：`E2ETestsBaseDir = "tests/e2e"` → `TestBaseDir = "tests"`，`E2EStagingDir` / `E2EGraduatedDir` 删除。

### Layer 1: Go 代码层术语/路径统一

1. **constants.go + paths.go 重命名**: `E2E*` 常量 → surface-neutral 名称（`TestBaseDir = "tests"`、`GetTestResultsDir()`、`GetTestConfigPath()`），删除 `E2EStagingDir`/`E2EGraduatedDir` 和对应的 `GetE2EStagingDir()`/`GetE2EGraduatedMarker()`/`GetE2ETargetDir()` 函数
2. **testrunner 路径更新**: MR #173 已修复 `WriteRegressionRawOutput` 路径为 `tests/results/`；需确认 `WriteUnitTestRawOutput` 路径一致性
3. **e2eprobe 包重命名**: `pkg/e2eprobe` → `pkg/serverprobe`，探针配置路径从 `tests/e2e/config.yaml` → `tests/config.yaml`
4. **quality_gate 更新**: `runTestRegressionSurface` 第 286 行仍有 `tests/e2e/results/raw-output.txt`（MR #173 只修了 legacy 路径，surface 路径遗漏）；`GetE2EStagingDir`/`GetE2EGraduatedMarker` 调用替换；`mobile-test-setup` 集成到 `runSurfaceLifecycle()`；注释 "promoted scripts in tests/e2e/" 更新
5. ~~**init.go 路径更新**~~: MR #173 已完全修复，移除此项
6. **废弃模板和类型清理**: 删除 `pkg/prompt/data/test-verify-regression.md` 和 `pkg/task/data/test-verify-regression.md` 两处模板；删除 `infer.go` 中 `T-quick-verify-regression` 死代码；清理 `types.go` 中 `TypeTestVerifyRegression` 常量及其 `ValidTypes`/`SystemTypes` 注册项（描述 "after graduation" 一并移除）；修复 `pkg/template/data/coding.fix.md` 中 "E2E Fix Boundaries" section 和 `just test-e2e` 引用；更新 `docsync_test.go` 中对旧常量名的引用
7. **Go build tag 重命名为 surface-type-specific 标签**: `//go:build e2e` → `//go:build <surface-type>-<type>`（如 CLI 项目用 `cli-functional`、Web 项目用 `web-e2e`），与 `test-type-model.md` 的 tag 命名完全对齐。变更传播链：Convention 文件（定义 build tag 规范）→ init-justfile surface rules（recipe 模板中 `-tags=<surface-type>-<type>`）→ 生成的 justfile（`<surface-key>-test` recipe 内部用 `-tags=<surface-type>-<type>`）→ run-tests（调用 `just <surface-key>-test <journey>`）。覆盖：`tests/` 目录下所有测试文件、Forge 项目自身 justfile、Convention 文件（`go.md`、`ginkgo.md`、`vitest.md`、`pytest.md`、`junit.md`、`rust.md`、`index.md` 共 7 个）、6 个 init-justfile surface rule recipe 模板、`test-guide/rules/` 中 build tag 表格
8. **prompt/task 模板术语更新**: `pkg/prompt/data/test-gen-scripts.md`、`pkg/prompt/data/test-run.md`、`pkg/task/data/test-gen-scripts.md` 中的 "profile"/"active profile"/"profile resolution" 替换为 "Convention"/"surface" 术语

### Layer 2: Skill 文档层术语统一

9. **gen-contracts + gen-test-scripts 术语替换**: "e2e 测试管道" → "Forge 测试管道"；`tests/e2e/` 旧路径引用 → `tests/<journey>/`；`gen-test-scripts/rules/step-1-contract-loading.md` 中 `tests/e2e/step1_test.go` 示例路径更新；`gen-test-scripts/rules/convention-guide.md` 中 "e2e tests" 通用引用替换；`gen-contracts/rules/journey-contract-model.md` 第 159 行 "language profile" 属于旧模型对比表，保留不改动
10. **breakdown-tasks + quick-tasks 执行顺序修正**: `quick-tasks/SKILL.md` 执行顺序从 `gen-journeys → gen-contracts → gen-test-scripts` 修正为 `gen-journeys → run-test`（`breakdown-tasks/SKILL.md` 的顺序正确，不动）；两处 "Integration Test Impact Assessment" → "Test Impact Assessment"
11. **其他 Skill/Command 文档修正**: `commands/fix-bug.md` 路径更新（`tests/e2e/features/` → `tests/<journey>/`）；`commands/run-tasks.md` 中 `T-test-verify-regression` 和 "e2e verification" 引用清理；test-guide 术语修正（含 `rules/draft-generation.md` 和 `rules/pattern-extraction.md` 中 build tag 表格）；`submit-task/data/record-format-test.md` 删除 `test.verify-regression` 类型并更新 `tests/e2e/` 示例路径；gen-sitemap 配置文件重命名（`e2e-config.yaml` → `test-config.yaml`）及 SKILL.md 中路径引用更新；`consolidate-specs/SKILL.md` 第 22 行 "e2e tests are promoted" → "all tests pass"；`init-justfile/SKILL.md` 中 `tests/e2e/` 示例路径更新；`init-justfile/templates/` 下全部 6 个 justfile 模板（`python.just`、`rust.just`、`node.just`、`go.just`、`mixed.just`、`generic.just`）中 `tests/e2e/` 路径更新；`run-tests/rules/test-isolation.md` 中 4 处 `tests/e2e/` 路径引用更新
12. **ARCHITECTURE.md 修复**: Quick 模式流程图移除 gen-contracts/gen-scripts；任务 ID 从 `T-test-1~5` 更新为描述性 ID；移除 `T-test-promote` 幽灵条目；"profile type" → "Convention" / "surface type"（第 275、305 行）；"profile 路由" → "Convention 路由"（第 516 行）；第 305 行 "所有 profile type 的 gen-journeys 并行执行后汇聚到 gen-contracts" 按实际代码修正
13. **OVERVIEW/WORKFLOW 文档同步（含中文版）**: 替换所有 "e2e" 泛用（约 15+10 处）为正确的 surface-specific 术语；更新 graduation/staging 描述为 tag-based promotion；移除 "profile" 旧术语引用；`OVERVIEW.zh.md` 和 `WORKFLOW.zh.md` 同步更新
14. **surface-test-type-model 提案 recipe 命名部分更新**: 该已批准提案第 73 行和第 85 行的 recipe 命名 `test-<surface-type>-<scope>`（如 `test-cli-functional`）被 `<surface-key>-test`（如 `cli-test`）supersede；第 107 行的多 surface recipe 命名同步更新；第 73 行的 alias 过渡方案不再适用（v3.0.0 大版本允许破坏性变更，alias 直接删除而非 2 版本过渡期保留，该提案 NFR1 的向后兼容要求被本提案覆盖）

## Alternatives

### A. 仅修复高层级问题（Critical + High）

只修复 Go 代码层和核心文档，Medium/Low 的 Skill 文档不一致留到后续。优点：更快完成，风险更低。缺点：Skill 层术语不一致继续误导 AI agent。

### B. 不修复，仅输出审计报告

生成完整审计报告作为参考文档，不做代码修改。优点：零风险。缺点：问题持续存在，后续开发继续受影响。

## Scope

### In Scope

- Go 代码层（`pkg/feature/`、`pkg/testrunner/`、`pkg/e2eprobe/`、`internal/cmd/`、`pkg/prompt/data/`、`pkg/task/data/`、`pkg/template/data/`）的术语、常量、路径、函数名统一
- 物理路径 `tests/e2e/` 扁平化为 `tests/`（config.yaml、results/ 上提一级，删除 staging/graduated 目录）。注：MR #173 已修复 `init.go` 和 `testrunner` 路径，但 `tests/e2e/results/` 空目录仍存在，`e2eprobe`/`quality_gate` 等仍引用旧路径
- Skill 文档层（`plugins/forge/skills/`）和 Command 文档（`plugins/forge/commands/`）的术语和路径引用统一
- ARCHITECTURE.md、OVERVIEW.md、OVERVIEW.zh.md、WORKFLOW.md、WORKFLOW.zh.md 文档与代码对齐
- 废弃代码清理（死代码、未使用的模板和类型注册项）
- `mobile-test-setup` 在 `quality_gate.go` 中的集成
- Go 测试文件的同步更新（反映重命名后的包名、常量、build tag），包括 `docsync_test.go` 中对旧常量名的引用
- Go build tag 从 `e2e` 重命名为 surface-type-specific 标签
- `validate-ux` 依赖链修复（`T-validate-ux` 应依赖最后一个 `run-test`）
- `TypeTestVerifyRegression` 类型系统条目完整清理（常量、描述、注册项）
- deprecated alias 直接删除
- `surface-test-type-model` 提案 recipe 命名部分被本提案 supersede

### Out of Scope

- 历史特征任务文件（`docs/features/*/tasks/quick-graduate-*.md` 等）不修改
- 已 Superseded 的 proposal 文档不修改（`surface-test-type-model` 除外，仅更新 recipe 命名部分）
- Legacy 测试 fixture 中的旧术语不修改
- 配置迁移代码中的 `e2eTest` 旧键名（向后兼容层，仅添加移除目标版本注释）

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Go 包重命名（e2eprobe）导致外部引用断裂 | Medium | Medium | 全面 grep 检查所有 import；CI 编译验证 |
| Go build tag 重命名为 surface-specific 标签 | Medium | Medium | Convention 文件 + init-justfile + justfile 三处同步更新；grep 验证无遗漏 |
| 物理路径变更影响用户项目 | Medium | High | v3.0.0 大版本允许破坏性变更；`forge init` 重建目录结构 |
| `surface-test-type-model` 已批准提案被部分 supersede | Low | Medium | 仅 supersede recipe 命名部分（`test-<surface-type>-<scope>` → `<surface-key>-test`），测试类型映射和术语定义不变 |
| `TypeTestVerifyRegression` 删除后影响依赖它的代码 | Low | Medium | 全面 grep 确认无运行时引用（该类型从未被 autogen 生成） |
| Skill 文档修改后 AI agent 行为变化 | Medium | Low | 术语统一方向与现有 HARD-RULE 一致，是强化而非改变行为 |
| ARCHITECTURE.md 更新引入新错误 | Low | Medium | 以 autogen.go 代码为 ground truth，交叉验证每个流程图 |

## Success Criteria

- [ ] `grep -rn "tests/e2e" forge-cli/pkg/ forge-cli/internal/`（排除 _test.go 和 worktrees）返回 0 结果
- [ ] `grep -rn "test-e2e\|E2E Fix" forge-cli/pkg/template/data/` 返回 0 结果
- [ ] `grep -rn '"e2e"' forge-cli/pkg/ forge-cli/internal/`（排除 _test.go、worktrees、config 迁移代码）返回 0 结果
- [ ] `grep -rn "graduated\|graduation" forge-cli/pkg/ forge-cli/internal/ --include="*.go"`（排除 _test.go、worktrees、config 迁移代码 `forgeconfig/config.go`）返回 0 结果
- [ ] `grep -rn "staging" forge-cli/pkg/ forge-cli/internal/ --include="*.go"`（排除 _test.go、worktrees）返回 0 结果
- [ ] `grep -rn "profile" forge-cli/pkg/prompt/data/ forge-cli/pkg/task/data/ --include="*.md"` 返回 0 结果
- [ ] `e2eprobe` 包已重命名为 `serverprobe`，`grep -rn "e2eprobe" forge-cli/ --include="*.go"` 返回 0 结果
- [ ] `quality_gate.go` 中 `runSurfaceLifecycle()` 包含 `mobile-test-setup` 调用
- [ ] `autogen.go` 中 `T-validate-ux` 依赖最后一个 `run-test` 任务（与 `T-validate-code` 同级）
- [ ] `types.go` 中不存在 `TypeTestVerifyRegression` 常量及其 ValidTypes/SystemTypes 注册项
- [ ] ARCHITECTURE.md Quick 模式流程图：不含 gen-contracts/gen-scripts；任务 ID 为描述性名称；不含 T-test-promote；"profile type" 已替换为 surface/Convention 术语
- [ ] Skill 文档中 "integration test"（`breakdown-tasks`、`quick-tasks` 的 "Integration Test Impact Assessment"）已替换为 "Test Impact Assessment"；`gen-test-scripts/types/ui.md` 中的 "Integration Test" 保留（UI 组件集成测试的专门概念）
- [ ] `infer.go` 中 `T-quick-verify-regression` 死代码已移除
- [ ] `grep -rn "//go:build e2e" tests/ forge-cli/` 返回 0 结果
- [ ] `grep -rn '\-tags=e2e' justfile plugins/forge/` 返回 0 结果
- [ ] Convention 文件中 build tag 与 surface 类型对齐（`go.md`/`ginkgo.md`/`index.md` 等 7 个文件中 `tags=e2e` 替换为 surface-specific 值）
- [ ] `surface-test-type-model/proposal.md` 第 73、85 行 recipe 命名已更新为 `<surface-key>-test`，NFR1 向后兼容要求标记为 v3.0.0 已覆盖
- [ ] 所有 deprecated alias（`test-e2e → <surface>-test`）已删除，`grep -rn "alias test-e2e" plugins/forge/` 返回 0 结果
- [ ] 物理目录 `tests/e2e/` 不再存在（`tests/config.yaml`、`tests/results/` 直接在 `tests/` 下）
- [ ] `go build ./...` 和 `go test ./...` 全部通过
- [ ] `grep -rn "tests/e2e" forge-cli/docs/OVERVIEW.zh.md forge-cli/docs/WORKFLOW.zh.md` 返回 0 结果（中文版文档同步更新）
