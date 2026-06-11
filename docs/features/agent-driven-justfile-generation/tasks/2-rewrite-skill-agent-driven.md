---
id: "2"
title: "Rewrite SKILL.md for agent-driven justfile generation"
priority: "P0"
estimated_time: "2h"
dependencies: [1]
type: "doc"
mainSession: false
# Note: surface-key and surface-type fields are intentionally absent from doc tasks.
# Doc tasks produce non-compilable output (markdown, specs, templates) and do not
# interact with the quality gate or test pipeline, so surface routing is unnecessary.
---

# 2: Rewrite SKILL.md for agent-driven justfile generation

## Description

Rewrite `SKILL.md` to remove template-driven generation and replace with agent-driven generation. The agent directly detects language/framework via marker files, reads recipe contracts from surface rules, acquires framework knowledge from Convention files, and generates concrete recipe commands — no language templates needed. This is the core architectural change of the proposal.

Key changes: remove `--type` parameter, remove project type detection step (Step 1a + `rules/project-detection.md` reference), make `forge surfaces` a prerequisite (empty → prompt user to run `forge init`), replace Step 0's template loading with agent-driven language detection + Convention loading, replace Step 3's template-based generation with agent-driven generation referencing `rules/server-lifecycle.md`, and add post-generation consistency verification.

## Reference Files
- `docs/proposals/agent-driven-justfile-generation/proposal.md` — Proposed Solution, Key Scenarios, Non-Functional Requirements, Scope > In Scope, Success Criteria (ref: ## Proposed Solution, ## Requirements Analysis > ### Key Scenarios, ## Scope > ### In Scope, ## Success Criteria)
- `plugins/forge/skills/init-justfile/SKILL.md` — full file being rewritten (ref: # /init-justfile)
- `plugins/forge/skills/init-justfile/rules/server-lifecycle.md` — new server lifecycle rule from Task 1
- `docs/conventions/forge-distribution.md` — plugin path resolution conventions (ref: ## 5. 路径解析机制)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/init-justfile/SKILL.md` | Core rewrite: remove template-driven, add agent-driven generation flow |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `--type` 参数已移除：frontmatter `argument-hint` 更新、Parameters 表移除 `--type` 行、全文无 `--type` 引用
- [ ] 项目类型检测步骤已移除：Step 1a（project type detection）删除、`rules/project-detection.md` 引用删除、`FRONTEND_DIR`/`BACKEND_DIR`/`BACKEND_ENTRY`/`FRONTEND_RUN_SCRIPT` 检测逻辑删除
- [ ] surfaces 为前提条件：Step 1s 的 Outcome B（无 surfaces）改为提示用户运行 `forge init` 配置 surfaces，而非静默跳过
- [ ] Step 3 改为 agent 驱动生成：agent 根据 marker files（`go.mod`/`package.json`/`Cargo.toml`/`pyproject.toml`/`pom.xml`/`build.gradle` 等）检测语言，从 surface rule 读取 recipe 契约，从 Convention 获取框架知识，直接生成 recipe 命令体
- [ ] 引用 `rules/server-lifecycle.md` 替代模板中的 server lifecycle 代码，Step 0 HARD-RULE 中的模板加载流程已移除
- [ ] 包含 post-generation 一致性验证步骤：生成 justfile 后，自动比对 surface rule 文件中的 Recipe Invocation Contract 与实际生成的 recipe 名称/参数，确保 init-justfile 和 run-tests 双消费者一致性

## Hard Rules
- 保持向后兼容：生成的 justfile 结构（boundary markers、recipe 命名、[linux]/[windows] 双平台、user-customized 标记、退出码语义）与当前输出一致
- `# user-customized` 标记的 recipe 在 re-generation 时必须保留（已有行为不变）
- 仅修改 `plugins/forge/skills/init-justfile/SKILL.md` 一个文件

## Implementation Notes

### 结构级一致性要求（来自 proposal Non-Functional Requirements）
- 相同项目多次运行，recipe 名称、分组（group marker）、边界标记（boundary marker）、退出码语义必须完全相同
- 命令体允许因 LLM 变化而不同，但必须满足相同的语义契约（输入→输出→副作用）

### agent 驱动生成的知识来源优先级
1. Surface rule 文件 → recipe 契约（名称、参数、exit code 语义）
2. Convention 文件 → 框架特定知识（test runner、build tool、lint tool）
3. Agent 自身知识 → 主流语言/框架的通用命令模式

### 回退策略
- 空 surfaces → 提示 `forge init`
- 无 Convention 文件 → agent 使用主流语言默认知识（冷启动）
- 罕见语言/框架 → agent 生成 error stub（与当前 generic.just 行为一致）

### 文件范围
仅修改 `plugins/forge/skills/init-justfile/SKILL.md`，不涉及其他文件。
