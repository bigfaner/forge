---
created: "2026-05-10"
updated: "2026-05-10"
author: faner
status: Draft
source: disc-4-stuck forensics
scope: forge-harness
---

# Proposal: task-executor 骨架化 — Execution Workflow 统一所有任务类型

## Problem

task-executor 硬编码了 TDD 工作流（Step 2-3），所有非 MAIN_SESSION 任务走同一条路径。但任务类型多样，执行型任务（如 T-test-3 "Run e2e Tests"）被强行塞入 TDD 循环，导致 14 分钟无效重试。

同时，`noTest` 标志作为权宜之计混入模板 frontmatter，含义模糊——它不表示"不需要测试"，而是"不走 TDD 循环"。这种语义不匹配导致新模板的 `noTest` 设置常出错。

### Evidence

- Forensic analysis (disc-4-stuck): T-test-3 被派发到 task-executor 后，agent 的 TDD 本能驱使其反复重试失败测试，而非分析失败并创建 fix task
- `noTest: true` 出现在"生成文档"（T-test-1）和"汇总记录"（phase-summary）等完全不同类型的任务上，表明该字段无法准确描述任务行为
- T-test-2 (`gen-test-scripts.md`) 设置了 `noTest: false`，但它也不是 TDD 任务——它生成代码后验证编译，不写测试先

### Urgency

每次执行 T-test-3/T-test-4.5 浪费约 14 分钟（对比 Execution Workflow 预期耗时 <5 分钟）。当前有 4 个非 TDD 任务模板受影响（T-test-3, T-test-4, T-test-4.5, T-test-5），graduate/consolidate-specs 将持续增加。`noTest` 的歧义会导致新模板配置错误。

## Proposed Solution

**两条主线：Execution Workflow 统一 + noTest 移除。**

### 主线 1: Execution Workflow 统一所有任务类型

每个任务模板声明 `## Execution Workflow` 段落，task-executor 读取并执行。不再硬编码 TDD：

```
BEFORE:
  task-executor Step 2: TDD (硬编码)
  task-executor Step 3: Quality Gate (硬编码)
  noTest=true → 跳过 Step 2-3（但 agent 仍会困惑）

AFTER:
  task-executor Step 2: 读取 ## Execution Workflow → 执行
  无 ## Execution Workflow → 回退 TDD（向后兼容旧任务）
  noTest 不存在 → 由 workflow 自身决定是否跑测试
```

task-executor 变成纯骨架：读任务 → 执行 workflow → 记录 → 提交。

#### Execution Workflow 解析机制

task-executor 使用 markdown 标题检测提取 workflow 段落：

1. **段落检测**: 在任务 markdown 中搜索 `## Execution Workflow` 标题行。使用标题层级（`##`）作为哨兵——标题行之后、下一个同级或更高级标题之前的所有内容即为 workflow 正文。
2. **注入方式**: task-executor 将提取的 workflow 正文原样注入到 agent prompt 的 Step 2 指令区域，替换当前硬编码的 TDD 步骤。agent 直接将这段 markdown 作为执行指令。
3. **回退逻辑**: 若任务文件不包含 `## Execution Workflow` 标题，task-executor 回退到当前 TDD + Quality Gate 步骤（向后兼容）。
4. **异常处理**: 若标题存在但正文为空（标题后紧接下一个 `##` 标题或文件结束），task-executor 视为配置错误，记录警告到执行记录并回退到 TDD，而非静默跳过。

### 主线 2: 移除 noTest

`noTest` 的职责被 `## Execution Workflow` 完全取代。移除范围：

- task-cli: `types.go` 删除字段、`record.go` 删除相关逻辑、claim 输出不再包含 `NO_TEST`
- agents/task-executor.md: 删除 `NO_TEST` input 和所有 `NO_TEST` 条件分支
- commands/run-tasks.md: 删除 `NO_TEST` 从 claim 输出解析和 dispatch prompt
- commands/execute-task.md: 删除 `NO_TEST` 相关内容
- 所有任务模板: 删除 `noTest` frontmatter 字段
- index.schema.json (breakdown + quick): 删除 `noTest` 字段定义

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 零成本 | T-test-3 每次浪费 ~14 分钟，noTest 歧义持续 | Rejected: 成本太高 |
| 只加 Execution Workflow，保留 noTest | 改动最小 | noTest 变成僵尸字段，新模板仍会困惑 | Rejected: 半吊子方案 |
| 新建 execution-type agent | 完全隔离 | 多一个 agent 维护，重复 Record/Commit 逻辑 | Rejected: 重复代码 |
| dispatcher 路由 + prompt 模板 | dispatcher 控制 | dispatcher 需改，多一条维护线 | Rejected: 复杂度高 |
| **Execution Workflow + 移除 noTest** | 干净、统一、可扩展 | 改动面广（20+ 文件） | **采用** |

