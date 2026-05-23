# Eval Report: Forge CLI Clean Code Proposal

**Iteration**: 1
**Evaluator**: CTO Expert (Adversarial)
**Date**: 2026-05-24
**Score**: 575/1000

---

## Phase 1: Reasoning Audit

### Argument Chain Trace

1. **Problem -> Solution**: The problem (technical debt in 92-file Go codebase) maps well to a phased refactoring. However, the solution contains factual errors in its evidence base (see Debugf, SetFeature, frontmatter counts) that undermine the audit's reliability.

2. **Solution -> Evidence**: Multiple evidence claims are factually incorrect or misleading:
   - `Debugf` in `cmd/output.go` is claimed "unused" but is used 10+ times in `quality_gate.go`
   - `SetFeature()` is classified as dead code but has 7+ active call sites
   - Frontmatter parsing claimed "3 files" but actual codebase shows 4+ implementations with different signatures
   - `mapXxxToSlugLens` claimed as "4 functions" but only 3 exist (1 variant)

3. **Evidence -> Success Criteria**: Success criteria are mostly testable but "0 处非顶层函数中的 os.Exit 调用" is ambiguous because "top-level" is undefined. The `quality_gate.go` has two `os.Exit(0)` calls in a `RunE` handler that may or may not qualify.

4. **Self-contradiction check**: Phase ordering contradicts the "bottom-up" claim (Phase 2 and Phase 3 should be reversed). The scope claims "keep exports unchanged" for testbridge but also says "clean up the pattern + migrate getTaskPhase" -- these goals conflict.

### Pre-Score Anchors

- Evidence integrity is the core problem: 4+ factual errors in the 18-item audit list
- Phase ordering is inverted between Phase 2 and Phase 3, risking double-churn
- `validateRecordData` refactor scope is severely underestimated (30+ test cases affected)
- The proposal's honesty about being "no innovation" is commendable and should not be penalized

---

## Phase 2: Rubric Scoring

### 1. Problem Definition: 65/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 30/40 | Core problem (tech debt accumulation in forge-cli) is clear and unambiguous. However, the categorization of "dead code" is mixed with "deprecated-but-active" code (SetFeature), blurring the problem boundary. |
| Evidence provided | 15/40 | **Major deduction**: The 18-item audit contains factual errors. Quote: "cmd/output.go 重复定义了 base.Debugf，且该重导出未被使用" -- this is false; `quality_gate.go` uses `cmd.Debugf()` at 10+ call sites. Quote: "YAML frontmatter 解析在 3 个文件中独立实现" -- actual count is 4+ with different signatures. Quote: "4 个 mapXxxToSlugLens 函数" -- only 3 exist. Evidence is the foundation of this proposal; errors here cast doubt on the entire audit. |
| Urgency justified | 20/30 | The v3.0.0-rc.19 pre-release timing is a valid urgency argument. However, no quantification of "post-release refactoring cost increase" is provided -- the claim "重构成本将显著增加" is vague without data. |

### 2. Solution Clarity: 75/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 30/40 | Four phases with specific file names and function names. A reader can explain back what will be done. Deduction: Phase descriptions mix mechanical and design-level tasks without distinguishing them (e.g., frontmatter consolidation is design work, not mechanical dedup). |
| User-facing behavior described | 40/45 | Quote: "所有 CLI 命令的输入输出保持不变" -- clear user-facing contract. Quote: "纯重构：不引入新依赖、不改变外部行为" -- unambiguous. Good. |
| Technical direction clear | 5/35 | **Major deduction**: Phase ordering is wrong. Quote: "按自底向上顺序执行四阶段纯重构：死代码消除 -> 重复逻辑合并 -> 超大文件拆分 -> 反模式修复". But Phase 2 (re-export cleanup) requires changing import paths across packages, while Phase 3 (file splitting) stays within the same package. Executing Phase 2 before Phase 3 means every file touched in Phase 3 must have its imports re-verified, doubling the churn on the same call sites. The technical direction is concrete but contains a sequencing error that undermines the entire "bottom-up" rationale. |

