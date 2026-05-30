---
created: "2026-05-30"
author: "faner"
status: Draft
intent: "refactor"
---

# Proposal: Forge CLI 代码库重组与规范建立

## Problem

Forge CLI 代码库（`forge-cli/`）缺少全面的编码规范，导致代码风格不一致、魔法值散布、死代码残留、包结构无明确组织原则。当前 `docs/conventions/` 下的规范仅覆盖枚举常量、扁平控制流和错误输出格式三个局部领域，无法指导新代码的编写和现有代码的演进。

### Evidence

**魔法值散布**：
- `"tests/results/raw-output.txt"` 在 `quality_gate.go` 中出现 7 次
- 颜色值 `#7DCFFF`、`#FF8700`、`#9ECE6A` 硬编码在 `init.go` 和 `init_surfaces.go`
- 哨兵数 `99999` 用于依赖环检测，无命名常量
- 八进制权限混用 `0644`（旧式）和 `0o644`（Go 1.25 式）
- 重试参数（`3` 次、`5*time.Second`）内联在函数调用中

**死代码残留**：
- deprecated `Scope` 字段保留在 `FrontmatterData` 和 `Task` 类型中
- `Debugf` 函数在 `internal/cmd/output.go` 和 `internal/cmd/base/output.go` 各有一份
- `checkExistingTaskState`、`getTaskPhase`、`compareVersionIDs` 等仅为 API 兼容保留的别名函数
- `.out` 构建产物（`cmd.out`、`cout.out`、`coverage.out`、`just.out`）提交在仓库中

**包组织无原则**：
- `internal/cmd/` 下 10+ 顶层命令文件散落在根目录，与已子包化的 6 个组不一致
- `pkg/` 有 19 个包，粒度不均（`project/` 仅 1 个文件，`task/` 有 33+ 个文件）
- `infocmd/` 定位模糊（shared utilities 还是特定领域）

### Urgency

v3.0.0 是唯一的大版本重构窗口。发布后 API 和包结构将趋于稳定，技术债修复成本指数增长。当前分支已重命名完成（`task` → `forge`），是建立规范的最后时机。

## Proposed Solution

两阶段推进：**Phase 1 输出规范**，**Phase 2 按规范重组代码**。

**Phase 1 — 规范建立**：分析现有代码库模式，扩展 `docs/conventions/` 下的规范文件，新增包组织、命名、常量管理、死代码管理等领域规范。这些规范将成为 `forge-cli/` 代码的唯一权威标准。

**Phase 2 — 代码重组与清理**：以新规范为指导，全面重新设计 `internal/cmd/` 和 `pkg/` 两层包结构，同时彻底清除所有已识别的死代码和魔法值。不保留兼容层。

### Innovation Highlights

此方案并非创新，而是工程实践的标准操作——在重构窗口期建立规范并执行。借鉴了 Go 标准库的包组织哲学（领域合并、依赖方向严格）和 `golangci-lint` 的 `goconst` 规则思想。

## Requirements Analysis

### Key Scenarios

- **开发者编写新命令**：查阅 `package-organization.md` 确定文件放置位置和包归属
- **Code Review 审查**：依据 `constants.md` 和 `naming.md` 判断 PR 是否符合规范
- **包结构重组**：依据 `package-organization.md` 中的依赖方向规则迁移文件
- **清理魔法值**：依据 `constants.md` 中的分类规则提取命名常量

### Non-Functional Requirements

- **向后兼容**：此为 v3.0.0 内部重构，不影响已发布 API（二进制尚未正式发布）
- **构建稳定性**：重组过程中每个提交都保持 `go build` 和 `go test` 通过
- **规范可发现性**：所有规范文档通过 `docs/conventions/` 统一入口访问

### Constraints & Dependencies

- Go 1.25 语言特性（`0o644` 八进制字面量等）
- Cobra 框架的命令注册模式
- 现有依赖方向规则：`cmd -> internal -> pkg`
- `pkg/types/` 作为 leaf package 不导入其他 forge-cli 包

## Alternatives & Industry Benchmarking

### Industry Solutions

Go 社区的标准实践是：
- `golang-standards/project-layout` 定义了 `cmd/`、`internal/`、`pkg/` 的职责
- Go 标准库自身使用领域合并策略（如 `net/http` 包含 HTTP 协议所有子领域）
- `goconst` linter 自动检测重复字符串常量

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零工作量 | 技术债持续增长，v3.0.0 后更难修复 | Rejected: v3.0.0 是最后窗口 |
| 仅输出规范文档 | golang-standards | 低风险，快速交付 | 无实际代码改进，规范可能与实践脱节 | Rejected: 用户要求实际清理 |
| 规范先行 + 代码重组 | 本方案 | 规范指导实践，审计可追溯 | 工作量较大 | **Selected: 两阶段确保方向正确** |

