---
iteration: 1
title: "CTO Adversarial Scoring — Independent Fresh Evaluation"
date: "2026-05-30"
scorer: "adversarial-cto"
rubric: "proposal.md (1000 pts)"
target: 900
document: "proposal.md"
---

# Eval Report — Iteration 1 (Fresh Independent Evaluation)

## Score
SCORE: 690/1000

## Dimension Breakdown
| Dimension | Score | Max | Notes |
|-----------|-------|-----|-------|
| Problem Definition | 90 | 110 | Evidence strong but some claims imprecise; urgency well-argued |
| Solution Clarity | 95 | 120 | Approach concrete but user-facing behavior underdescribed |
| Industry Benchmarking | 70 | 120 | Weak industry references; selected approach justification thin |
| Requirements Completeness | 85 | 110 | Good scenario coverage; NFRs present but some gaps |
| Solution Creativity | 40 | 100 | Self-described as "standard practice"; minimal innovation |
| Feasibility | 85 | 100 | Technically sound; timeline reasonable but optimistic |
| Scope Definition | 65 | 80 | In-scope items concrete but some vague; out-of-scope good |
| Risk Assessment | 75 | 90 | Risks identified; mitigations mostly actionable but incomplete |
| Success Criteria | 65 | 80 | SC testable but coverage gaps and internal inconsistencies |
| Logical Consistency | 70 | 90 | Multiple cross-section contradictions found |

## Attack Points

1. **[Logical Consistency]** Cost estimates contradict across sections. Urgency section states "错过后成本上升 9-17 天" while the Do Nothing comparison table verdict says "成本上升约 10-19 天兼容层维护". These are different ranges (9-17 vs 10-19) for the same metric with no explanation. Must reconcile to a single justified range.

2. **[Logical Consistency]** Package count contradicts. Evidence section states "pkg/ 有 17 个包" but Assumptions Challenged table says "pkg/ 19 个包的粒度总体合理". The actual count verified is 17. The assumption flip entry references a phantom number (19) that does not exist, undermining the credibility of the assumptions analysis.

3. **[Logical Consistency]** Comparison table Do Nothing verdict contains duplicated text: "且用户明确要求实际清理（约 10-19 天额外开销），且用户明确要求实际清理而非仅文档输出" — the clause "且用户明确要求实际清理" is repeated verbatim in the same cell. This is a proofreading error in a proposal asking for 6-10 days of execution time.

4. **[Problem Definition — Evidence]** The evidence claims `"tests/results/raw-output.txt"` appears 2 times in `quality_gate.go` production code, but verification shows 6 occurrences of `tests/results/` paths in that file (including variants like `tests/results/unit-raw-output.txt`). The claim undercounts by using exact-string matching while the actual scope is broader. The SCs (SC-1) correctly use a broader grep pattern `tests/results/`, revealing the evidence section is narrower than the actual work scope.

5. **[Problem Definition — Evidence]** The evidence classifies `getTaskPhase` as a "test-bridge 别名函数" and notes "其中 getTaskPhase 在 validate_index.go 生产代码中亦有 5 处调用，非纯粹死代码". Verification confirms 5 production calls in `validate_index.go` via the `var getTaskPhase = task.GetTaskPhase` indirection. However, the proposal later (scope item 11) lumps this into "test-bridge 别名函数" cleanup alongside purely dead aliases. Deleting this indirection variable requires replacing 5 production call sites — this is not "test-bridge cleanup" but a production refactor disguised as housekeeping. The proposal underplays the blast radius.

6. **[Success Criteria — Coverage]** SC-5 exempts `root.go`, `output.go`, and `surfaces.go` from "zero top-level command files" but `surfaces_detect.go` (182 lines of surface detection logic) is not exempted. Since `surfaces_detect.go` is a companion to `surfaces.go` (shared surface detection, not a command implementation), it should be either exempted or explicitly addressed. SC-5 is ambiguous about its classification.

7. **[Success Criteria — Coverage]** No SC covers the "PR review checklist" deliverable (scope item 14). The proposal commits to "在 docs/conventions/package-organization.md 中附加 PR review checklist 条目" but there is no testable criterion confirming its existence, content quality, or adoption.

8. **[Success Criteria — Coverage]** Scope item 10a promises file splitting for `internal/cmd/` files over 500 lines during Phase 2c, and SC-10 covers this. But the same item mentions `pkg/` layer large files (5 files: `forgeconfig/config.go` 1272 lines, `task/pipeline.go` 1097 lines, `forgeconfig/detect_surface.go` 962 lines, `task/build.go` 638 lines, `task/autogen.go` 518 lines) as "后续迭代目标" without any SC tracking this deferral. There is no success criterion confirming the decision to defer was intentional and documented, leaving it as an untracked commitment.

