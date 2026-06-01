---
id: "T-quick-doc-drift"
title: "Detect Spec Drift"
priority: "P2"
estimated_time: "15min"
dependencies: ["5"]
type: "doc.drift"
surface-key: ""
surface-type: ""
---

Detect spec drift between existing project specs and current code for the sitemap-surface-guard feature.

## Discovery Strategy
1. Run `git diff --name-only main...HEAD` to identify files changed by this feature
2. List all spec files in docs/business-rules/ and docs/conventions/
3. For each spec file, read its `domains` frontmatter
4. Only verify specs whose domains overlap with the changed files
5. Skip specs with no overlap — they are unaffected by this feature

Do NOT scan all spec files blindly. Use git diff to narrow scope first.
If git diff returns no changes, skip — nothing to drift against.

Auto-fix drifted specs and commit with [auto-specs] tag.

## Acceptance Criteria
- [ ] `git diff --name-only main...HEAD` 已执行，变更文件清单已列出
- [ ] 仅检查 domains 与变更文件有重叠的 spec 文件，无重叠的跳过
- [ ] 若发现 drift 已自动修复并提交；无 drift 则记录"no drift detected"
