---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["9"]
type: "doc.review"
scope: "all"
---

Review documentation quality for the v3-release-audit feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-split-gen-test-scripts-skill
- [ ] SKILL.md ≤ 350 行
- [ ] 新 rules/ 文件被 SKILL.md 通过 Load 引用（入度 ≥ 1）
- [ ] 拆分后 SKILL.md 流程完整，无断裂引用
- [ ] `wc -l plugins/forge/skills/gen-test-scripts/SKILL.md` ≤ 350


### 10-update-pipeline-docs
- [ ] guide.md Pipeline 描述与当前 skill 流程（/write-prd → /tech-design → /breakdown-tasks 等）一致
- [ ] forge-distribution.md Pipeline 描述与当前分发机制一致
- [ ] 无 v2 时代命令名残留


### 11-add-subsystem-summaries
- [ ] 9 个子系统各有独立概述段落
- [ ] 每段落含：架构角色描述 + SKILL.md 链接
- [ ] 新增内容 ≤ 180 行
- [ ] 每子系统 ≤ 20 行


### 12-template-consistency
- [ ] 所有模板变量使用统一格式
- [ ] 模板 frontmatter 字段完整且格式一致


### 13-hook-and-misc-fixes
- [ ] Hook 参数格式已审查，问题已记录或修复
- [ ] validate-ux-pipeline.md 位于 rules/ 目录
- [ ] 未暴露 skill 的 command 入口已评估，决策已记录


### 14-spec-extraction-and-fixes
- [ ] consolidate-specs SKILL.md 体积减小 ≥50 行
- [ ] UTF-8 字符处理问题已评估，影响已记录
- [ ] CLI 测试中无过时 skill 路径引用


### 15-readme-v3-features
- [ ] README 含 v3.0.0 新特性段落
- [ ] 列举的主要特性与 forge-cli-v3 proposal 和 v3 实际变更一致
- [ ] 段落长度适中（~20-30 行），不喧宾夺主


### 2-split-eval-skill
- [ ] SKILL.md ≤ 350 行
- [ ] 新 rule 被 SKILL.md 通过 Load 引用（入度 ≥ 1）
- [ ] 拆分后 SKILL.md 流程完整，无断裂引用
- [ ] `wc -l plugins/forge/skills/eval/SKILL.md` ≤ 350


### 3-remove-harness-prerequisite
- [ ] 所有 SKILL.md Prerequisites 表不含 harness 类型
- [ ] harness 相关文件保留不删除
- [ ] `grep -ri "harness" plugins/forge/skills/*/SKILL.md` 仅剩非 Prerequisites 上下文引用（如有）


### 4-fix-architecture-md
- [ ] Agent 计数与 `ls plugins/forge/agents/ | wc -l` 一致
- [ ] Hook 列表与 `ls plugins/forge/hooks/` 一致
- [ ] Skill 计数与 `ls plugins/forge/skills/ | wc -l` 一致
- [ ] 无 PostToolUse 引用（`grep -c "PostToolUse" docs/ARCHITECTURE.md` = 0）
- [ ] 路径引用指向实际存在的目录


### 5-fix-cli-references
- [ ] `grep -r "forge config get surface" plugins/forge/` 返回 0 结果
- [ ] `grep -r "test\.execution" plugins/forge/` 返回 0 结果（或仅 config-schema 定义性引用）
- [ ] 所有替换命令与 `forge --help` 输出一致


### 6-rewrite-readme
- [ ] 版本号 = `cat plugins/forge/scripts/version.txt` 或当前 RC 版本
- [ ] 技能计数 = `ls plugins/forge/skills/ | wc -l`
- [ ] 任务类型表覆盖所有 dot-notation 类型
- [ ] 命令速查与 `forge --help` 一一对应
- [ ] 无幽灵命令（如已移除的 web/raycast 引用）
- [ ] 安装步骤指向正确 Go 版本要求
- [ ] 路径引用与实际目录匹配


### 7-complete-cli-flags
- [ ] 所有 `forge <subcommand> --help` 输出的 flags 在文档中有对应描述
- [ ] 无文档中存在但 `--help` 不输出的幽灵 flags


### 8-fix-paths-and-orphans
- [ ] `grep -r "hardcoded/path/pattern" plugins/forge/skills/run-tests/` 返回 0（具体模式待确认）
- [ ] 所有 rules/ 文件被至少一个 SKILL.md 引用（入度 ≥ 1）
- [ ] 6 个真孤儿 rules 已添加 Load 指令
- [ ] 5 个参数化 surface rules 已标注引用关系


### 9-dead-code-cleanup
- [ ] `grep -r "sitemap-example" plugins/forge/` 返回 0（排除删除操作本身）
- [ ] init-justfile 6 个 .just 模板已评估，使用状态已记录
- [ ] 删除的文件通过 `grep -r` 全仓库确认无引用


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/v3-release-audit/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/v3-release-audit/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
