---
created: 2026-05-24
author: "faner"
status: Draft
---

# Proposal: Forge CLI 代码清理

## Problem

forge-cli（Go 代码库，92 个源文件）积累了大量技术债务：死代码、重复逻辑、超大文件和反模式，导致可维护性下降。

### Evidence

代码审计发现 15 个具体问题：

**死代码（2 项）**：
- `run.go` 中 `GetVersion()`、`GetName()`、`IsTestMode()` 无任何调用者
- 仓库根目录残留构建产物（cmd.out, cout.out, coverage.out, just.out）

**待迁移废弃代码（1 项）**：
- `pkg/feature/feature.go` 中 `SetFeature()` 标记 Deprecated，仍有 7+ 处调用需要迁移

**重复逻辑（6 项）**：
- `cmd/output.go` 重复定义了 `base.Debugf`（`quality_gate.go` 有 10+ 处通过 `cmd.Debugf()` 调用，冗余的间接层需清理）
- YAML frontmatter 解析在 4+ 个文件中独立实现，签名各异
- 依赖检查逻辑（含 `.x` 通配符）在 4 个文件中重复
- `defaultRunClaude()` 在 `claude.go` 和 `worktree.go` 中完全相同
- 3+1 个 `mapXxxToSlugLens` 函数（3 个独立函数 + 1 个变体）做同样的事，可用泛型替代
- `errors.go` 和 `output.go` 将 `base/` 的所有符号重导出到父包（包括 `Debugf`，被 `quality_gate.go` 等 10+ 处引用）

**超大文件（2 项）**：
- `forensic.go` 994 行（20+ 结构体定义、3 个命令、300 行嵌套函数）
- `worktree.go` 1069 行（6 个命令 + 补全 + TUI + 文件操作）

**反模式（4 项）**：
- `validateRecordData()` 内部调用 `os.Exit()`，使函数不可测试
- `askAutoBehavior()` 130 行，13 个相同的 `askConfirm` 块
- `runQualityGate()` 中 e2e 回归逻辑 4 层嵌套
- `testbridge.go` 导出 37 个内部符号供测试使用，但文件在正式构建中

### Urgency

代码库处于 v3.0.0-rc.19，即将正式发布。在 API 稳定前做结构清理是最佳时机，发布后重构成本将显著增加。

## Proposed Solution

按自底向上顺序执行四阶段纯重构：死代码消除 → 超大文件拆分 → 重复逻辑合并 → 反模式修复。每个阶段完成后确保所有测试通过。

### Innovation Highlights

无创新。这是标准的代码健康维护——Go 社区推荐的渐进式重构实践。

## Requirements Analysis

### Key Scenarios

- 开发者在 forensic 或 worktree 命令中定位问题：文件拆分后每个文件职责单一，更容易导航
- 新增命令需要依赖检查：统一到单一函数后只需调用一处
- CI 运行测试：所有 108 个测试文件必须继续通过，零行为变更

### Non-Functional Requirements

- **向后兼容**：所有 CLI 命令的输入输出保持不变
- **构建稳定**：`go build ./...` 和 `go test ./...` 全部通过
- **代码量缩减**：消除重复后总行数应减少
- **性能不退化**：构建和测试执行时间不退化（基准：当前 `go test ./...` 耗时）

### Constraints & Dependencies

- Go 1.25 工具链
- 纯重构：不引入新依赖、不改变外部行为
- 需要阅读 `docs/conventions/forge-distribution.md` 了解分发约束

## Alternatives & Industry Benchmarking

### Industry Solutions

Go 社区标准做法：golangci-lint 发现问题 + 手动重构修复。结构性问题（大文件拆分、重复逻辑合并）无法自动化。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 技术债继续累积，发布后重构更贵 | Rejected: 成本只会增加 |
| 只做 lint 修复 | golangci-lint | 自动化，~0.5 天 | 仅覆盖 15 项中的 ~5 项（死代码 + 简单模式），无法解决结构性问题（重复、大文件、os.Exit 滥用） | Rejected: 覆盖面不足（~33%） |
| golangci-lint + 选择性重构 | Go 社区实践 | 覆盖全部 15 项，风险可控 | 需人工审查，预估 2-3 天 | **Selected: 平衡效果与风险** |

## Feasibility Assessment

### Technical Feasibility

Go 工具链原生支持重构（gofmt, go vet, go test）。所有修改都是文件内或包内操作，无跨模块依赖。

### Resource & Timeline

工作量约 15 个独立任务，每个任务可独立验证。

### Dependency Readiness

无外部依赖。所有修改在本地完成。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "re-export 层是必要的" | 5 Whys | Overturned: 子包已直接 import base/，重导出层增加混淆而非简化。清理时需将调用点改为直接引用 |
| "testbridge.go 是合理的测试模式" | XY Detection | Refined: 目标是测试可访问性。方案：底层函数迁移到 `pkg/task/`，testbridge 保留为薄别名，避免修改测试调用点 |
| "SetFeature() 可以先 Deprecated 再慢慢迁移" | Assumption Flip | Overturned: Deprecated 超过一个版本仍未迁移，仍有 7+ 处调用，需主动迁移后删除 |

## Scope

### In Scope