### 3. Industry Benchmarking: 60/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 20/40 | Quote: "Go 社区标准做法：golangci-lint 发现问题 + 手动重构修复". A single reference to general Go community practice. No specific tools, libraries, or documented approaches cited beyond golangci-lint. |
| At least 3 meaningful alternatives | 20/30 | Three alternatives presented (do nothing, lint-only, lint + selective refactoring). The "do nothing" alternative is straw-manned: quote "技术债继续累积，发布后重构更贵" without evidence of accumulation rate. The "lint-only" alternative is also weak -- it is dismissed for not solving structural issues, but no analysis of what fraction of the 18 issues lint tools actually can address. |
| Honest trade-off comparison | 10/25 | The comparison table is perfunctory. Quote: "需阅读 docs/conventions/forge-distribution.md 了解分发约束" is listed as a dependency but no trade-off analysis of how this constraint affects the chosen approach vs. alternatives. No analysis of the cost of not doing this work (e.g., developer velocity impact, bug rate correlation with file size). |
| Chosen approach justified against benchmarks | 10/25 | Quote: "Selected: 平衡效果与风险". The justification is a single phrase. No quantification of expected impact, no measurement framework to compare results against the "do nothing" baseline post-implementation. |

### 4. Requirements Completeness: 55/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 20/40 | Three scenarios listed (developer navigation, dependency check, CI testing). Missing: (1) what happens when a refactor breaks a test -- rollback plan? (2) edge case: what if SetFeature migration misses a call site? (3) error scenario: what if frontmatter consolidation changes behavior for a specific YAML edge case? |
| Non-functional requirements | 20/40 | Quote: "向后兼容：所有 CLI 命令的输入输出保持不变". Good. Quote: "代码量缩减：消除重复后总行数应减少". This is imprecise -- "应减少" is not measurable. How much? 5%? 20%? No target. |
| Constraints & dependencies | 15/30 | Quote: "Go 1.25 工具链" and "需阅读 docs/conventions/forge-distribution.md". The forge-distribution constraint is mentioned but not analyzed -- how does it constrain the file splitting or re-export cleanup? The proposal does not explain. |

### 5. Solution Creativity: 65/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 25/40 | The proposal is explicitly honest: "无创新。这是标准的代码健康维护". This is appropriate for a refactoring proposal. The Assumptions Challenged table (5 Whys, XY Detection, Assumption Flip) is a modest methodological contribution. |
| Cross-domain inspiration | 15/35 | No cross-domain inspiration. The approach is standard Go community practice. Not penalized heavily since the proposal explicitly disclaims innovation. |
| Simplicity of insight | 25/25 | The phased approach (safe deletions first, structural changes last) is a sound and elegant principle, even if the execution order has a flaw. Full marks for simplicity of insight. |

### 6. Feasibility: 60/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 25/40 | Go tooling supports the approach. However, several items are underestimated: (1) validateRecordData refactor requires restructuring 30+ test cases (quote: "修复 validateRecordData() 中的 os.Exit -> 返回 error"), (2) frontmatter consolidation requires API design across 4 different signatures, (3) Debugf deletion would break quality_gate.go. The proposal treats design-level work as mechanical refactoring. |
| Resource & timeline feasibility | 20/30 | Quote: "工作量约 15 个独立任务". No time estimate. No team size specified. "15 independent tasks" is a count, not a timeline. |
| Dependency readiness | 15/30 | Quote: "无外部依赖。所有修改在本地完成。". The forge-distribution.md dependency is mentioned but not analyzed. The dependency on existing test infrastructure (submit_test.go, integration_test.go) for validation is not discussed as a constraint. |

