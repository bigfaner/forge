---
created: 2026-05-16
author: "fanhuifeng"
status: Draft
---

# Proposal: Quick Mode Test Pipeline Slim -- 合并生成与运行，缩减任务链

## Problem

Quick mode 的测试管线对小型 feature 来说过重。核心矛盾：

**6 个串行测试任务服务于 3-5 个业务任务的 feature，测试流程比业务实现更重。**

### 具体表现

#### 1. 任务数量倒挂

当前 quick mode 的测试任务链（go-test 单 profile）：

```
T-quick-1: gen-test-cases        (30min-1h)
T-quick-2: gen-test-scripts      (30min-1h)
T-quick-3: run-e2e-tests         (15-30min)
T-quick-4: graduate-tests        (15min)
T-quick-5: verify-regression     (15min)
T-quick-6: drift-detection       (15min)
```

一个 3 个业务任务的小 feature，测试任务有 6 个。总 pipeline：9 个 task，测试占比 67%。

#### 2. Subagent 切换开销

每个 task 由独立 subagent 执行。切换代价：读 manifest → 读 task .md → 重建代码上下文 → 执行。每次切换 ~2-5 分钟的 context 重建。6 个测试 task = 6 次切换。

尤其 gen-test-scripts → run-e2e-tests 的切换最浪费：gen 刚刚建立了完整的代码上下文，run 需要完全重建同一上下文来理解生成的测试脚本。

#### 3. gen-scripts 和 run-tests 天然连续

gen-test-scripts 生成脚本后立即需要验证。当前拆成两个 task 的唯一理由是"每个 task 做一件事"，但这对 quick mode 是过度拆分。Quick mode 的哲学是效率优先。

### Evidence

- 最近的 quick mode feature（`cli-list-reverse-chronological`）有 5 个业务 task + 6 个测试 task = 11 个总 task
- gen-test-scripts 执行完后，run-e2e-tests 的 subagent 需要重新读取刚生成的测试脚本来理解上下文
- T-quick-4 (graduate) 对 go-test profile 来说是空操作（go-flat-staging 提案已论证）

### Urgency

go-flat-staging 提案将取消 flat profile 的毕业任务。趁此机会同步简化 quick mode 的任务结构，避免两次修改同一文件（`testgen.go`）。两个提案一起落地，测试管线从 6 步缩减到 4 步（flat profile）或 5 步（nested profile）。

## Proposed Solution

**核心思路**：将 gen-test-scripts 和 run-e2e-tests 合并为单个 task。生成脚本后立即运行测试，失败在同一 task 内修复重跑。

### 方案细节

#### 1. 合并 T-quick-2 + T-quick-3 为 "Generate & Run Tests"

新 task 同时执行两个 skill：
1. 调用 `/gen-test-scripts` 生成测试脚本
2. 调用 `/run-e2e-tests` 运行测试
3. 如果测试失败，在同一 subagent 会话内修复并重跑

**新增任务类型**：`test-pipeline.gen-and-run`，对应新 prompt template。

#### 2. 重编号 quick mode 任务链

**合并后（flat profile，配合 go-flat-staging）**：

```
T-quick-1: gen-test-cases           (不变)
T-quick-2: gen-scripts + run-tests  (合并)
T-quick-3: verify-regression        (原 T-quick-5)
T-quick-4: drift-detection          (原 T-quick-6)
```

**合并后（nested profile，保留毕业）**：

```
T-quick-1: gen-test-cases           (不变)
T-quick-2: gen-scripts + run-tests  (合并)
T-quick-3: graduate-tests           (原 T-quick-4)
T-quick-4: verify-regression        (原 T-quick-5)
T-quick-5: drift-detection          (原 T-quick-6)
```

#### 3. Per-type 模式的合并

当 test-cases.md 检测到多种类型（如 tui + api）时：

```
T-quick-1: gen-test-cases
T-quick-2-tui: gen-scripts + run-tests (tui)   ← 独立生成并验证
T-quick-2-api: gen-scripts + run-tests (api)   ← 独立生成并验证
T-quick-3: verify-regression                   ← 依赖所有 T-quick-2-*
T-quick-4: drift-detection
```

每个 per-type task 自包含（生成 + 运行），verify-regression 等待所有类型完成后做集成验证。

#### 4. 上下文利用