**Phase 1: 死代码消除**
- 删除 `run.go` 中 `GetVersion()`、`GetName()`、`IsTestMode()`（零调用者）
- 将构建产物添加到 `.gitignore`

**Phase 2: 超大文件拆分**
- 拆分 `forensic.go`（994 行）为 `types.go`、`search.go`、`extract.go`、`subagents.go`
- 拆分 `worktree.go`（1069 行）为按命令分文件 + `helpers.go`

**Phase 3: 重复逻辑合并**
- 扩展现有 `pkg/task.ParseFrontmatter()` 为共享解析器（两层 API：ParseFrontmatter 返回 raw YAML bytes 作为第二返回值，调用者按需 unmarshal 到各自的结构体）
- 统一依赖检查逻辑到 `pkg/task/` 单一函数
- 提取 `defaultRunClaude()` 到共享位置
- 完成 `SetFeature()` 迁移（7+ 调用点）并删除废弃函数
- 清理重导出层（`errors.go`、`output.go`），将 `cmd.Debugf` 调用点改为直接引用 `base.Debugf`

**Phase 4: 反模式修复**
- 重构 `askAutoBehavior()` 为数据驱动循环
- 修复 `validateRecordData()` 中的 `os.Exit` → 返回 error（scope 包含 `submit_test.go` 中 30+ 个测试用例的重构）
- 提取 `runE2ERegression()` 减少嵌套
- 统一错误处理模式（非顶层函数用 `return error`；`quality_gate.go` 中 RunE 处理器内的 2 处 `os.Exit(0)` 视为顶层调用，保留不动）
- testbridge 清理：将底层函数（如 `getTaskPhase()`）迁移到 `pkg/task/`，testbridge 保留为薄别名层，确保测试调用点无需修改

### Out of Scope

- 用泛型替代 `mapXxxToSlugLens` 函数：净节省约 4 行，泛型在发布前稳定阶段增加认知负担，收益不足以支撑风险
- 新功能或行为变更
- API 接口变更
- 性能优化
- 新增测试（保持现有测试通过即可；`submit_test.go` 的 30+ 用例重构属于 `validateRecordData` 签名变更的必要适配，不算新增测试）
- Go 依赖升级
- `quality_gate.go` 中 2 处 `os.Exit(0)` 调用（属于顶层 RunE 处理器，视为合法的 CLI 入口退出）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 拆分文件时遗漏导出符号 | M | M | 每步运行 `go build ./...` 验证编译 |
| 重导出层清理导致子包编译失败 | L | H | 先用 `go vet` 确认所有调用点已直接引用 base/ |
| os.Exit 移除改变错误处理流程 | L | M | 只改内部函数，顶层 RunE 保持 Exit 行为 |
| testbridge 重构影响测试 | M | L | 保持导出接口不变，仅重新组织 |
| 审计数据错误导致误删活跃代码 | M | H | 执行前用 `grep -r` 或 LSP findReferences 逐一验证每个删除目标的调用点数量 |

## Rollback Strategy

每个 Phase 作为一个独立的 git commit 提交，commit message 标注 phase 编号。若某 Phase 出现问题：

1. **Phase 1 回滚**：`git revert <commit>` 即可恢复，因为只涉及删除操作
2. **Phase 2 回滚**：`git revert <commit>` 恢复文件拆分，所有变更在包内完成无跨包影响
3. **Phase 3 回滚**：Phase 3 包含 5 个独立任务，默认按任务粒度提交（每个任务一个 commit），回滚时仅 revert 出问题的任务 commit
4. **Phase 4 回滚**：`git revert <commit>` 恢复 os.Exit 调用和 testbridge 原状

若 Phase 内部某个任务失败，可在 Phase commit 之前创建更细粒度的 checkpoint commit。原则：任何 commit 都必须通过 `go build ./...` 和 `go test ./...`。Phase 3 因任务数量较多（5 个），强制按任务粒度提交。

## Success Criteria

- [ ] `go build ./...` 零错误零警告
- [ ] `go test ./...` 全部通过（108 个测试文件）
- [ ] 0 个死代码函数（golangci-lint deadcode 检查通过）
- [ ] 0 处重复的 YAML frontmatter 解析（所有调用者使用 `pkg/task.ParseFrontmatter()` 或其提取的 bytes）
- [ ] 0 处重复的依赖检查逻辑
- [ ] forensic.go 从 994 行降至 <300 行（按功能拆分）
- [ ] worktree.go 从 1069 行降至 <300 行（按命令拆分）
- [ ] 0 处非顶层函数中的 `os.Exit` 调用（顶层定义：直接作为 cobra.Command 的 RunE/Erorr 处理器执行的函数）
- [ ] `askAutoBehavior()` 从 130 行降至 <30 行（数据驱动循环）
- [ ] `defaultRunClaude()` 仅存在 1 处定义（当前 2 处完全相同）
- [ ] testbridge.go 仅保留薄别名，底层函数全部在 `pkg/task/` 中有实体定义
- [ ] `SetFeature()` 迁移完成：0 处 Deprecated 调用点
- [ ] 总行数减少 >= 5%（消除重复 + 删除死代码）

## Next Steps

- Proceed to `/quick-tasks` to generate task breakdown
