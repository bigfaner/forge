---
created: "2026-05-30"
author: "faner"
status: Approved
intent: "refactor"
---

# Proposal: Forge CLI 代码库重组与规范建立

## Problem

Forge CLI 代码库（`forge-cli/`）缺少全面的编码规范，导致代码风格不一致、魔法值散布、死代码残留、包结构无明确组织原则。当前 `docs/conventions/` 下虽有 11 个规范文件，但覆盖面偏向 Forge plugin 生态（skill 结构、分发模型、dispatcher 质量等），对 forge-cli 自身的包组织、命名、常量管理等编码标准覆盖不足，无法指导新代码的编写和现有代码的演进。

### Evidence

**魔法值散布**：
- `"tests/results/raw-output.txt"` 在生产代码中出现 2 次（`quality_gate.go`），测试文件中出现 17 次
- 颜色值 `#7DCFFF`、`#FF8700`、`#9ECE6A` 硬编码在 `init.go` 和 `init_surfaces.go`
- 哨兵数 `99999` 用于依赖环检测，无命名常量
- 八进制权限统一使用旧式 `0644`，应迁移为 Go 1.25 式 `0o644`
- 重试参数（`3` 次、`5*time.Second`）内联在函数调用中，经审计确认存在语义相同的重复调用点（同模块内多处以相同重试策略调用外部服务）

**死代码残留**：
- deprecated `Scope` 字段保留在 `FrontmatterData` 类型中（`pkg/task/frontmatter.go`）
- `Debugf` 函数在 `internal/cmd/output.go` 和 `internal/cmd/base/output.go` 各有一份
- `checkExistingTaskState`、`getTaskPhase`、`compareVersionIDs` 等为 test-bridge 别名函数（用于测试注入，其中 `getTaskPhase` 在 `validate_index.go` 生产代码中亦有 5 处调用，非纯粹死代码）
- `.out` 构建产物（`cmd.out`、`cout.out`、`just.out`）提交在仓库中

**包组织无原则**：
- `internal/cmd/` 下 15 个顶层命令文件散落在根目录，与已子包化的 7 个组不一致
- `pkg/` 有 17 个包，粒度不均（`project/` 仅 3 个文件，`task/` 有 22+ 个非测试文件）
- `infocmd/` 定位模糊（shared utilities 还是特定领域）

### Urgency

v3.0.0 是成本最低的重构窗口（非唯一，但错过后成本上升 9-17 天）。发布后 API 和包结构将趋于稳定，后续任何包移动都需同时维护兼容层（估计每个包移动增加 0.5-1 天兼容层维护工作，17 个包即 9-17 天额外开销）。当前分支已重命名完成（`task` → `forge`），代码库尚未被外部消费者依赖，是建立规范的成本最优时机。

## Proposed Solution

四阶段推进，按 blast radius 从小到大排序：**Phase 1 输出规范**，**Phase 2a 清理死代码**，**Phase 2b 提取魔法值**，**Phase 2c 重组包结构**。

**Phase 1 — 规范建立**：分析现有代码库，提炼目标态定义和偏差分析，扩展 `docs/conventions/` 下的规范文件，新增包组织、命名、常量管理、死代码管理等领域规范。产出必须包含每个领域的目标态定义（规范性，而非描述性）以及当前代码与目标态的偏差分析表。包组织规范须明确依赖方向规则：`cmd → internal → pkg`（严格单向，禁止反向依赖）；`pkg/` 内采用**三层层级模型**——`pkg/types/`（leaf，零内部依赖）、基础设施层（如 `pkg/index/`、`pkg/git/`，仅依赖 types，不依赖其他领域包）、领域层（如 `pkg/task/`、`pkg/forgeconfig/`，可依赖基础设施层和 types，但领域包之间禁止互相依赖）。此分层模型承认并规范化当前 `pkg/` 层的横向依赖现状（如 `pkg/infocmd` 被 4 个领域包导入），而非简单禁止。Phase 1 的第一个产出是用 `go list -json ./pkg/...` 自动生成 `pkg/` 层的完整依赖图作为事实基线。偏差分析范围控制：优先完成包组织和常量管理两个核心领域的逐文件审计，其余领域（命名、死代码）以模块级摘要覆盖，避免全量逐文件扫描导致时间线膨胀。