### 7. Scope Definition: 55/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 20/30 | Most items are concrete (specific files, specific functions). However, some are ambiguous: "清理 testbridge.go 模式 + 迁移 getTaskPhase() 到 pkg/task/" -- what does "clean up the pattern" mean exactly? "统一错误处理模式" -- which error handling pattern, unified to what? |
| Out-of-scope explicitly listed | 20/25 | Five items explicitly out of scope. Good. Missing: no mention of whether quality_gate.go's os.Exit(0) calls are in or out of scope. |
| Scope is bounded | 15/25 | The 18-item audit with 4 phases is bounded. However, the scope creep risk is real: Phase 1 includes SetFeature migration (7+ call sites, not dead code), Phase 2 includes frontmatter consolidation (design work, not mechanical), Phase 4 includes validateRecordData refactor (30+ test cases affected). Each of these expands the scope beyond what the phase label implies. |

### 8. Risk Assessment: 55/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 20/30 | Four risks listed. Missing: (1) factual errors in the evidence base leading to destructive deletions (e.g., deleting Debugf that is actually used), (2) Phase ordering risk (double-churn from Phase 2 before Phase 3), (3) test architecture impact from validateRecordData refactor (30+ test cases). These are significant omissions. |
| Likelihood + impact rated | 15/30 | Ratings use M/L/H but are not calibrated. Quote: "拆分文件时遗漏导出符号 -- Likelihood: M, Impact: M". But the re-export cleanup risk (L/H) is rated Low likelihood despite the evidence showing the codebase already has inconsistent import patterns. The Debugf deletion risk (factual error -> broken build) is not listed at all. |
| Mitigations are actionable | 20/30 | Most mitigations are actionable ("每步运行 go build ./... 验证编译"). However, the testbridge mitigation ("保持导出接口不变，仅重新组织") contradicts the scope ("清理 testbridge.go 模式 + 迁移 getTaskPhase()"). If you migrate functions to pkg/task/ and keep the bridge as aliases, that is "reorganization". If you remove the bridge, that breaks integration tests. The mitigation does not resolve this contradiction. |

### 9. Success Criteria: 35/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 25/55 | Most criteria are measurable (build passes, test passes, line counts). However: (1) "0 处非顶层函数中的 os.Exit 调用" -- "top-level" is undefined. quality_gate.go has os.Exit(0) in a RunE handler -- is that "top-level"? (2) "总行数减少" -- by how much? 1 line counts. (3) "0 处重复的 YAML frontmatter 解析" -- after consolidation, what counts as "duplicate"? If pkg/task has one and pkg/frontmatter has another, is that still "duplicate"? |
| Coverage is complete | 10/25 | Missing criteria: (1) no criterion for SetFeature migration completion, (2) no criterion for testbridge cleanup, (3) no criterion for askAutoBehavior refactoring, (4) no criterion for defaultRunClaude deduplication. Several in-scope items have no corresponding success criterion. |

### 10. Logical Consistency: 50/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 20/35 | The solution addresses the problem categories (dead code, duplication, large files, anti-patterns). However, the solution contains factual errors in its evidence that could lead to incorrect actions (deleting Debugf that is used, treating SetFeature as dead code). |
| Scope <-> Solution <-> Success Criteria aligned | 10/30 | **Misalignment**: Phase 1 includes "完成 SetFeature() 迁移并删除废弃函数" but no success criterion mentions SetFeature. Phase 4 includes "askAutoBehavior() refactoring" but no success criterion measures it. Phase 4 includes "testbridge.go cleanup + migrate getTaskPhase()" but no criterion covers it. Phase 2 includes "提取 defaultRunClaude() 到共享位置" but no criterion verifies it. Multiple scope items lack corresponding success criteria. |
| Requirements <-> Solution coherent | 20/25 | Requirements (backward compatibility, build stability, code reduction) are coherent with the solution approach. Minor deduction: the "代码量缩减" requirement ("总行数应减少") conflicts with Phase 2's frontmatter consolidation (creating a new package may add boilerplate) and Phase 4's validateRecordData error-return (adding error handling code may increase line count). |

---

## Phase 3: Blindspot Hunt

