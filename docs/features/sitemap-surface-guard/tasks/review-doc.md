---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["1", "2", "3", "4", "5"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the sitemap-surface-guard feature (quick mode).

## Acceptance Criteria
- [ ] 所有 doc task 的产出物已扫描并对照 AC 逐项验证
- [ ] 守卫措辞一致性已检查：Task 1-4 的 surface 检测指令模式统一
- [ ] Review 结果记录到 tasks/records/

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-gen-sitemap-surface-check
- [ ] Step 0 执行 `forge surfaces --json`，解析返回的 surface 列表
- [ ] 无 `web` surface 类型时，STOP 并输出明确提示（如 "No web surface detected. gen-web-sitemap is only applicable to web projects."）
- [ ] 有 `web` surface（含 monorepo 多 surface 场景）时，Step 0 通过，正常进入 Step 1
- [ ] `forge surfaces --json` 返回空或命令失败时，等同于无 web surface，中止执行


### 2-write-prd-sitemap-guard
- [ ] SKILL.md Step 1 读取 sitemap.json 前，先检查项目是否有 web surface（通过 `forge surfaces --json`），无则跳过读取
- [ ] ui-functions.md Placement Rules 读取 sitemap 前，增加 web surface 前置条件检查
- [ ] self-check.md 的 Placement consistency 检查增加 web surface 前置条件；Sitemap availability 检查同样受 surface 条件守卫


### 3-breakdown-tasks-surface-guard
- [ ] Placement Validation 段落在检查 sitemap.json 存在性之前，先检查项目是否有 web surface
- [ ] 无 web surface 时，跳过 route 验证并输出适当提示（而非警告 sitemap 缺失）


### 4-eval-surface-guard
- [ ] Web 行的 sitemap.json 引用增加显式守卫：仅在项目有 web surface 时使用 sitemap 作为辅助信息
- [ ] 非 web 类型（CLI、TUI）行的辅助信息不含 sitemap 相关内容，确认无遗漏


### 5-audit-ui-test-guard
- [ ] 审查已有守卫 "Only execute when the project has `web-ui` interface" 是否与 `forge surfaces --json` 检测方式一致
- [ ] 记录审查结论：守卫充分无需修改，或实施修改并记录变更内容
- [ ] 若修改，守卫措辞与 Task 1-4 中建立的模式保持一致


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/sitemap-surface-guard/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/sitemap-surface-guard/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