**Phase 2a — 死代码删除**：以新规范为指导，删除纯粹的死代码——deprecated `Scope` 字段、重复的 `Debugf` 定义、`.out` 构建产物。此阶段不涉及包移动，blast radius 最小。

**Phase 2b — 魔法值提取与 test-bridge 清理**：提取所有魔法值为命名常量（路径、颜色、超时、哨兵数、八进制权限）。同时清理 test-bridge 别名函数，区分两类：纯粹重导出（可直接删除）和内部函数导出（需评估测试迁移策略后再删除）。

**Phase 2c — 包结构重组**：以新规范为指导，重新设计 `internal/cmd/` 和 `pkg/` 两层包结构，消除重复工具函数。不保留兼容层。此阶段 blast radius 最大，放在最后。包合并前必须检查合并是否引入循环依赖——若 A 包依赖 B 包且 B 包依赖 A 包的部分功能，则不可合并，改为提取共享部分到新包或重新划分边界。执行顺序：先处理 leaf 包（无内部依赖者，如 `pkg/version/`、`pkg/types/`），再处理中间包，最后处理被依赖最多的包（如 `pkg/task/`），确保每步移动后立即验证编译。

### Innovation Highlights

此方案的核心价值不在于技术创新，而在于工程时机的把握和执行策略的设计：(1) 利用 v3.0.0 未发布窗口消除兼容层开销，(2) 用机器可验证的依赖图替代人工审计作为规范的事实基线，(3) 三层层级模型将"禁止横向依赖"的教条替换为可 CI 验证的分层规则。这些不是新发明，而是将 Go 标准库的包组织哲学（领域合并、依赖方向严格）和 `goconst` linter 的检测思想组合应用到特定窗口期。

## Requirements Analysis

### Key Scenarios

- **开发者编写新命令**：查阅 `package-organization.md` 确定文件放置位置和包归属。具体步骤：(1) 在 `internal/cmd/<command-group>/` 下创建新文件（如 `internal/cmd/deploy/deploy.go`），(2) 在 `root.go` 中注册命令，(3) 业务逻辑调用 `pkg/` 层对应领域包
- **Code Review 审查**：依据 `constants.md` 和 `naming.md` 判断 PR 是否符合规范，使用 PR review checklist 确认包结构变更合规
- **包结构重组**：依据 `package-organization.md` 中的依赖方向规则迁移文件
- **清理魔法值**：依据 `constants.md` 中的分类规则提取命名常量

### Non-Functional Requirements

- **向后兼容**：v3.0.0 二进制尚未正式发布，且 forge-cli 未发布 Go module（`go.mod` 中无外部可引用的 module path），因此不存在 monorepo 外部消费者。跨模块依赖审计（Phase 2a 前置条件）将验证此假设——若审计发现 monorepo 内存在对 `forge-cli` 的跨模块 import，则 Phase 2c 改为条件性执行（见 Dependency Readiness fallback），向后兼容 NFR 仅在审计确认无跨模块依赖时成立
- **可回滚性**：每个 Phase 为独立可回退的提交（或提交组），`git revert` 可恢复到重组前状态；Phase 2c 中每个包的重组为独立提交，可单独回退
- **构建稳定性**：重组过程中每个提交都保持 `go build` 和 `go test` 通过
- **规范可发现性**：所有规范文档通过 `docs/conventions/` 统一入口访问
- **规范可执行性**：将 `goconst`、`gofmt`、`go vet` 集入 `make lint` 作为 CI gate，防止新的魔法值和格式违规引入（见 scope item 13）；包组织规范通过 PR review checklist 人工执行（包结构变更需要 review 确认符合 `package-organization.md`，见 scope item 14）

### Constraints & Dependencies