### Beyond-Rubric Findings

1. **Evidence reliability crisis**: The proposal's audit claims 18 issues but at least 3 have factual errors (Debugf unused, SetFeature dead code, mapXxx count). This raises the question: were the other 15 items verified? The proposal does not include a verification methodology. If any other claims are wrong, the implementation could cause breakage.

2. **No rollback plan**: The proposal has no explicit rollback strategy. Quote: "每个阶段完成后确保所有测试通过". This is a forward-gating strategy, not a rollback plan. If Phase 3 (file splitting) is half-complete and a problem is discovered, how do you roll back? Git revert works for single commits, but the proposal does not specify commit granularity per phase or per task.

3. **Test coverage as invisible dependency**: The proposal says "新增测试 -- 保持现有测试通过即可" is out of scope. But the validateRecordData refactor requires restructuring 30+ test cases. This is effectively "new test code" -- it falls through the cracks between in-scope (refactor production code) and out-of-scope (new tests).

4. **Frontmatter consolidation is feature work disguised as cleanup**: Creating a new `pkg/frontmatter/` package means designing a public API. The existing implementations have different signatures, different return types, and different error handling. Choosing the canonical API is a design decision, not a mechanical refactoring. This belongs in a separate proposal.

5. **Debugf deletion would be catastrophic**: The proposal states Debugf in cmd/output.go is unused. If a developer executes Phase 1 as written and deletes this function, quality_gate.go breaks at 10+ call sites. This is not a theoretical risk -- it is a factual error that would cause immediate build failure.

---

## Phase 4: Summary of Injected Freeform Findings

| Finding | Severity | Disposition |
|---------|----------|-------------|
| Phase 2/3 ordering inversion | high | Incorporated into Solution Clarity (technical direction) and Feasibility deductions |
| Debugf "unused" claim is false | high | Incorporated into Problem Definition (evidence) and Risk Assessment deductions |
| validateRecordData test refactoring scope | high | Incorporated into Feasibility and Scope Definition deductions |
| Contradictory cleanup goals (testbridge) | medium | Incorporated into Logical Consistency deductions |
| Frontmatter consolidation = design work | medium | Incorporated into Feasibility and Scope Definition deductions |
| SetFeature not dead code | medium | Incorporated into Problem Definition (evidence) and Logical Consistency deductions |
| quality_gate.go os.Exit(0) not mentioned | medium | Incorporated into Success Criteria (coverage) and Scope Definition deductions |
| mapXxxToSlugLens low value | low | Noted but minimal score impact; acknowledged in Feasibility |

---

## Final Score Breakdown

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 65 | 110 |
| Solution Clarity | 75 | 120 |
| Industry Benchmarking | 60 | 120 |
| Requirements Completeness | 55 | 110 |
| Solution Creativity | 65 | 100 |
| Feasibility | 60 | 100 |
| Scope Definition | 55 | 80 |
| Risk Assessment | 55 | 90 |
| Success Criteria | 35 | 80 |
| Logical Consistency | 50 | 90 |
| **Total** | **575** | **1000** |

---

## Gate Decision: FAIL (threshold: 750)

### Top 3 Actions Required

1. **Fix all factual errors in evidence**: Re-verify every claim in the 18-item audit. The Debugf "unused" claim is demonstrably false. SetFeature is not dead code. Frontmatter implementations are 4+, not 3, with different signatures. mapXxx count is 3, not 4. Correct these before any implementation begins.

2. **Reverse Phase 2 and Phase 3 ordering**: File splitting (same package, no import changes) should come before re-export cleanup (cross-package import changes). This eliminates double-churn on the same call sites and aligns with the stated "bottom-up" principle.

3. **Add success criteria for every in-scope item**: Currently missing criteria for SetFeature migration, askAutoBehavior refactoring, testbridge cleanup, defaultRunClaude deduplication, and error handling unification. Also define "top-level" for the os.Exit criterion and address quality_gate.go's two os.Exit(0) calls explicitly.
