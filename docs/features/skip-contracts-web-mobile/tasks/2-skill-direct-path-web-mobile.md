---
id: "2"
title: "Add direct test generation path for web/mobile journeys in gen-test-scripts"
priority: "P0"
estimated_time: "1-2h"
complexity: "medium"
dependencies: [1]
surface-key: ""
surface-type: "cli"
breaking: false
type: "doc"
mainSession: false
---

# 2: Add direct test generation path for web/mobile journeys in gen-test-scripts

## Description

gen-test-scripts SKILL.md 当前假设所有 journey 都有 Contract 前置产物。对于纯 web/mobile journey，没有 contract 文件，导致静默跳过。需要修改 SKILL.md 路由逻辑：当 journey 的 surface_types 仅含 web/mobile 且无对应 contract 文件时，跳过 contract 前置检查，直接从 journey.md + types/web.md 生成测试脚本（直达路径）。同时在 types/web.md 和 types/mobile.md 补充直达生成规则，并添加覆盖率自检机制。

## Reference Files
- `docs/proposals/skip-contracts-web-mobile/proposal.md` — Proposed Solution (Skill层), Scope > In Scope (gen-test-scripts路由, 覆盖率自检, types规则), Success Criteria SC-5/6/8
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — 需修改 Step 2 路由，新增直达路径分支 (ref: Step 2: Read Contract Specifications, Step 2.5: Load Type Rules)
- `plugins/forge/skills/gen-test-scripts/types/web.md` — 需补充直达映射规则 (ref: Golden Rules, Classification Indicators)
- `plugins/forge/skills/gen-test-scripts/types/mobile.md` — 同 web.md 补充直达映射规则 (ref: Golden Rules, Classification Indicators)
- `docs/conventions/forge-distribution.md` — 分发模型规范，修改 plugin 文件前必读

## Acceptance Criteria
- [ ] SC-5: 直达路径生成的脚本包含与 journey 步骤对应的用户动作调用（click/type/navigate）和至少一个可视化断言（非空、非骨架）
- [ ] SC-6: 按 surface type 自检覆盖率（count journeys_of_type == count test-scripts_of_type），缺口或类型不匹配时 FAIL 并输出缺口列表
- [ ] SC-8: types/web.md 和 types/mobile.md 直达规则能产出包含有意义断言的测试脚本

## Implementation Notes
- SKILL.md Step 2 新增路由：检查 journey.md 的 surface_types 字段。仅含 web/mobile 且无对应 contract → 走直达路径
- 直达映射规则：步骤描述 → step-action（click/type/navigate），前置条件描述 → fixture_spec，预期结果描述 → Outcome（assertVisible/assertText）
- 覆盖率自检：Surface → Test Type 映射定义在 SKILL.md（web→Web E2E Test, mobile→Mobile E2E Test, api→API Functional Test, cli→CLI Functional Test, tui→Terminal Functional Test）
- 直达路径标识输出：`Generating test scripts for <journey> via direct path (surface: web, no contract required)`
- 失败回退：输出 `Direct path generation failed for <journey>: <reason>`，而非静默跳过
- types/web.md 新增 "Direct Path Generation Rules" 段落，定义 journey 步骤到测试动作的映射模板
- types/mobile.md 同样新增直达生成规则段落
- 修改 plugin 文件需遵循 forge-distribution.md 规范
