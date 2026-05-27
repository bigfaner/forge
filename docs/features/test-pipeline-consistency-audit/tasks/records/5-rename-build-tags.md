---
status: "completed"
started: "2026-05-27 20:00"
completed: "2026-05-27 20:00"
time_spent: ""
---

# Task Record: 5 重命名 Go build tag 为 surface-type-specific 标签

## Summary
重命名 Go build tag 从 e2e 为 cli_functional，同步更新 Convention 文件、justfile、init-justfile 模板、test-guide 规则，并删除 deprecated alias

## Changes

### Files Created
无

### Files Modified
- forge-cli/tests/justfile-integration/forge_detection_test.go
- forge-cli/tests/justfile-integration/init_justfile_test.go
- forge-cli/tests/justfile-integration/mixed_cli_test.go
- forge-cli/tests/justfile-integration/execution_test.go
- forge-cli/tests/justfile-integration/main_test.go
- forge-cli/tests/error-handling/error_handling_test.go
- forge-cli/tests/error-handling/main_test.go
- forge-cli/tests/task-lifecycle/task_stage_gates_test.go
- forge-cli/tests/task-lifecycle/submit_test.go
- forge-cli/tests/task-lifecycle/fix_task_claim_priority_test.go
- forge-cli/tests/task-lifecycle/main_test.go
- forge-cli/tests/skill-ops/prompt_test.go
- forge-cli/tests/skill-ops/clean_code_skill_test.go
- forge-cli/tests/skill-ops/plugin_content_test.go
- forge-cli/tests/skill-ops/forensic_test.go
- forge-cli/tests/skill-ops/main_test.go
- forge-cli/tests/testkit/helpers_test.go
- forge-cli/tests/testkit/helpers.go
- forge-cli/tests/scope-resolution/scope_resolution_test.go
- forge-cli/tests/scope-resolution/main_test.go
- forge-cli/tests/forge-commands/forge_info_commands_test.go
- forge-cli/tests/forge-commands/e2e_commands_test.go
- forge-cli/tests/forge-commands/discovery_test.go
- forge-cli/tests/forge-commands/forge_init_install_just_test.go
- forge-cli/tests/forge-commands/main_test.go
- forge-cli/tests/task-type-system/task_types_dispatch_test.go
- forge-cli/tests/task-type-system/task_type_refinement_test.go
- forge-cli/tests/task-type-system/main_test.go
- tests/surface-aware-recipe-generation/smoke_test.go
- tests/surface-aware-recipe-generation/step2_init_justfile_test.go
- tests/surface-aware-recipe-generation/helpers_test.go
- tests/surface-aware-recipe-generation/step1_configure_surfaces_test.go
- tests/surface-aware-recipe-generation/step4_user_customized_test.go
- tests/surface-aware-recipe-generation/main_test.go
- tests/surface-aware-recipe-generation/step3_verify_recipes_test.go
- tests/surface-aware-recipe-generation/step5_mixed_project_test.go
- tests/task-lifecycle/task_lifecycle_test.go
- tests/task-lifecycle/main_test.go
- tests/task-lifecycle/task_record_test.go
- tests/test-suite-health/simplify_e2e_tests_test.go
- tests/test-suite-health/gen_journeys_skill_test.go
- tests/test-suite-health/risk_density_test.go
- tests/test-suite-health/e2e_test_quality_cleanup_test.go
- tests/test-suite-health/main_test.go
- tests/automated-test-orchestration/step3_execute_dev_test.go
- tests/automated-test-orchestration/step1_run_tests_frontmatter_test.go
- tests/automated-test-orchestration/smoke_test.go
- tests/automated-test-orchestration/step4_execute_probe_test.go
- tests/automated-test-orchestration/step5_execute_test_test.go
- tests/automated-test-orchestration/step7_alternative_surfaces_test.go
- tests/automated-test-orchestration/step6_teardown_test.go
- tests/automated-test-orchestration/helpers_test.go
- tests/automated-test-orchestration/main_test.go
- tests/testkit/forge_binary.go
- tests/automated-test-orchestration/step2_load_strategy_rule_test.go
- tests/testkit/helpers.go
- tests/quality-gate/main_test.go
- tests/feature-management/proposal_status_lifecycle_test.go
- tests/quality-gate/quality_gate_test.go
- tests/feature-management/feature_set_test.go
- tests/feature-management/main_test.go
- tests/test-generation/test_guide_test.go
- tests/test-generation/quick_test_slim_test.go
- tests/test-generation/test_scripts_per_type_test.go
- tests/test-generation/gen_test_scripts_test.go
- tests/test-generation/forge_commands_test.go
- tests/test-generation/integration_test.go
- tests/test-generation/main_test.go
- tests/surface-key-migration/step6_fix_task_test.go
- tests/surface-key-migration/step2_task_struct_test.go
- tests/surface-key-migration/step5_task_add_test.go
- tests/surface-key-migration/smoke_test.go
- tests/surface-key-migration/step3_resolve_scope_test.go
- tests/surface-key-migration/helpers_test.go
- tests/surface-key-migration/step4_breakdown_tasks_test.go
- tests/surface-key-migration/step1_surfaces_cli_test.go
- tests/surface-key-migration/step7_zero_regression_test.go
- tests/surface-key-migration/main_test.go
- tests/task-type-system/task_type_refinement_test.go
- tests/task-type-system/main_test.go
- justfile
- docs/conventions/testing/go.md
- docs/conventions/testing/ginkgo.md
- docs/conventions/testing/index.md
- plugins/forge/skills/init-justfile/templates/generic.just
- plugins/forge/skills/init-justfile/templates/go.just
- plugins/forge/skills/init-justfile/templates/mixed.just
- plugins/forge/skills/init-justfile/templates/node.just
- plugins/forge/skills/init-justfile/templates/python.just
- plugins/forge/skills/init-justfile/templates/rust.just
- plugins/forge/skills/init-justfile/rules/surfaces/cli.md
- plugins/forge/skills/init-justfile/rules/surfaces/api.md
- plugins/forge/skills/init-justfile/rules/surfaces/tui.md
- plugins/forge/skills/init-justfile/rules/surfaces/web.md
- plugins/forge/skills/init-justfile/rules/surfaces/mobile.md
- plugins/forge/skills/test-guide/rules/convention-structure.md

### Key Decisions
- Build tag 使用 cli_functional（下划线）而非 cli-functional（连字符），因为 Go build constraint 语法不支持连字符
- alias test-e2e 直接删除而非保留过渡期（v3.0.0 大版本允许破坏性变更）
- Convention 文件（go.md/ginkgo.md/index.md）中 -tags=e2e 更新为 -tags=cli_functional

## Test Results
- **Tests Executed**: Yes
- **Passed**: 29
- **Failed**: 0
- **Coverage**: 84.3%

## Acceptance Criteria
- [x] 所有 //go:build e2e 替换为 //go:build cli_functional
- [x] Convention 文件中 tags=e2e 替换为 surface-specific 值
- [x] init-justfile/templates/ 中 build tag 引用更新
- [x] test-guide/rules/ 中 build tag 表格更新
- [x] 所有 deprecated alias 已删除
- [x] grep -rn "//go:build e2e" tests/ forge-cli/ 返回 0
- [x] grep -rn '\-tags=e2e' justfile plugins/forge/ 返回 0
- [x] grep -rn "alias test-e2e" plugins/forge/ 返回 0
- [x] go build ./... 通过

## Notes
Go build tag 语法不允许连字符，因此使用 cli_functional（下划线）替代任务文件中示例的 cli-functional（连字符）。所有测试（29 个包）全部通过，平均覆盖率 84.3%。
