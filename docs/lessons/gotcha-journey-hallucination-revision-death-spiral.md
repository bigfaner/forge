---
created: "2026-05-29"
tags: [testing]
---

# Journey 生成幻觉与修订死循环

## Problem

`unify-enum-constants` feature 的 eval-journey 评分持续不达标（466→630→585，目标 850），3 轮迭代后反而下降。核心问题维度：
- **Fact Alignment (45/150)**：捏造常量名（StatusCancelled 实为 StatusSuspended）、捏造 CLI 命令（`forge task validate` 实为 `forge validate-index`）、未验证的文件计数
- **Surface Fitness (50/150)**：full-verification 缺少 not-found/already-exists 派生结果，journey 描述开发者工作流而非 CLI 子进程调用
- **Internal Consistency (90/150)**：Invariant 声称"绝不使用原始字符串"但 statemachine.go 使用 `types.Status("*")` 通配符模式

## Root Cause

**Level 1 — 生成器不验证事实**：gen-journeys skill 基于 feature 文档（tech design, user stories）生成 journey，但不交叉参照实际代码。常量名、CLI 命令、文件名、字符串计数全部是"合理推测"而非代码验证的结果。

**Level 2 — 修订循环无验证闭环**：eval-journey 的 adversary 能发现事实错误，但 reviser 在修正时同样无法访问代码验证。结果是"用新的猜测修复旧的猜测"——每轮修订可能引入新的未验证声明。评分下降（630→585）正是因为 reviser 修复了问题 A 但引入了问题 B。

**Level 3 — 结构性问题无法通过迭代修复**：Surface Fitness 和 Semantic Purity 的扣分是结构性的——journey 格式将 derived outcomes 写成独立 Steps 而非同一 Step 的多 Outcome，描述的是人工开发工作流而非自动化 CLI 测试。这类问题需要重写 journey 骨架，而非局部修补。迭代式修订对结构性问题无效。

**Level 4（根源）— Pipeline 对纯重构 feature 强制走 journey 流程**：`unify-enum-constants` 是纯内部重构（字符串字面量→类型常量），PRD 明确要求"行为零变更"。从 CLI 用户视角看，命令的输入输出完全不变——没有可观测的新行为，就不存在可测试的 journey。但 pipeline 把 `gen-journeys` 作为所有 feature 的必经步骤，gen-journeys 机械地将每个 PRD user story 1:1 映射为 journey，没有判断"这个 story 是不是用户可观测的工作流"。PRD 的 4 个 stories（编译期类型安全、常量集中定义、枚举完整性、验证 map 合并）主语都是"CLI 开发者/维护者"，描述的是代码质量改进而非用户交互场景，可被 `go build` + `go test` 直接验证，无需端到端 journey。

## Solution

1. **Pipeline 增加 journey 适应性检查**：gen-journeys 在执行前先判断 feature 类型——如果 feature 是纯内部重构（行为零变更、无新 CLI 命令、无新 API 端点），跳过 journey 生成，用 `go build` + `go test` 替代验证。判定信号：PRD 中出现"零行为变更"、"重构"、"内部质量"等关键词，且所有 user stories 的主语是"开发者/维护者"而非"用户"。
2. **生成前必须验证**：gen-journeys 在写入任何事实声明（常量名、CLI 命令、文件路径、计数）之前，必须先 grep/read 实际代码确认。未验证的声明必须标记为 `source: inferred`。
3. **修订必须回查代码**：reviser 收到 adversary 反馈后，应在修改 journey 之前先用代码验证正确的值，而非凭 adversary 描述推测。
4. **区分"可修补"与"需重写"**：adversary 反馈如果涉及结构性问题（格式、骨架、策略比例），应标记为 `structural`，触发重新生成而非局部修订。局部修订只适用于事实性错误和措辞调整。
5. **Invariant 必须反映真实代码**：编写 journey invariant 时先检查代码中是否有例外情况（如通配符模式 `types.Status("*")`），在 invariant 中明确排除这些例外。

## Reusable Pattern

当 eval 循环出现"评分先升后降"模式时，立即停止迭代，区分问题类别：
- **事实错误**（常量名、命令名）→ 单次修复即可，但必须验证代码
- **结构缺陷**（格式、骨架、策略比例）→ 停止迭代，重新设计骨架
- **Invariant 矛盾** → 先修正 invariant 与代码的对齐，再继续

不要用"多轮迭代"解决"生成器本身不查代码"的问题——问题在输入端，不在迭代端。

更根本的判断：**纯重构类 feature 不需要 journey**。如果 feature 的 PRD 明确要求"行为零变更"，且 user stories 的主语是"开发者"而非"用户"，那么 journey 是错误的测试策略——应该用 `go build` + `go test` + `go vet` 验证编译期保证，而非强行构造端到端用户旅程。

## Example

```
# 错误模式：生成器凭文档推测
"117 Status string literals across 22 files"  # 未验证

# 正确模式：生成器先查代码
grep -r '"pending"\|"in_progress"\|"completed"...' → 实际计数 → 写入 journey 并标注 source
```

## Related Files

- `docs/features/unify-enum-constants/testing/eval/eval-journey-report.md`
- `docs/features/unify-enum-constants/testing/eval/iteration-3.md`

## References

- `docs/reference/test-type-model.md` — CLI Surface 测试要求
