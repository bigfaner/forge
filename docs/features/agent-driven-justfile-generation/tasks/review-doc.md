---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["1", "2", "3", "4"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the agent-driven-justfile-generation feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-create-server-lifecycle-rule
- [ ] `plugins/forge/skills/init-justfile/rules/server-lifecycle.md` 文件已创建
- [ ] 包含 PID 追踪模式：PID 文件路径约定（`.forge/<surfaceKey>.pid`）、原子写入、stale PID 检测（进程不存在时的清理逻辑）
- [ ] 包含幂等启动模式：检测已有进程（避免重复启动）、端口占用检查（检测后选择备选端口或报错）、启动/重启逻辑
- [ ] 包含健康检查模式：HTTP/TCP probe 实现、重试策略（最多 3 次、5 秒间隔）、超时处理
- [ ] 包含 multi-service 场景指导：per-service PID 文件隔离、端口感知启动顺序、启动顺序依赖声明
- [ ] 提供可直接使用的 bash 代码片段（带插槽占位符如 `<PORT>`、`<START_CMD>`），agent 优先复用而非从头生成


### 2-rewrite-skill-agent-driven
- [ ] `--type` 参数已移除：frontmatter `argument-hint` 更新、Parameters 表移除 `--type` 行、全文无 `--type` 引用
- [ ] 项目类型检测步骤已移除：Step 1a（project type detection）删除、`rules/project-detection.md` 引用删除、`FRONTEND_DIR`/`BACKEND_DIR`/`BACKEND_ENTRY`/`FRONTEND_RUN_SCRIPT` 检测逻辑删除
- [ ] surfaces 为前提条件：Step 1s 的 Outcome B（无 surfaces）改为提示用户运行 `forge init` 配置 surfaces，而非静默跳过
- [ ] Step 3 改为 agent 驱动生成：agent 根据 marker files（`go.mod`/`package.json`/`Cargo.toml`/`pyproject.toml`/`pom.xml`/`build.gradle` 等）检测语言，从 surface rule 读取 recipe 契约，从 Convention 获取框架知识，直接生成 recipe 命令体
- [ ] 引用 `rules/server-lifecycle.md` 替代模板中的 server lifecycle 代码，Step 0 HARD-RULE 中的模板加载流程已移除
- [ ] 包含 post-generation 一致性验证步骤：生成 justfile 后，自动比对 surface rule 文件中的 Recipe Invocation Contract 与实际生成的 recipe 名称/参数，确保 init-justfile 和 run-tests 双消费者一致性


### 3-simplify-surface-rules
- [ ] 所有 5 个 surface rule 文件（`api.md`/`cli.md`/`tui.md`/`web.md`/`mobile.md`）已更新
- [ ] `## Recipe Template (Dual Platform)` section（含 TODO stub 代码块和 `<Test-Dir-Path>` 块）已替换为 `## Recipe Generation Requirements` section
- [ ] 以下 section 保留不变：Orchestration Sequence、Recipe Invocation Contract、Journey Filter Strategy
- [ ] `## Recipe Generation Requirements` section 包含：agent 生成 recipe 时需遵循的结构约束（recipe 命名规则、`[linux]`/`[windows]` 双平台属性、`# user-customized` 标记、exit code 0/1 语义、test 目录路径规则 single vs multi surface）
- [ ] 双消费者一致性保留：init-justfile（生成 recipe）和 run-tests（消费 recipe）都能从同一份 surface rule 的 Recipe Invocation Contract 获取所需信息


### 4-delete-templates-and-detection
- [ ] 6 个语言模板文件已删除：`go.just`、`node.just`、`python.just`、`rust.just`、`mixed.just`、`generic.just`
- [ ] `rules/project-detection.md` 已删除
- [ ] `templates/` 目录已清空或删除（确认无其他文件残留）
- [ ] grep 验证：`SKILL.md` 和所有 surface rule 文件中无对已删除文件的引用（搜索 `templates/`、`project-detection`、`generic.just` 等关键词）


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/agent-driven-justfile-generation/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/agent-driven-justfile-generation/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.

## Acceptance Criteria

- [ ] All acceptance criteria met
