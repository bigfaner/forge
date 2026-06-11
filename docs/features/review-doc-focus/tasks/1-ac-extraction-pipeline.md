---
id: "1"
title: "AC extraction pipeline (BodyContext + extract.go + renderBody)"
priority: "P0"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 1: AC Extraction Pipeline

## Description

构建 AC 提取管线：扩展 BodyContext schema 新增 `DocTaskCriteria map[string]string` 字段，在 `extract.go` 中新增 `extractDocTaskCriteria()` 函数遍历 doc 任务文件按 header 提取 AC section，在 `autogen.go` 的 `renderBody()` 中新增 `{{DOC_TASK_AC}}` 占位符通过 `strings.ReplaceAll` 注入序列化后的 AC markdown。

## Reference Files
- `proposal.md#Proposed-Solution` — 定义构建时嵌入 AC 的核心方案和双层防护模型
- `proposal.md#Feasibility-Assessment` — Technical Feasibility 四条主线中的 BodyContext schema 和 AC 提取管线
- `proposal.md#Scope` — In Scope 中 build.go、extract.go、autogen.go 的具体修改要求
- `proposal.md#AC-汇总区域格式` — AC 汇总区域的 markdown schema 和规则定义

## Acceptance Criteria

- [ ] `BodyContext` struct 新增 `DocTaskCriteria map[string]string` 字段（key=任务名称, value=AC markdown）
- [ ] `extract.go` 新增 `extractDocTaskCriteria(taskDir string) map[string]string` 函数
- [ ] 提取算法：逐行扫描 `.md` 文件，找到 `## Acceptance Criteria` 后收集至下一个 `##` 行之前的所有内容
- [ ] `renderBody()` 支持 `{{DOC_TASK_AC}}` 占位符，将 map 序列化为 markdown 后通过 `strings.ReplaceAll` 注入
- [ ] `BuildIndex()` 在生成 review-doc 任务时调用提取函数并填充 BodyContext
- [ ] 渲染策略使用 `strings.ReplaceAll`（非 Go `text/template`），保持与现有模板一致

## Hard Rules

- 使用 `strings.ReplaceAll` 而非 Go `text/template` 渲染，不引入新依赖
- 字段命名 `DocTaskCriteria` 避免与现有 `AcceptanceCriteria []string` 冲突

## Implementation Notes

- 提取函数需处理 section 不存在的情况（返回空 string）
- 逐行扫描需注意 fenced code block 内的 `##` 行不应作为 section 边界
- `renderBody()` 序列化 map 时需按 AC 汇总区域格式的 markdown schema 生成（`### task-name` 子 section）