合并后的 subagent 持有 gen-test-scripts 建立的完整上下文：
- 代码结构理解（来自 reconnaissance）
- 刚生成的测试脚本（已在内存中）
- Profile 策略（已解析）

run-e2e-tests 直接复用这些上下文，无需重建。测试失败时，修复循环在同一上下文内完成。

#### 5. 与 go-flat-staging 的交互

两个提案修改同一函数 `GetQuickTestTasks`，修改点不重叠：

| 提案 | 修改点 |
|------|--------|
| go-flat-staging | 当 staging-mode=flat 时，跳过 graduate 任务生成 |
| 本提案 | 合并 gen-scripts + run 为单任务，重编号 |

**应用顺序**：先 go-flat-staging（取消毕业），再本提案（合并生成运行）。或同时应用。

**最终效果对比**：

| Profile | 当前 | go-flat-staging 后 | 本提案后 | 两者都应用后 |
|---------|------|-------------------|---------|------------|
| flat (go-test) | 6 步 | 5 步 | 5 步 | **4 步** |
| nested (web-playwright) | 6 步 | 6 步 | 5 步 | **5 步** |

### Developer Walkthrough

#### Before

```
Task executor: T-quick-2 (gen-test-scripts)
  → Read manifest, read codebase, generate scripts
  → Subagent exits

Task executor: T-quick-3 (run-e2e-tests)
  → Read manifest, read generated scripts (rebuild context), run tests
  → FAIL: assertion error
  → Read generated script again, fix, re-run
  → PASS
  → Subagent exits
```

2 个 subagent，2 次 context 重建。

#### After

```
Task executor: T-quick-2 (gen-scripts + run-tests)
  → Read manifest, read codebase, generate scripts
  → Run tests (context already in memory)
  → FAIL: assertion error
  → Fix script (context already in memory), re-run
  → PASS
  → Subagent exits
```

1 个 subagent，1 次 context 重建。失败修复在同一会话内完成。

## Requirements Analysis

### Key Scenarios

#### Happy Path

1. **单 profile 无 per-type**：T-quick-2 生成脚本并运行，一次通过。verify-regression 确认无回归。
2. **单 profile per-type（tui + api）**：T-quick-2-tui 和 T-quick-2-api 分别独立生成并运行。全部通过后 T-quick-3 做回归验证。
3. **测试失败修复**：生成后运行失败，在同一 subagent 内读错误、修复脚本、重跑。最多 3 次重试。

#### Error & Edge Cases

4. **生成阶段失败**：gen-test-scripts 生成失败（如无法解析 test-cases.md）。subagent 正常报错，任务标记为 blocked。
5. **运行阶段反复失败**：超过 3 次修复仍失败。任务标记为 blocked，由 run-tasks dispatcher 创建 fix task。
6. **per-type 并发执行**：T-quick-2-tui 和 T-quick-2-api 由不同 subagent 并行执行，各自独立。不会写入同一文件（不同 type 的测试文件独立）。
7. **合并后 verify-regression 依赖变更**：verify-regression 依赖 T-quick-2（merged）而非旧的 T-quick-3（run-only）。依赖链正确指向合并后的 task ID。

### Non-Functional Requirements

- **上下文效率**：合并后的 subagent context 峰值不超过 gen-test-scripts 单独执行的 1.5 倍（因为 run 步骤只增加测试输出和错误信息，不增加代码理解负担）
- **向后兼容**：已存在的旧编号任务文件（如正在执行中的 T-quick-3）通过 `infer.go` 的 type 映射继续正常工作。新编号仅对新创建的 feature 生效
- **pipeline 加速**：subagent 切换次数减少 1 次（flat: 6→4, nested: 6→5），每次切换节省 2-5 分钟

### Constraints & Dependencies

- 新增任务类型常量 `TypeTestPipelineGenAndRun`（在 `infer.go` 中注册）
- 新增 prompt template `test-pipeline-gen-and-run.md`
- `testgen.go` 的 `GetQuickTestTasks` 函数签名不变（仍接收 `profiles []string, detectedTypes []string`）
- 依赖 `infer.go` 的 type→prompt 映射机制

## Alternatives & Industry Benchmarking

### Industry Solutions

多数 CI/CD 系统支持 stage 合并以减少 agent 分配开销：

