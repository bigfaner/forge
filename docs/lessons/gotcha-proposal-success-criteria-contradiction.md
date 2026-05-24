---
created: "2026-05-24"
tags: [architecture, testing]
---

# Proposal Success Criteria 内部矛盾导致实现无法同时满足

## Problem

Feature `pipeline-integration-stitch` 的 proposal 通过了 897/1000 的 adversarial eval，但实现完成后发现无法同时满足所有 Success Criteria：

- `grep -r "gen-and-run" forge-cli/ --exclude-dir=docs/proposals` 要求零结果
- 同时 In Scope 要求 `validate_index.go` 保留迁移错误提示（检测到 `test.gen-and-run` 时输出 deprecated 信息）
- 同时 In Scope 要求 `Synthesize()` 对 gen-and-run 文件名返回迁移指引

保留迁移防护 = grep 非零。清除 grep = 删除迁移防护。两个要求互斥。

## Root Cause

**因果链（3 层）**：

1. **L1 — Proposal 写入时未做 Success Criteria 内部一致性校验**：多个 SC 条目由不同 In Scope 条目分别推导，推导过程独立，未交叉验证是否可同时满足。

2. **L2 — "加法思维"导致功能范围与清理范围重叠**：P1 要求"添加迁移防护"（加法），P2 要求"完全移除 gen-and-run 引用"（减法），两者针对同一代码区域但方向相反。Proposal 写作时将两者视为独立条目，未识别冲突。

3. **L3 — Eval rubric 不检测 Success Criteria 的逻辑可满足性**：897/1000 的 adversarial eval 聚焦于 proposal 的完整性、清晰度、行业对标，但不检查 SC 条目之间是否存在逻辑矛盾（类似 SAT solver 的约束冲突检测）。这是 eval 的盲区——它验证"写得好不好"，不验证"能不能做到"。

## Solution

用户最终明确要求"完全去掉 gen-and-run"，选择了**减法优先**策略：删除所有迁移防护代码（`prompt.go` 迁移指引、`validate_index.go` 迁移守卫、相关测试用例），使 grep 真正零结果。

旧 index.json 引用 `test.gen-and-run` 时将收到通用的 "read template" 或 "invalid type" 错误，而非有意义的迁移提示。

## Reusable Pattern

**写 proposal 时，对 Success Criteria 执行"同时满足性检查"**：

1. **SC 冲突检测**：列出所有 SC 条目，标记每个条目对代码的"方向"（加法/减法/修改）。若两个条目方向相反且作用于同一文件，标记为冲突候选。
2. **互斥声明**：若业务上确实需要"同时保留迁移防护 + 清除残留"，在 proposal 中显式声明互斥区域（如 `--exclude-dir` 排除特定文件），而非留给实现者猜测优先级。
3. **SC 可验证性原则**：每个 SC 条目必须能独立通过或失败，不应存在"A 通过则 B 必然失败"的组合。

**对于废弃代码移除类任务**，在 proposal 阶段就决定是"干净切断"还是"保留迁移桥"——两者只能选一个。将决策记录到 proposal 的 Alternatives 表中，而非让两个方向同时出现在 In Scope 中。

## Example

**矛盾 SC 示例**：
```
- [ ] grep -r "gen-and-run" forge-cli/ 返回零结果          ← 减法
- [ ] validate_index.go 对 test.gen-and-run 返回迁移错误   ← 加法（保留字符串）
- [ ] Synthesize() 对 gen-and-run 文件名返回迁移指引        ← 加法（保留字符串）
```

**修正后 SC（选择干净切断）**：
```
- [ ] grep -r "gen-and-run" forge-cli/ 返回零结果
- [ ] （迁移防护移除，旧 index.json 将收到通用错误）
```

**修正后 SC（选择保留迁移桥）**：
```
- [ ] grep -r "gen-and-run" forge-cli/ --exclude=validate_index.go --exclude=prompt.go 返回零结果
- [ ] validate_index.go 对 test.gen-and-run 返回迁移错误
- [ ] Synthesize() 对 gen-and-run 文件名返回迁移指引
```

## References

- `docs/proposals/pipeline-integration-stitch/proposal.md` — 原始 proposal（eval 897/1000）