- Go 1.25 语言特性（`0o644` 八进制字面量等）
- Cobra 框架的命令注册模式
- 现有依赖方向规则：`cmd -> internal -> pkg`
- `pkg/types/` 作为 leaf package 不导入其他 forge-cli 包

## Alternatives & Industry Benchmarking

### Industry Solutions

Go 社区的标准实践是：
- `golang-standards/project-layout` 定义了 `cmd/`、`internal/`、`pkg/` 的职责（注意：该仓库 README 明确声明 "This is NOT an official Go standard"，而是社区参考；本方案仅参考其目录职责定义，不视为标准规范）
- Go 标准库自身使用领域合并策略（如 `net/http` 包含 HTTP 协议所有子领域）
- `goconst` linter 自动检测重复字符串常量

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零工作量，零回归风险 | 技术债持续增长，v3.0.0 后更难修复 | Rejected: v3.0.0 是成本最低的重组窗口（错过后成本上升约 9-17 天兼容层维护），且用户明确要求实际清理而非仅文档输出 |
| 仅输出规范文档 | golang-standards | 低风险，快速交付，规范可指导后续迭代 | 无实际代码改进，规范可能与实践脱节，历史证明仅靠文档无法阻止违规累积 | Rejected: 用户要求实际清理；且现有 `docs/conventions/` 已有 3 个规范但代码仍大量违规，证明规范不配合代码改进则效果有限 |
| Lint 驱动渐进重构（无规范前置） | 行业常见 | 自动化程度高，增量改进风险低 | 缺乏统一目标态，lint 规则碎片化无法指导包结构重组；`goconst` 可解决魔法值但无法解决包组织无原则问题 | Rejected: 包结构重组需要全局视角，lint 是局部工具；且 v3.0.0 窗口要求一次性确立结构，非渐进试错 |
| 增量每包迁移（gopls workspace refactor） | Go 工具链 | 利用 `gopls` 内置的 rename/move 重构，零手动 import 更新，每个包独立迁移可回退 | 缺乏全局目标态，每次迁移的方向可能不一致；`gopls` 的 move refactoring 对跨目录移动的可靠性取决于 LSP server 状态 | Rejected: 逐包迁移缺少整体规划，容易在方向上反复；且本方案已在 Phase 2c 采用 `gopls` 作为执行工具，差异在于是否有前置规范指导方向 |
| 规范先行 + 代码重组 | 行业参考：(1) Go 标准库 `net/http` 在 Go 1.0 前进行了大规模包结构重组（参见 Go 1 release notes "Changes to the Go standard library" 及 rsc.io/go-package-narrative 对标准库包演变的分析）；(2) `golangci-lint` 在引入新 linter 前先在 CONTRIBUTING.md 中定义目标行为和接受标准（参见 golangci-lint/golangci-lint CONTRIBUTING.md "Adding a new linter" 章节）；(3) `helm` v3 重构时重新定义 `pkg/` 层包职责（参见 helm/helm commit history 2019-2020 期间 chartutil 和 `pkg/` 目录重构） | 规范指导实践，审计可追溯，blast radius 可隔离 | 总计约 6-10 天工作量，需前后一贯执行 | **Selected: 四阶段确保方向正确、风险可控** |

## Feasibility Assessment

### Technical Feasibility

完全可行。Go 的包重组主要是文件移动和 import 路径更新，工具链（`gopls` 内置重构、IDE refactor）支持良好。

### Resource & Timeline

单人可完成。Phase 1（规范输出 + 目标态定义 + 偏差分析）约 2-3 天（偏差分析优先覆盖包组织和常量两个核心领域的逐文件审计，其余领域以模块级摘要覆盖），Phase 2a（死代码删除 + 跨模块依赖审计）约 0.5-1 天，Phase 2b（魔法值提取 + test-bridge 清理 + CI gate 集成）约 1.5-2.5 天（含 `goconst` linter 配置及现有违规修复），Phase 2c（包结构重组）约 2-3 天。总计约 6-10 天。

### Dependency Readiness