- **GitHub Actions**：通过单一 job 内多 step 避免跨 job 的 artifact 传递开销
- **Bazel**：通过 single executable 测试策略在相同进程内完成编译+运行+验证
- **Go test**：`go test` 命令本身是 compile+link+run 的单命令融合

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动 | 测试占比 67%，执行慢 | Rejected: 问题明确且与 go-flat-staging 有协同优化窗口 |
| **Merge gen+run (recommended)** | CI/CD stage merging | 减少 1 次切换，上下文复用，修复循环更快 | 合并 task 的 prompt 更长，需新类型常量 | **Selected: 直接解决切换开销问题** |
| 合并 prompt 但拆 task | 自定义 | 保持 task 粒度 | 仍然是 2 个 subagent 切换，核心问题未解 | Rejected: 假合并 |
| 保持独立但加 context cache | Claude Code | 最小改动 | 依赖平台特性，不可控 | Rejected: 平台耦合 |

## Feasibility Assessment

### Technical Feasibility

改动集中在 `forge-cli` 的 3 个文件：

| 文件 | 改动 |
|------|------|
| `pkg/task/testgen.go` | 合并 gen+run 为单 task 定义，调整依赖链 |
| `pkg/task/infer.go` | 新增 `TypeTestPipelineGenAndRun` 类型映射 |
| `pkg/prompt/data/test-pipeline-gen-and-run.md` | 新 prompt template，调用两个 skill |

### Resource & Timeline

1 个任务即可完成：修改 testgen.go + infer.go + 新增 prompt template + 更新测试。估计 1-2 小时。

### Dependency Readiness

- `infer.go` 的 type→template 映射机制已成熟（`prompt.go` line 22-37）
- `TestTaskDef` struct 已有 `Type` 字段，无需扩展
- `resolveQuickDeps` 函数可直接修改依赖链逻辑

## Scope

### In Scope

- 合并 gen-test-scripts + run-e2e-tests 为单 task（quick mode only）
- 新增任务类型 `TypeTestPipelineGenAndRun`
- 新增合并 prompt template
- 重编号 quick mode 任务链（flat: 4步, nested: 5步）
- 更新 `testgen_test.go` 测试用例

### Out of Scope

- breakdown mode 的任务结构变更（full mode 保持独立 gen/run）
- graduation 任务移除（由 go-flat-staging 处理）
- gen-test-scripts 或 run-e2e-tests skill 本身的修改
- per-type 检测逻辑变更

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 合并 task 的 prompt 过长导致执行质量下降 | L | M | prompt 结构清晰（先 gen 后 run），quick mode 的 feature 本身小，生成内容有限 |
| 上下文峰值过高（gen 的 reconnaissance + run 的输出） | L | L | Quick mode feature 小，reconnaissance 范围有限。实际测量：当前 gen 单独执行的 context 足以覆盖合并后的需求 |
| 与 go-flat-staging 的合并冲突 | M | M | 两个提案修改 testgen.go 的不同区域（go-flat-staging 改 graduate 跳过逻辑，本提案改 gen/run 合并逻辑）。建议同时应用 |
| 旧编号任务文件兼容性 | L | L | 已存在的 feature 使用旧编号，由 infer.go 的现有映射处理。新 feature 使用新编号 |
| per-type 并行执行时的 helpers.go 冲突 | L | M | 仅影响 flat staging（共享 helpers.go）。gen-test-scripts 的 merge 机制已处理冲突检测 |

## Success Criteria

- [ ] Quick mode 生成 4 个测试 task（flat profile）或 5 个（nested profile），不再有 6 个
- [ ] T-quick-2 的 prompt 同时调用 `/gen-test-scripts` 和 `/run-e2e-tests`
- [ ] 合并 task 的 subagent 在一次会话内完成生成和运行（含失败修复）
- [ ] verify-regression 正确依赖 T-quick-2（merged）而非独立 run task
- [ ] drift-detection 作为最后一步，依赖 verify-regression
- [ ] per-type 模式下 T-quick-2-tui 和 T-quick-2-api 各自独立生成并运行
- [ ] `testgen_test.go` 新增测试验证合并后的 task 数量和依赖链
- [ ] 与 go-flat-staging 同时应用后，flat profile quick mode 为 4 步

## Next Steps

- 通过 `/quick` 流程生成实施任务（预计 1-2 个任务）
- 或直接实施（预计 1-2 小时）
