# Freeform Review: Forge CLI 代码库重组与规范建立

**审阅者视角**: Go 大规模重构工程师，8+ 年代码库健康治理经验
**审阅文档**: `docs/proposals/forge-cli-codebase-standards/proposal.md`
**审阅日期**: 2026-05-30

---

## 第一部分：背景评估

### 它声称解决什么问题

本提案声称 Forge CLI 代码库（`forge-cli/`）缺乏全面的编码规范，导致三类症状：魔法值散布（路径字符串如 `"tests/results/raw-output.txt"` 出现 19 次、颜色值 `#7DCFFF` 等硬编码、哨兵数 `99999`、旧式八进制权限 `0644`、内联重试参数）；死代码残留（deprecated `Scope` 字段、重复的 `Debugf` 定义、test-bridge 别名函数、`.out` 构建产物）；以及包结构无明确组织原则（`internal/cmd/` 下 15 个文件散落根目录、`pkg/` 层 17 个包粒度不均）。

提案选择 v3.0.0 发布前作为重构窗口，声称：

> "发布后 API 和包结构将趋于稳定，后续任何包移动都需同时维护兼容层（估计每个包移动增加 0.5-1 天兼容层维护工作，17 个包即 9-17 天额外开销）"

### 核心技术路线

四阶段推进：Phase 1 输出规范文档（含目标态定义和偏差分析），Phase 2a 删除死代码并审计跨模块依赖，Phase 2b 提取魔法值并配置 CI linter，Phase 2c 重组包结构。每阶段独立可回退，按 blast radius 从小到大排列。这是一个务实且保守的工程方案，底色扎实。

### 它依赖的假设

提案建立在五个核心假设上：

1. forge-cli 没有外部 Go module 消费者（`go.mod` 中无外部可引用的 module path）
2. monorepo 内不存在对 `forge-cli/internal/` 或 `forge-cli/pkg/` 的跨模块 import
3. `pkg/types/` 是 leaf package，不导入任何其他 forge-cli 包
4. v3.0.0 是成本最低的重构窗口
5. Go 工具链足以安全完成包移动

我逐一验证了这些假设。假设 2 在当前代码库中成立（`tests/go.mod` 是独立 module 不引用 forge-cli，`plugins/` 下无 Go 代码）。假设 3 成立（`pkg/types/` 不 import 任何 forge-cli 包）。但假设 1 的安全性依赖于 Phase 2a 的前置审计是否被执行且被持久化——这是一个悬而未决的问题。

---

## 第二部分：关键风险识别

### 风险 1：`pkg/` 横向依赖与"禁止横向依赖"规则的根本矛盾

提案在 Phase 1 中定义了严格的依赖方向规则：

> "pkg/ 内禁止包间横向依赖；pkg/types/ 作为 leaf package 不导入任何其他 forge-cli 包"

这一规则在原则上是正确的，但与我审计到的实际代码存在严重冲突。`pkg/` 层当前的横向依赖关系如下：

- `pkg/infocmd` 被 `pkg/research`、`pkg/proposal`、`pkg/task`、`pkg/lesson` 共 4 个领域包导入
- `pkg/index` 被 `pkg/task`（3 处）、`pkg/feature` 共 4 个包导入
- `pkg/forgeconfig` 被 `pkg/task`（3 处）、`pkg/prompt` 共 4 个包导入
- `pkg/feature` 被 `pkg/proposal`、`pkg/prompt` 共 3 个包导入
- `pkg/just` 被 `pkg/testrunner` 导入
- `pkg/task` 被 `pkg/prompt` 导入

这不是 3 个小工具包合并到 `pkg/util/` 就能解决的表面问题。`pkg/infocmd` 被四个领域包依赖，如果移入 `internal/cmd/`（提案给出的选项之一），这四个 `pkg/` 包就会违反 `pkg/` 不依赖 `internal/` 的规则。如果创建 `pkg/shared/`，那本质上就是承认了横向依赖的存在——只是换了个名字。

> "评估合并至 pkg/util/（~3 个）：pkg/index/、pkg/serverprobe/ 等小工具包——pkg/util/ 作为唯一允许被多个 pkg/ 包共享的工具包"

但 `pkg/index` 不是一个通用工具包——它是文件锁和原子操作的基础设施，被 `pkg/task`（核心领域包）和 `pkg/feature` 同时依赖。将它并入 `pkg/util/` 只是掩盖了依赖关系，并没有真正解决领域包之间的耦合问题。