**前置条件**：Phase 2a 启动前必须完成跨模块依赖审计——检查 monorepo 内是否存在其他 Go 模块（如 plugin 相关代码）import `forge-cli` 的 `internal/` 或 `pkg/` 包。若存在跨模块依赖，需先解耦或纳入重组范围。审计方法：`go list -m` 和 `go mod graph` 检测 Go module 级依赖，辅以 `grep -rn 'forge-cli/internal\|forge-cli/pkg' --include='*.go'` 搜索 monorepo 根目录作为文本级补充检查。**Fallback**：若审计发现无法解耦的跨模块依赖，则 Phase 2c 改为保留必要的导出接口（内部标记 `// Deprecated`），而非执行完整包重组。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 全面重组是解决"规范缺失"的最佳路径 | XY Problem Detection | **Overridden**: 用户确认 v3.0.0 是唯一的重组窗口，规范与重组需同步推进 |
| 死代码应保留兼容层以平滑迁移 | 5 Whys | **Overturned**: v3.0.0 二进制未正式发布，且 forge-cli 未作为 Go module 发布（无外部可引用的 module path），monorepo 内跨模块依赖审计将在 Phase 2a 前验证此假设 |
| `pkg/` 17 个包的粒度总体合理 | Assumption Flip | **Overturned**: 探索发现 `project/`（3 文件小包）、`infocmd/`（定位模糊，被 4 个领域包横向依赖）等不合理案例，需领域合并或分层归属 |

## Scope

### In Scope

1. **扩展 `docs/conventions/enum-constants.md`**：增加非枚举常量管理规则（路径常量、超时值、颜色值）
2. **扩展 `docs/conventions/code-structure.md`**：增加包组织相关的结构规则
3. **新增 `docs/conventions/package-organization.md`**：`internal/cmd/` 和 `pkg/` 的包职责划分、依赖方向、文件组织原则
4. **新增 `docs/conventions/naming.md`**：文件名、函数名、常量名、包名命名规范
5. **新增 `docs/conventions/constants.md`**：魔法值全面管理策略（分类、提取规则、集中管理位置）
6. **新增 `docs/conventions/dead-code.md`**：死代码识别标准、deprecation 策略、清理流程
7. **重组 `internal/cmd/` 包结构**：顶层散落的命令文件子包化，统一命令注册模式
8. **重组 `pkg/` 层**：按领域合并小包，明确每个包的职责边界。当前→目标包映射表：