9. **[Success Criteria — Consistency]** SC-12f is a fallback condition that modifies SC-5 and SC-9 with the phrase "按实际可达目标调整" without specifying concrete fallback values. If cross-module dependencies exist, what does SC-5 become? "No more than N top-level files"? The fallback leaves SC-5 undefined in the failure path, making the SC set partially unsatisfiable-by-design.

10. **[Solution Clarity — User-facing behavior]** The proposal does not describe the developer experience after reorganization. When a developer wants to add a new CLI command post-Phase-2c, what is the step-by-step workflow? Which directory? How do they register the command? The package organization rules are described abstractly ("cmd -> internal -> pkg" dependency direction, three-layer model) but the concrete developer workflow is missing — a gap for a proposal whose primary output is conventions.

11. **[Industry Benchmarking]** The selected approach cites `golangci-lint` and `helm` as references but both are explicitly disclaimed: "此为作者对公开仓库结构的解读，非官方声明". This means the industry benchmarking section has zero verified industry references. Even `golang-standards/project-layout` is noted as "NOT an official Go standard". The entire benchmarking rests on the author's interpretation without verifiable external validation.

12. **[Industry Benchmarking]** Only 3 alternatives are listed (Do Nothing, Docs Only, Lint-driven), plus the selected approach. The rubric requires "at least 3 meaningful alternatives". "Do Nothing" and "Docs Only" are minimal-effort baselines. The third (Lint-driven) is the only genuinely different approach. Missing: (a) automated refactoring via `gopls` workspace actions as primary driver, (b) adopting an existing Go project layout template wholesale, (c) incremental per-package migration with feature flags, (d) Big Bang rewrite of targeted packages.

13. **[Requirements Completeness — Constraints]** The proposal claims Go 1.25 as a constraint for `0o644` octal literals but does not verify whether all CI environments, developer tooling, and IDE integrations support Go 1.25 syntax. The go.mod confirms `go 1.25`, but the constraint section does not address toolchain readiness beyond the module declaration.

14. **[Risk Assessment]** No risk covers the scenario where Phase 1's dependency graph analysis (`go list -json ./pkg/...`) reveals unexpected cross-domain dependencies that invalidate the proposed three-layer model (types/infrastructure/domain). The proposal assumes the layer model can be retrofitted, but if `pkg/infocmd` (currently imported by 4 domain packages: `research`, `proposal`, `task`, `lesson`) itself depends on another domain package, the proposed merge to `pkg/util/` may create circular dependencies.

15. **[Scope Definition — In-scope vagueness]** Scope item 8 contains a catch-all: "待 Phase 1 偏差分析裁决（~2 个）" without naming which 2 packages. The reader cannot determine the full scope of work without first completing Phase 1. This is scope-as-a-function-of-future-discovery — acceptable for exploration but problematic for effort estimation. The 6-10 day timeline cannot be validated if 2 unknown packages may or may not require reorganization.

16. **[Solution Creativity]** The proposal explicitly states "此方案并非创新，而是工程实践的标准操作". The three-layer model is standard DDD layering. The blast-radius ordering is standard risk management. There is no cross-domain inspiration (no borrowing from JS/Rust/Java ecosystems), no novel tooling integration, no creative insight. The honesty is commendable but the creativity score reflects the self-assessment.

17. **[Problem Definition — Urgency]** The urgency argument states "v3.0.0 是成本最低的重构窗口（非唯一，但错过后成本上升 9-17 天）". The parenthetical "非唯一" contradicts the urgency framing — if it is not the only window, the urgency is reduced. The cost estimate of 9-17 days (or 10-19 days depending on section) assumes all 17 packages would need simultaneous post-release moves, which is a worst-case assumption not supported by evidence.

18. **[Feasibility — Timeline]** Phase 2c at "2-3 days" covers reorganizing `internal/cmd/` (15 top-level files into sub-packages) AND the `pkg/` layer (up to 6 package merges). For `pkg/infocmd` alone (imported by 4 packages: `research`, `proposal`, `task`, `lesson`), import path changes cascade through all dependents. The total scope — 15 cmd file moves + 6 pkg merges with cascading import updates + test updates + compilation verification — in 2-3 days is optimistic.

19. **[Requirements Completeness — Error scenarios]** No error scenario covers `goconst` violations in test files. The proposal evidence shows 17 test-file occurrences of magic paths. SC-1 only greps `forge-cli/internal/` and `forge-cli/pkg/`, excluding `forge-cli/tests/` and `*_test.go` files. Test-file magic values may or may not be in scope — the proposal is silent, and SC-1 as written would pass even if test files still contain hardcoded paths.

