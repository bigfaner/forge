# Freeform Review: Forge CLI 代码库重组与规范建立

**Reviewer Profile**: Pragmatic Go refactorer, 8+ years large-scale codebase health
**Document**: `docs/proposals/forge-cli-codebase-standards/proposal.md`
**Date**: 2026-05-30

---

## Section 1: Background Assessment

This proposal aims to address accumulated technical debt in the `forge-cli/` codebase by (a) establishing comprehensive coding conventions and (b) restructuring package layout and eliminating dead code. The core approach is a two-phase plan: Phase 1 produces convention documents, Phase 2 executes the reorganization under those conventions. The proposal correctly identifies the v3.0.0 pre-release window as a rare opportunity where breaking internal APIs is nearly free.

The evidence marshalled is specific and verifiable: hardcoded path strings, duplicate `Debugf` definitions, sentinel `99999`, mixed octal literal styles, orphan `.out` build artifacts, and an uneven `pkg/` layer with packages ranging from 1 to 22 non-test files. I confirmed most of these claims against the actual codebase and they hold up.

The proposal's fundamental assumption — that conventions must precede restructuring — is sound in principle but creates a sequencing risk I will elaborate on below.

---

## Section 2: Key Risk Identification

风险：Phase 2 将包重组、死代码删除、魔法值提取、重复消除四类变更混合在同一阶段内执行。引用原文：

> "Phase 2 — 代码重组与清理：以新规范为指导，全面重新设计 internal/cmd/ 和 pkg/ 两层包结构，同时彻底清除所有已识别的死代码和魔法值。不保留兼容层。"

"全面重新设计"加上"同时彻底清除"描述了一个大爆炸式的变更范围。包重组（文件移动 + import 路径更新）和死代码删除（符号移除）的 blast radius 完全不同。包重组的风险是编译错误，死代码删除的风险是运行时行为回归。将它们捆绑在一起意味着：如果一个重组步骤引入了编译问题，你无法判断是否同时有死代码误删导致的逻辑丢失。应该先删除死代码（减少移动目标），再重组包结构（移动更少的文件）。

问题：提案对 test-bridge 模式的处理存在严重盲区。引用原文的 In Scope 第 9 项：

> "消除重复：统一 Debugf 等重复工具函数到唯一位置"

以及 Scope 第 10 项：

> "删除所有死代码：deprecated Scope 字段、别名函数、兼容层、构建产物（.out 文件）"

我检查了 `internal/cmd/task/testbridge.go`，发现该文件是一个 125 行的 test-bridge 层，导出了大量符号（`ExportRunSubmit`、`ExportClaimNextTask`、`ExportExecuteClaim` 等），其中一半指向已迁移到 `pkg/task/` 的函数（纯粹的重导出），另一半引用 `internal/cmd/task/` 内部的未导出函数。同时 `claim.go` 中有三个 `var` 别名（`checkExistingTaskState`、`getTaskPhase`、`compareVersionIDs`）也是测试桥接的一部分。

提案将别名函数归入"死代码"要删除，但这些别名是生产代码中的 package-level `var` 声明，被同包的测试文件直接引用。删除它们会破坏测试编译，而测试编译失败会阻断 `go test` 的安全网。提案没有分析 test-bridge 的清理策略：是改为直接 import `pkg/task` 的导出函数？还是将测试迁移到 `internal/cmd/task/` 包内？还是用 `testify` 替代？这是一个独立的技术决策，不能与"删除死代码"混为一谈。

风险：`quality_gate.go` 是一个 1067 行的巨型文件，但提案没有将其列入处理范围。引用原文的 Scope：

> "重组 internal/cmd/ 包结构：顶层散落的命令文件子包化，统一命令注册模式"

这意味着 `quality_gate.go` 作为 `internal/cmd/` 下的顶层文件应该被"子包化"，但提案没有说明这个 1067 行文件是否应该拆分。如果只是移动到子目录而不拆分，那么重组只是把问题藏到了子目录里。类似地，`pkg/forgeconfig/config.go`（1272 行）和 `pkg/task/pipeline.go`（1097 行）也超过了 500 行的健康阈值，但提案的 Success Criteria 只检查包数量（SC-9），不检查单文件行数。