| 当前包 | 目标包 | 说明 |
|--------|--------|------|
| `pkg/project/`（3 文件） | 合并入 `pkg/task/` 或新 `pkg/workspace/` | 小包，领域归属待 Phase 1 确定 |
| `pkg/infocmd/` | 合并入 `pkg/shared/` 或 `internal/cmd/` | 定位模糊（共享通用工具），归入明确位置 |
| 其余 15 个包 | 分类处置 | **保留不变**（6 个）：`pkg/version/`、`pkg/types/`、`pkg/task/`、`pkg/forgeconfig/`、`pkg/git/`、`pkg/feature/`——领域清晰且职责单一。**评估合并至 `pkg/util/`**（~3 个）：`pkg/index/`、`pkg/serverprobe/` 等小工具包——`pkg/util/` 作为唯一允许被多个 `pkg/` 包共享的工具包，不引入横向依赖（`pkg/util/` 不依赖任何其他 forge-cli `pkg/` 包，仅依赖标准库和第三方库）。**保留并优化内部结构**（~4 个）：`pkg/just/`、`pkg/testrunner/` 等领域包——内部文件组织优化但不改名。**待 Phase 1 偏差分析裁决**（2 个）：`pkg/prompt/`（与 `pkg/feature/` 职责边界模糊——prompt 含 feature 查询逻辑）、`pkg/research/`（与 `pkg/proposal/` 共享 infocmd 工具链，可能合并）——Phase 1 给出明确归向 |
9. **消除重复**：统一 `Debugf` 等重复工具函数到唯一位置
10. **删除死代码**：deprecated `Scope` 字段、重复的 `Debugf` 定义、构建产物（`.out` 文件）
10a. **拆分超大文件**：将超过 500 行的 `.go` 文件按职责拆分（如 `quality_gate.go` 1067 行拆分为质量检查核心逻辑 + 报告生成逻辑）
11. **清理 test-bridge 别名函数**：区分两类——纯粹重导出别名（可直接删除，如 `checkExistingTaskState` 等仅转发调用的包装）和内部函数导出（需评估测试代码迁移后再删除）
12. **提取所有魔法值为命名常量**：路径、颜色、超时、哨兵数、八进制权限统一使用 `0o` 前缀
13. **CI gate 集成**：将 `goconst`、`gofmt`、`go vet` 集入 `make lint`，确保 `golangci-lint` 配置中启用 `goconst` linter；修复 `goconst` 报告的所有现有违规（与 Phase 2b 魔法值提取同步完成）；预计 0.5 天，纳入 Phase 2b 时间线
14. **PR review checklist**：在 `docs/conventions/package-organization.md` 中附加 PR review checklist 条目（包结构变更需 review 确认符合依赖方向规则和包职责定义），预计 0.5 天，纳入 Phase 1 时间线
15. **跨模块依赖 CI check**：Phase 2a 完成时在 `Makefile` 中添加 `make check-cross-module-deps` target（`grep -rn '"forge-cli/internal'` 搜索 monorepo 排除 `forge-cli/` 自身，返回非零则构建失败），为 Phase 2c 的"不保留兼容层"决策提供持续保护

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
| 包重组导致 import 路径大量变更，引入编译错误 | M | M | Phase 2a/2b 不涉及包移动，Phase 2c 每步重组后立即 `go build` + `go test` 验证 |
| 规范过于理想化，与实际代码模式冲突 | M | H | Phase 1 产出须由项目维护者 review 后方可进入 Phase 2；若 review 发现规范与实际代码严重冲突，回退至纯描述性文档并基于冲突点修订规范。**下游影响**：若规范回退为描述性文档，Phase 2c 的目标包映射表将失去权威依据，此时 Phase 2c 缩减为仅执行 Phase 1 中共识度最高的合并项（3 个明确合并目标），其余包保持现状。**回退产出定义**：即使规范回退，Phase 1 仍产出两个有价值的交付物——(1) `pkg/` 层的机器可验证依赖图（事实基线），(2) 每个领域的偏差分析表（即使规范本身是描述性的，偏差分析仍具有审计价值）。因此 Phase 1 不会因规范回退而完全浪费 |
| pkg/ 领域合并导致包内职责模糊 | L | M | 每个合并后的包须在 `doc.go` 中列出包含的子领域及其职责边界。**合并回退标准**：若合并后的包内出现循环 import、单一包超过 15 个文件、或 `go vet` 报告命名冲突，则该合并回退（`git revert`），改为保留原包结构并仅在 `doc.go` 中补充职责说明 |
| 重组过程中破坏 golangci-lint 配置 | L | L | 重组后运行 `make lint` 验证 |
| Phase 2c 引入运行时回归（编译通过但行为变化） | L | H | 每个 Phase 为独立可回退的提交组，`git revert` 可恢复到重组前状态；Phase 2c 中每个包的重组为独立提交，可单独回退 |
| Phase 1 偏差分析范围超出预期，导致时间线膨胀 | M | M | 偏差分析按优先级排序，先完成包组织和常量两个核心领域，其余领域可降级为简表；若 Phase 1 超过 3 天则缩减范围 |
| Phase 1 依赖图分析揭示三层模型无法覆盖的耦合模式（如 `pkg/infocmd` 被领域包依赖但自身也依赖领域包，形成双向耦合） | M | H | Phase 1 的依赖图产出将暴露此类问题；若发现双向耦合，将涉及的包标记为"待解耦"并从 Phase 2c 的合并范围中排除，改为在 Phase 2c 中仅做 `doc.go` 职责标注而不移动文件 |

## Success Criteria