## Feasibility Assessment

### Technical Feasibility

完全可行。Go 的包重组主要是文件移动和 import 路径更新，工具链（`gorename`、IDE refactor）支持良好。

### Resource & Timeline

单人可完成。Phase 1（规范输出）约 1-2 天，Phase 2（代码重组 + 清理）约 3-5 天。

### Dependency Readiness

无外部依赖。所有涉及的包都是 `forge-cli` 内部包。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 全面重组是解决"规范缺失"的最佳路径 | XY Problem Detection | **Overridden**: 用户确认 v3.0.0 是唯一的重组窗口，规范与重组需同步推进 |
| 死代码应保留兼容层以平滑迁移 | 5 Whys | **Overturned**: v3.0.0 未正式发布，无外部消费者需要兼容 |
| `pkg/` 19 个包的粒度总体合理 | Assumption Flip | **Overturned**: 探索发现 `project/`（1 文件）、`infocmd/`（定位模糊）等不合理案例，需领域合并 |

## Scope

### In Scope

1. **扩展 `docs/conventions/enum-constants.md`**：增加非枚举常量管理规则（路径常量、超时值、颜色值）
2. **扩展 `docs/conventions/code-structure.md`**：增加包组织相关的结构规则
3. **新增 `docs/conventions/package-organization.md`**：`internal/cmd/` 和 `pkg/` 的包职责划分、依赖方向、文件组织原则
4. **新增 `docs/conventions/naming.md`**：文件名、函数名、常量名、包名命名规范
5. **新增 `docs/conventions/constants.md`**：魔法值全面管理策略（分类、提取规则、集中管理位置）
6. **新增 `docs/conventions/dead-code.md`**：死代码识别标准、deprecation 策略、清理流程
7. **重组 `internal/cmd/` 包结构**：顶层散落的命令文件子包化，统一命令注册模式
8. **重组 `pkg/` 层**：按领域合并小包，明确每个包的职责边界
9. **消除重复**：统一 `Debugf` 等重复工具函数到唯一位置
10. **删除所有死代码**：deprecated Scope 字段、别名函数、兼容层、构建产物（`.out` 文件）
11. **提取所有魔法值为命名常量**：路径、颜色、超时、哨兵数、八进制权限统一使用 `0o` 前缀

### Out of Scope

- Plugin 代码（`plugins/forge/` 下的 skills、commands、hooks）——非 Go 代码
- `docs/features/` 目录结构
- 测试逻辑重写（仅更新 import 路径）
- 新功能开发
- CLI UX 变更（命令名、参数名、输出格式）
- `internal/embedded/` 层
- `forge-cli/CLAUDE.md` 更新（属于 `/learn` 范畴）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 包重组导致 import 路径大量变更，引入编译错误 | M | M | 每步重组后立即 `go build` + `go test` 验证 |
| 规范过于理想化，与实际代码模式冲突 | M | L | 规范基于现有代码模式提炼，而非凭空设计 |
| pkg/ 领域合并导致包内职责模糊 | L | M | 每个包内通过文件名和注释区分子职责 |
| 重组过程中破坏 golangci-lint 配置 | L | L | 重组后运行 `make lint` 验证 |

## Success Criteria

- [ ] SC-1: `grep -rn '"tests/results/' forge-cli/internal/ forge-cli/pkg/` 返回零结果（所有路径常量已提取）
- [ ] SC-2: `grep -rn 'lipgloss.Color("#' forge-cli/internal/ forge-cli/pkg/` 返回零结果（所有颜色常量已提取）
- [ ] SC-3: `grep -rn '\b99999\b' forge-cli/` 返回零结果（哨兵常量已命名）
- [ ] SC-4: `grep -rn '0644\|0755' forge-cli/internal/ forge-cli/pkg/` 返回零结果（统一 `0o` 前缀）
- [ ] SC-5: `internal/cmd/` 下零个顶层命令文件（所有命令均已子包化）
- [ ] SC-6: `Debugf` 函数在整个 `forge-cli/` 中仅存在一个定义
- [ ] SC-7: deprecated `Scope` 字段、所有别名函数、所有 `.out` 构建产物已删除
- [ ] SC-8: `docs/conventions/` 包含 6 个与 forge-cli 相关的规范文件（扩展 2 个 + 新增 4 个）
- [ ] SC-9: `pkg/` 层包数量不超过 12 个（当前 19 个，领域合并后减少）
- [ ] SC-10: `go build ./...` 和 `go test ./...` 在重组后全部通过

consistency_check_result:
  status: pass
  pairs_checked: 45
  conflicts_found: 0

## Next Steps

- Proceed to `/write-prd` to formalize requirements