问题：SC-5 的成功标准自相矛盾。引用原文：

> "SC-5: internal/cmd/ 下零个顶层命令文件（所有命令均已子包化）"

但 `internal/cmd/root.go`（Cobra 根命令）、`output.go`（输出工具函数）、`surfaces.go`、`surfaces_detect.go`（surface 检测）不是"命令文件"——它们是共享基础设施。如果这些也必须子包化，那么每个子包中的命令如何 import 共享的 `root.go` 中的 `RootCmd`？如果它们不算"命令文件"可以留在根目录，那么 SC-5 的表述需要明确哪些文件豁免。

风险：引用原文关于 `pkg/` 层重组的目标：

> "SC-9: pkg/ 层包数量不超过 12 个（当前 19 个，领域合并后减少）"

我验证了实际的 `pkg/` 目录结构：18 个子目录（而非提案声称的 19 个）。更重要的是，其中 6 个包只有 1 个非测试 Go 文件（`facttable`、`infocmd`、`lesson`、`proposal`、`research`、`serverprobe`、`version`），这些都是合并候选。但合并它们需要分析依赖方向。例如，`infocmd/` 的"定位模糊"被提案指出但未解决——如果 `infocmd/` 合并到某个更大的包中，它的导出符号是否会被其他包 import？提案缺少一个目标包结构的具体映射表（当前包 -> 目标包）。

问题：提案声称 `"tests/results/raw-output.txt"` 在 `quality_gate.go` 中出现 7 次。我实际验证发现：`quality_gate.go` 中出现 2 次（第 255 行和第 284 行），其余 7 次出现在 `quality_gate_test.go` 中。原文表述为：

> "tests/results/raw-output.txt 在 quality_gate.go 中出现 7 次"

这是事实性错误。测试文件中的魔法值也需要提取，但将两者混为一谈会误导实施者对工作量的估计。

风险：引用原文的约束：

> "向后兼容：此为 v3.0.0 内部重构，不影响已发布 API（二进制尚未正式发布）"

但随后又说：

> "不保留兼容层"

这两者在 v3.0.0 的语境下是一致的，但如果 `forge-cli` 被其他 Go 模块通过 `go.mod` replace 指令引用（即使在 monorepo 内部），删除导出符号会破坏那些引用。提案没有检查是否有其他 Go 模块 import `forge-cli/pkg/` 或 `forge-cli/internal/` 下的包。在 monorepo 中，`go.mod` replace 指令可以创建隐形的跨模块依赖。

问题：提案在 Feasibility Assessment 中引用了 `gorename` 工具：

> "Go 的包重组主要是文件移动和 import 路径更新，工具链（gorename、IDE refactor）支持良好。"

`gorename` 已经不再维护（最后更新 2018 年），且不支持 Go modules。现代替代方案是 `gopls` 的重命名功能或 IDE 内置重构。引用过时工具说明技术可行性评估不够严谨。

风险：Phase 1 和 Phase 2 之间的规范-实践脱节风险。引用原文的 Phase 1 描述：

> "Phase 1 — 规范建立：分析现有代码库模式，扩展 docs/conventions/ 下的规范文件，新增包组织、命名、常量管理、死代码管理等领域规范。这些规范将成为 forge-cli/ 代码的唯一权威标准。"

规范"基于现有代码库模式提炼"意味着 Phase 1 的产出是对当前状态的描述性文档，而非对目标状态的规范性文档。如果 Phase 1 描述的是"现状"而 Phase 2 要执行的是"理想态"，那么 Phase 2 缺乏权威指导。如果是"理想态"，那么 Phase 1 的描述就不是"提炼"而是"设计"。这个矛盾没有解决。

---

## Section 3: Improvement Suggestions

建议：将 Phase 2 拆分为至少三个子阶段，按 blast radius 从小到大排序：

