---
id: "9"
title: "Fix small documentation completeness issues (M-3, M-5, M-6, M-7)"
priority: "P2"
estimated_time: "30m"
dependencies: [5]
type: "doc"
mainSession: false
---

# 9: Fix small documentation completeness issues (M-3, M-5, M-6, M-7)

## Description

4 处小型文档完整性问题：(M-3) breakdown-tasks intent 读取路径不明确，只说 "If proposal.md has intent" 未给出完整路径；(M-5) test-isolation.md 跨 skill 依赖未声明，run-tests SKILL.md 未引用但 gen-test-scripts 通过 INLINE 引用；(M-6) brainstorm proposal 模板 {{AUTHOR}} 占位符无赋值逻辑；(M-7) write-prd manifest 模板使用 {{FEATURE_SLUG}} 与 quick-tasks 的 {{SLUG}} 不一致。

## Reference Files
- `docs/proposals/forge-skill-audit/proposal.md` — M-3, M-5, M-6, M-7 sections, Proposed Solution, Success Criteria
- `plugins/forge/skills/breakdown-tasks/SKILL.md`: Add complete path for intent read (ref: M-3: breakdown-tasks intent 读取路径不明确)
- `plugins/forge/skills/run-tests/rules/test-isolation.md`: Add OWNER/CONSUMERS comment (ref: M-5: run-tests test-isolation.md 跨 skill 依赖未声明)
- `plugins/forge/skills/brainstorm/SKILL.md`: Add {{AUTHOR}} assignment guidance (ref: M-6: proposal 模板 {{AUTHOR}} 占位符无赋值逻辑)
- `plugins/forge/skills/write-prd/templates/manifest.md`: Rename {{FEATURE_SLUG}} to {{SLUG}} (ref: M-7: manifest slug 占位符命名不一致)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Add complete path `docs/proposals/<slug>/proposal.md` for intent read |
| `plugins/forge/skills/run-tests/rules/test-isolation.md` | Add `<!-- OWNER: run-tests \| CONSUMERS: gen-test-scripts (INLINE) -->` header comment |
| `plugins/forge/skills/brainstorm/SKILL.md` | Add `Set {{AUTHOR}} to git config user.name output, or ask user if not available.` in Step 5 |
| `plugins/forge/skills/write-prd/templates/manifest.md` | Replace `{{FEATURE_SLUG}}` with `{{SLUG}}` |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] breakdown-tasks SKILL.md intent 读取部分指定完整路径 `docs/proposals/<slug>/proposal.md`
- [ ] test-isolation.md 头部有 `<!-- OWNER: run-tests | CONSUMERS: gen-test-scripts (INLINE) -->` 注释
- [ ] brainstorm SKILL.md Step 5 包含 `{{AUTHOR}}` 赋值指导（git config user.name 或询问用户）
- [ ] write-prd/templates/manifest.md 使用 `{{SLUG}}` 而非 `{{FEATURE_SLUG}}`

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`
- Only modify markdown files; no Go code changes

## Implementation Notes
- M-7 修改后需验证 `quick-tasks/templates/manifest-quick.md` 确认两者一致（quick-tasks 已使用 {{SLUG}}）
- 回归验证：`grep -r "FEATURE_SLUG" plugins/forge/skills/write-prd/templates/` 确认无残留