20. **[Scope Definition]** Scope item 15 (`make check-cross-module-deps`) is scoped as Phase 2a but exists to "为 Phase 2c 的'不保留兼容层'决策提供持续保护". If Phase 2c is descoped per SC-12f fallback, this CI check becomes orphaned infrastructure. The proposal does not address whether it should be removed or retained in the fallback scenario.

## Reasoning Audit

### Problem -> Solution Chain
The argument: "codebase lacks coding standards" -> "establish standards + reorganize code in 4 phases". The link is sound — standards without cleanup is explicitly rejected as an alternative. However, the proposal conflates two distinct problems: (1) missing documented standards and (2) existing code quality issues (magic values, dead code, package disorganization). Dead code and magic values exist even in projects WITH conventions; they are maintenance failures, not convention failures. The solution addresses both simultaneously, which is appropriate for the v3.0.0 window but conflates causation.

### Solution -> Evidence Chain
Evidence is well-grounded in verifiable codebase facts. I verified most claims against the actual code and found them substantially accurate, with exceptions noted in attacks #4 and #5. The evidence section is one of the strongest parts of this proposal.

### Evidence -> Success Criteria Chain
Most SCs trace directly to evidence items (SC-1 through SC-4 -> magic value examples, SC-6 -> Debugf, SC-7 -> dead code). However, SC-5 (command subpackaging) and SC-9 (package count) have no direct evidence antecedents — they are solution-driven rather than evidence-driven.

### Self-Contradiction Check
Found 3 distinct contradictions (attacks #1, #2, #3). The `consistency_check_result` at the bottom of the proposal acknowledges "2 conflicts_found" but this underestimates the actual count. The automated checker missed the cost-estimate discrepancy and the package-count discrepancy.

### SC Consistency Deep-Dive
Clustered by affected area:

- **Package structure group**: SC-5, SC-9, SC-12f — SC-5 and SC-9 have fallback modifications via SC-12f but fallback values are undefined for SC-5.
- **Magic values group**: SC-1 through SC-4, SC-8 (goconst) — Satisfiable as a set, but SC-1 scope excludes test files creating a coverage gap.
- **Dead code group**: SC-6, SC-7 — Satisfiable in isolation. However, the `Scope` field deletion claim in SC-7 is problematic: `Scope` is actively used by `CheckLegacyScope` in `query.go`, `list.go`, `add.go`, and `build.go` for legacy task migration detection. Deleting `Scope` without replacing the migration check would break those codepaths. The proposal does not address this dependency.
- **Documentation group**: SC-8 — "6 个与 forge-cli 相关的规范文件" (扩展 2 + 新增 4) — countable and satisfiable.
- **Build stability group**: SC-11 — Standard and testable.

## Blindspots

1. **Scope field deletion creates migration regression**: The proposal marks `Scope` as dead code for deletion (SC-7, scope item 10). But `CheckLegacyScope` in `pkg/task/migrate.go` and 4 production call sites in `query.go`, `list.go`, `add.go`, `build.go` actively read `Scope` to detect and block legacy tasks. Deleting `Scope` from `FrontmatterData` would silently break legacy migration detection. The proposal must either: (a) keep `Scope` until all legacy tasks are migrated and remove `CheckLegacyScope` first, or (b) design an alternative detection mechanism.

2. **pkg/util/ as a dumping ground**: The proposal creates `pkg/util/` and proposes merging ~3 small packages into it. The Go community (including standard library philosophy and `golang-standards/project-layout`) strongly discourages `util` packages — packages should be named by what they provide. The Phase 1 dependency graph might reveal these packages have distinct responsibilities better preserved under more specific names.

3. **No rollback plan for Phase 1 review gate**: The proposal says "Phase 1 产出须由项目维护者 review 后方可进入 Phase 2" but does not define review criteria, reviewer identity, SLA, or what happens if review takes 2 weeks. The Phase 2-4 timeline assumes immediate Phase 1 approval.

4. **Missing definition of "command implementation file"**: SC-5 says "internal/cmd/ 下零个顶层命令实现文件" with 3 exemptions. But the classification criteria are unstated. Is `verify_task_done.go` a command implementation or shared validation? Is `cleanup.go` a command or lifecycle utility? Without explicit criteria, SC-5 is subjective.

5. **Scope item numbering error**: The proposal has two items numbered "10" (scope items 10 and 10a). While 10a is a sub-item, the numbering is inconsistent with other hierarchical structures (Phase 2a/2b/2c).

6. **No compile-time impact NFR**: Package reorganization affects Go parallel compilation. Merging packages reduces granularity and may increase build times. For a proposal touching 17 packages, the absence of any build-time impact assessment is a minor gap.
