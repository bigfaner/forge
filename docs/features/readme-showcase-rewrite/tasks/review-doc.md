---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["1", "2"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the readme-showcase-rewrite feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-write-readme-zh
- [ ] 首屏（前 30 行）包含核心定位 + 痛点共鸣叙述（场景化叙事，非空泛口号）
- [ ] 总长度 ≤ 200 行（当前 ~474 行）
- [ ] 包含功能矩阵对比表，覆盖 Superpowers / Spec Kit / OpenSpec 共 3 个竞品，≥ 6 个对比维度
- [ ] 4 大核心特性（质量门控 / 上下文持久化 / 知识沉淀 / Agent 自动编排）各有独立小节
- [ ] 无 CLI 命令速查表（命令信息由 `forge --help` 承担）
- [ ] 保留现有安装步骤的技术准确性（install 命令、前置依赖不变）


### 2-write-readme-en
- [ ] 英文版内容与中文版对应（非直译），符合英文母语表达习惯
- [ ] 包含与中文版相同的 7 节结构
- [ ] 总长度 ≤ 200 行
- [ ] 竞品功能矩阵与中文版一致（3 个竞品、≥ 6 个维度）
- [ ] 4 大核心特性各有独立小节


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/readme-showcase-rewrite/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/readme-showcase-rewrite/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.

## Acceptance Criteria

- [ ] All acceptance criteria met