风险：如果 Phase 1 产出"禁止横向依赖"的规范，Phase 2c 执行时就会发现这条规则在当前架构下不可达——除非进行远比"合并小包"更深层的架构重组。这将导致规范回退（提案已预见此情况），浪费 Phase 1 的时间投入。

### 风险 2：`getTaskPhase` 等函数被错误归类为死代码

提案在 Evidence 部分声明：

> "checkExistingTaskState、getTaskPhase、compareVersionIDs 等仅为 API 兼容保留的别名函数"

我在代码中验证了这些函数的实际使用情况。`getTaskPhase` 在 `validate_index.go`（生产代码）中被调用了 5 次（第 299、333、369、381、410 行）。这不是一个"仅为 API 兼容保留的别名"——它是生产逻辑中活跃使用的内部函数。

`checkExistingTaskState` 在 `claim_integration_test.go` 中被调用 7 次。`compareVersionIDs` 在 `claim_test.go` 中被测试。这些是测试桥接函数，不是死代码。

风险：如果执行者在 Phase 2a 中按提案描述将 `getTaskPhase` 作为"死代码"删除，`validate_index.go` 的生产逻辑会立即编译失败。更危险的是，如果执行者基于此分类在 Phase 2a 中大范围删除"别名函数"，可能会误删仍在生产中使用的函数。提案自身的死代码分类存在事实性错误，会误导实施者。

### 风险 3：SC-10 的 500 行阈值覆盖范围不完整

提案声明：

> "SC-10: 无超过 500 行的单个 .go 文件（当前 quality_gate.go 1067 行等巨型文件需按职责拆分）"

但实际代码库中超过 500 行的非测试文件不止 `quality_gate.go`：

| 文件 | 行数 |
|------|------|
| `pkg/forgeconfig/config.go` | 1272 |
| `pkg/task/pipeline.go` | 1097 |
| `internal/cmd/quality_gate.go` | 1067 |
| `pkg/forgeconfig/detect_surface.go` | 962 |
| `pkg/task/build.go` | 638 |
| `internal/cmd/init.go` | 591 |
| `internal/cmd/init_surfaces.go` | 550 |
| `internal/cmd/task/validate_index.go` | 521 |
| `pkg/task/autogen.go` | 518 |
| `internal/cmd/task/tree.go` | 504 |

共 10 个文件超过 500 行。提案只点名了 `quality_gate.go`，其他 9 个没有被提及。其中 `config.go`（1272 行）和 `pipeline.go`（1097 行）比 `quality_gate.go` 还大。拆分 `config.go` 可能涉及重组 `forgeconfig` 包的导出接口——这是设计变更，不是简单的文件切割。

风险：如果 SC-10 被视为硬性成功标准，实施者需要拆分 10 个文件而非提案暗示的 1-2 个。这将显著膨胀 Phase 2c 的工作量（2-3 天的估计可能需要翻倍），且部分拆分（如 `config.go`）可能引入接口变更，超出纯重构的边界。

### 风险 4：跨模块依赖审计缺乏持久化保障

提案声明：

> "Phase 2a 启动前必须完成跨模块依赖审计——检查 monorepo 内是否存在其他 Go 模块（如 plugin 相关代码）import forge-cli 的 internal/ 或 pkg/ 包"

审计方法仅列出了 `go list -m`、`go mod graph` 和 `grep -rn`。这些是一次性命令。我在当前代码库中验证了：`tests/go.mod` 是独立 module 不依赖 forge-cli，`plugins/` 下没有 Go 代码。所以审计在当下会通过。

风险：但提案声称"不保留兼容层"是基于"不存在跨模块依赖"这个假设。如果未来有人在 monorepo 中新增了一个 Go 模块并 import 了 `forge-cli/pkg/`，Phase 2c 的包重组就会破坏那个模块。一次性 grep 不能保护未来。提案没有要求将这个审计固化为 CI check，也没有要求在 `Makefile` 中添加自动化验证。

问题：提案的 fallback 条件是：

> "若审计发现无法解耦的跨模块依赖，则 Phase 2c 改为保留必要的导出接口（内部标记 // Deprecated），而非执行完整包重组"

这个 fallback 是合理的，但触发条件是一次性的手动审计。如果审计在 Phase 2a 时通过，Phase 2c 执行期间有人新增了跨模块依赖，fallback 就不会被触发。

### 风险 5：规范回退路径的下游影响未被充分分析

提案预见了一个风险：