- [ ] SC-1: `grep -rn '"tests/results/' forge-cli/internal/ forge-cli/pkg/` 返回零结果（所有路径常量已提取）
- [ ] SC-2: `grep -rn 'lipgloss.Color("#' forge-cli/internal/ forge-cli/pkg/` 返回零结果（所有颜色常量已提取）
- [ ] SC-3: `grep -rn '\b99999\b' forge-cli/` 返回零结果（哨兵常量已命名）
- [ ] SC-4: `grep -rn '0644\|0755' forge-cli/internal/ forge-cli/pkg/` 返回零结果（统一 `0o` 前缀）
- [ ] SC-5: `internal/cmd/` 下零个顶层命令实现文件。**豁免**（基础设施文件，非命令实现）：`root.go`（命令注册入口）、`output.go`（共享输出工具）、`surfaces.go`（共享 surface 定义）、`surfaces_detect.go`（surface 检测工具，被 `surfaces.go` 调用）。**分类标准**：文件属于豁免类当且仅当它不包含任何 Cobra 命令的 `Run`/`RunE` 实现
- [ ] SC-6: `Debugf` 函数在整个 `forge-cli/` 中仅存在一个定义
- [ ] SC-7: deprecated `Scope` 字段、重复 `Debugf` 定义、所有 `.out` 构建产物已删除；test-bridge 纯重导出别名已删除、内部导出别名已迁移
- [ ] SC-8: `docs/conventions/` 包含 6 个与 forge-cli 相关的规范文件（扩展 2 个 + 新增 4 个），每个包含目标态定义和偏差分析；`make lint` 包含 `goconst`、`gofmt`、`go vet` 且 CI 通过
- [ ] SC-9: `pkg/` 层包数量不超过 14 个（当前 17 个，3 个明确合并目标 + ~3 个小工具包合并至 `pkg/util/` 可减少至 ~13 个；待裁决 2 个包的最终归向由 Phase 1 确定，上限 14 个留有缓冲）
- [ ] SC-10: `internal/cmd/` 层无超过 500 行的单个 `.go` 文件（当前 5 个文件超标：`quality_gate.go` 1067 行、`init.go` 591 行、`init_surfaces.go` 550 行、`task/validate_index.go` 521 行、`task/tree.go` 504 行），在 Phase 2c 命令子包化过程中同步拆分。`pkg/` 层的超大文件（`forgeconfig/config.go` 1272 行、`task/pipeline.go` 1097 行等 5 个文件）作为后续迭代目标，不阻塞当前提案。拆分标准：每个文件围绕单一职责/功能组织，通过代码审查确认拆分后每个文件的内聚性（不因行数阈值而机械切割）
- [ ] SC-11: `go build ./...` 和 `go test ./...` 在重组后全部通过
- [ ] SC-12: 跨模块依赖审计已完成，审计结果记录在 Phase 2a 前置文档中
- [ ] SC-12f（fallback 条件）: 若跨模块依赖审计发现无法解耦的依赖，则 `pkg/` 中保留的导出接口均标记 `// Deprecated` 注释，且保留的包数量不超过 16 个（从 17 个减少至少 3 个明确合并目标：`project/`、`infocmd/`、~1 个小工具包）；SC-5 豁免列表不变；SC-9 上限调整为 16 个
- [ ] SC-13: `pkg/` 层超大文件（`forgeconfig/config.go` 1272 行、`forgeconfig/detect_surface.go` 962 行、`task/pipeline.go` 1097 行、`task/build.go` 638 行、`task/autogen.go` 518 行）已在 Phase 1 偏差分析中记录为后续迭代目标，并产出拆分可行性评估（评估每个文件的拆分是否涉及导出接口变更）
- [ ] SC-14: `docs/conventions/package-organization.md` 包含 PR review checklist（包结构变更需 review 确认符合依赖方向规则和包职责定义）

consistency_check_result:
  status: issues_found
  pairs_checked: 48
  conflicts_found: 2
  issues:
    - SC-9 目标可达性（已调整至 <= 14 并分类处置其余包，Phase 1 裁决 ~2 个待定包）
    - SC-12f fallback 条件下的 SC-5/SC-9 调整（已声明按实际可达目标调整，具体值在审计后确定）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