## Scope

### In Scope

- `agents/task-executor.md`: Step 2-3 合并为"执行 workflow"（从任务文件读取），删除 NO_TEST input
- 所有 breakdown-tasks 模板（10 个，不含 manifest-update-tasks.md 和 eval-test-cases.md）: 添加 `## Execution Workflow`，删除 `noTest`
- 所有 quick-tasks 模板（6 个，不含 manifest-quick.md）: 添加 `## Execution Workflow`，删除 `noTest`
- `index.schema.json` (breakdown + quick): 删除 `noTest` 字段
- `commands/run-tasks.md`: 删除 NO_TEST 从 claim 解析和 dispatch
- `commands/execute-task.md`: 删除 NO_TEST 相关内容
- `task-cli`: 删除 noTest 字段和相关逻辑
- `skills/record-task/SKILL.md`: 删除 noTest 引用
- `skills/quick-tasks/SKILL.md`: 删除 --no-test 标志
- `skills/consolidate-specs/SKILL.md`: 删除 noTest 引用

### Out of Scope

- 模板内容重写（仅添加 Execution Workflow 段落，不改动现有 Implementation Notes）
- mainSession 任务的路由逻辑（已由 dispatcher 处理，不涉及 task-executor）
- `## Execution Workflow` 的模板化/标准化（未来优化，将常用 workflow 提取为可引用片段）
- `task add --template` 自动注入 Execution Workflow（未来优化）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| task-executor 不遵循 Execution Workflow | Low | High | `<EXTREMELY-IMPORTANT>` / `<HARD-RULE>` 标签包裹指令："你必须严格遵循任务文件中 `## Execution Workflow` 段落指定的步骤。该段落定义了你要运行哪些命令、什么算成功、失败时如何处理。不得偏离指定的工作流。" 此机制与 "ONE TASK PER INVOCATION" 规则一致，后者自部署以来零违规。若 workflow 缺失则回退到 TDD。 |
| noTest 移除引入 task-cli bug | Medium | High | task-cli 有单元测试覆盖；移除后跑全量测试 |
| 20+ 文件改动导致遗漏 | Medium | Medium | grep 确认零匹配 + 对每个含过 noTest 条件分支的文件（task-cli types.go/record.go, task-executor.md, run-tasks.md, execute-task.md）逐一审查：条件逻辑中不再有基于 noTest/NO_TEST 的分支，无死代码残留 |
| 旧任务文件缺少 Execution Workflow | Low | Low | task-executor 有 fallback：无 workflow → TDD |
| Execution Workflow 写得不好导致 agent 迷失 | Medium | Medium | 每个 workflow 模板合并前须通过 dry-run 测试；模板审核清单：workflow 段落须包含明确的结束条件，不允许开放式指令如"继续测试直到通过" |

## Success Criteria

- [ ] task-executor Step 2 从任务文件读取 `## Execution Workflow` 并执行
- [ ] 无 `## Execution Workflow` 的旧任务回退到 TDD + Quality Gate（行为不变）
- [ ] 所有 16 个任务模板（10 breakdown + 6 quick）包含 `## Execution Workflow` 段落
- [ ] `noTest` 从所有文件中完全移除（grep 确认零匹配 + 语义审查无残留分支）
- [ ] task-cli 编译通过 + 单元测试通过
- [ ] dispatcher (`run-tasks.md`) 不再引用 NO_TEST
- [ ] `commands/execute-task.md`: NO_TEST 相关内容已删除，execute-task 流程不传递 NO_TEST 参数仍正常运行
- [ ] 手动运行 T-test-3，检查执行记录——Step 2 输出行必须包含 'Execution Workflow'，不得包含 'TDD implementation' 或 'RED/GREEN/REFACTOR' 关键词；执行时间 < 5 分钟（对比当前 ~14 分钟）
- [ ] `skills/record-task/SKILL.md`：noTest 引用已删除，skill 调用后仍能正确创建执行记录
- [ ] `skills/quick-tasks/SKILL.md`：`--no-test` 标志已删除，quick-tasks 流程不传该标志仍正常运行
- [ ] `skills/consolidate-specs/SKILL.md`：noTest 引用已删除，skill 调用后仍能正确提取和合并规格

## Next Steps

- Proceed to `/write-prd` to formalize requirements