> "若 review 发现规范与实际代码严重冲突，回退至纯描述性文档并基于冲突点修订规范"

以及下游影响：

> "若规范回退为描述性文档，Phase 2c 的目标包映射表将失去权威依据，此时 Phase 2c 缩减为仅执行 Phase 1 中共识度最高的合并项（3 个明确合并目标），其余包保持现状"

问题：如果 Phase 1 花费 2-3 天产出规范性文档，review 后发现规范与代码冲突需要回退，那 Phase 1 的产出变成了什么？是纯描述性文档（没有指导意义）？还是"3 个明确合并目标"的执行计划？如果是后者，那 Phase 1 的价值就从"全面规范建立"降级为"3 个包的合并方案"，投入产出比需要重新评估。

---

## 第三部分：改进建议

建议：在 Phase 1 开始前（或作为 Phase 1 的第一个产出），用自动化工具生成当前 `pkg/` 层的完整依赖图。具体做法：`go list -json ./pkg/... | jq '.ImportPath, .Imports'` 输出每个包的导入关系，整理为可视化的有向图。这个依赖图应该被提交到仓库中作为事实基线。它有两个直接价值：(1) 让 Phase 1 的偏差分析建立在机器可验证的事实上，而不是人工审计上，避免写出"禁止横向依赖"这种与现状矛盾的规则；(2) 为 Phase 2c 的包移动提供依赖安全的验证手段——每次移动后重新生成依赖图，与基线对比确认没有引入新的不当依赖。这不会增加 Phase 1 的总时间（生成依赖图只需几分钟），但能避免人工审计的系统性遗漏。参考风险 1。

建议：将 `pkg/` 内的依赖规则从"禁止横向依赖"改为"分层允许"。定义三个子层：`pkg/types/`（零内部依赖的 leaf）、基础设施层（如 `pkg/index/`、`pkg/git/`、未来的 `pkg/util/`——仅依赖 types，不依赖其他领域包）、领域层（如 `pkg/task/`、`pkg/forgeconfig/`——可依赖基础设施层和 types，但领域包之间禁止互相依赖）。这个分层规则可以用一个简单的 shell 脚本在 CI 中验证：对每个 `pkg/{domain}/` 包，检查其 imports 是否包含其他 `pkg/{other-domain}/` 包。当前代码中 `pkg/infocmd` 被 4 个领域包导入的问题，可以通过将 infocmd 的功能拆分为基础设施层（通用查询工具）和领域层（特定领域查询）来解决。参考风险 1。

建议：Phase 2a 的死代码分类必须在实施前用自动化工具交叉验证。"仅为 API 兼容保留的别名函数"这个标签不能仅靠人工判断。建议使用 `grep -rn '函数名' --include='*.go' | grep -v _test.go | grep -v 'func 函数名'` 来验证每个候选死代码函数在生产代码中的实际调用次数。对于 `getTaskPhase`，这个检查会立即揭示它在 `validate_index.go` 中有 5 处生产调用，不应被归类为死代码。同时建议在提案中修正这个分类错误，将 `getTaskPhase` 从 Phase 2a 的死代码清单中移除。参考风险 2。

建议：SC-10 的 500 行阈值应该区分两类处理方式。对于 `quality_gate.go`（1067 行）、`init.go`（591 行）、`init_surfaces.go`（550 行）等命令层文件，可以在 Phase 2c 的命令子包化过程中自然拆分——将不同职责的命令逻辑分入不同子包。但对于 `pkg/forgeconfig/config.go`（1272 行）、`pkg/task/pipeline.go`（1097 行）等核心领域文件，拆分需要设计性决策（哪些导出符号放哪个文件、是否引入子包），应该作为独立的 backlog 项而非阻塞当前提案的成功标准。建议 SC-10 限定范围为 `internal/cmd/` 层（与命令子包化同步拆分），`pkg/` 层的超大文件作为后续迭代目标。参考风险 3。

建议：跨模块依赖审计应该被固化为 CI 中的一道 check，而不是一次性命令。在 `Makefile` 或 CI pipeline 中添加一个 target：`grep -rn '"forge-cli/internal'` 搜索 monorepo（排除 `forge-cli/` 自身），如果返回非零结果则构建失败。这个 check 应该在 Phase 2a 完成时就添加，为 Phase 2c 的"不保留兼容层"决策提供持续保护。同时，审计结果应该以文档形式记录在 Phase 2a 的产出中（而不仅仅是执行了 grep 命令），作为 Phase 2c 启动的前置条件检查项。参考风险 4。