1. **Phase 2a — 死代码删除**：删除 `.out` 构建产物、deprecated `Scope` 字段、`internal/cmd/task/testbridge.go` 中指向 `pkg/task/` 的纯粹重导出。这一步只删除，不移动任何文件。每删一个类别就跑 `go build ./...` + `go test ./...`。

2. **Phase 2b — 魔法值提取**：将路径字符串、颜色值、哨兵数、八进制权限提取为命名常量。这一步只修改字面量，不改变任何函数签名或包结构。

3. **Phase 2c — 包结构重组**：移动文件、更新 import 路径、合并小包。这一步只移动，不删除任何逻辑。

这种排序确保每一步的失败模式是单一维度的。

建议：在 Scope 中增加一个明确的目标包结构映射表。当前的描述是定性的（"按领域合并小包"），但缺少定量的映射。例如：

| 当前包 | 目标包 | 理由 |
|--------|--------|------|
| `pkg/infocmd/` | 合并到 `pkg/project/` | 共享项目级查询职责 |
| `pkg/serverprobe/` | 合并到 `pkg/git/` 或独立保留 | 需要分析依赖 |

没有这个映射表，Phase 2 的实施者需要在执行时做出架构决策，这违反了"规范先行"的原则。

建议：引用原文的 test-bridge 处理：

> "删除所有死代码：deprecated Scope 字段、别名函数、兼容层、构建产物（.out 文件）"

`internal/cmd/task/testbridge.go` 中的别名函数不应归类为"死代码"。它们是活跃的测试基础设施。提案应该将 test-bridge 清理作为独立的 Scope 项，明确策略：

- 对于纯粹重导出（如 `ExportParseSegment = task.ParseSegment`）：直接删除，让测试 import `pkg/task`。
- 对于内部函数导出（如 `ExportRunSubmit = runSubmit`）：评估是否可以将这些函数移到 `pkg/task/`（消除 cmd 层依赖后），或接受它们作为必要的 test-bridge 模式。

建议：增加一个文件行数相关的成功标准。引用原文 SC-9：

> "SC-9: pkg/ 层包数量不超过 12 个（当前 19 个，领域合并后减少）"

包数量只是结构健康的维度之一。建议增加：

> `SC-11: forge-cli/ 下零个非测试 Go 文件超过 600 行`

这会迫使 `quality_gate.go`（1067 行）、`config.go`（1272 行）、`pipeline.go`（1097 行）等文件在重组中被拆分，而不仅仅是移动。

建议：引用原文的 Success Criteria 中缺少对 `golangci-lint` 的验证。当前只有：

> "SC-10: go build ./... 和 go test ./... 在重组后全部通过"

但 `CLAUDE.md` 中明确列出了 `golangci-lint run ./...` 作为常用命令，且提案在 Key Risks 中提到：

> "重组过程中破坏 golangci-lint 配置"

建议增加 `SC-12: golangci-lint run ./... 在重组后零新增 warning`，将 lint 作为硬性门控而非事后检查。

建议：引用原文的 evidence 部分关于 `"tests/results/raw-output.txt"` 的计数：

> "tests/results/raw-output.txt 在 quality_gate.go 中出现 7 次"

修正为准确描述：`quality_gate.go` 中 2 次，`quality_gate_test.go` 中 5+ 次。测试文件中的魔法值同样需要提取，但应作为独立的计数项目，以避免实施者误判生产代码中的实际重复密度。

建议：Phase 1 的规范产出应该包含一个明确的"目标态"描述，而不仅仅是"基于现有代码模式提炼"。引用原文：

> "分析现有代码库模式，扩展 docs/conventions/ 下的规范文件"

应该补充："规范文件必须包含两部分：(1) 当前状态的偏差分析（哪些代码违反规范），(2) 合规的迁移路径（不合规代码应如何修改）。"没有偏差分析的规范只是愿望清单，无法指导 Phase 2 的执行。

建议：提案应该检查并声明 monorepo 内是否存在跨 Go 模块的 import 依赖。在执行"不保留兼容层"之前，运行 `grep -r 'forge-cli/' --include='go.mod'` 确认没有其他模块通过 replace 指令引用 `forge-cli` 的导出符号。这应该作为 Phase 2 的前置条件。
