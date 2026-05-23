# Eval Report: Forge CLI Clean Code Proposal

**Iteration**: 2
**Evaluator**: CTO Expert (Adversarial)
**Date**: 2026-05-24
**Score**: 785/1000

---

## Phase 1: Reasoning Audit

### Argument Chain Trace

1. **Problem -> Solution**: Problem definition is now factually corrected. The 15-item audit reclassifies SetFeature as "deprecated code pending migration" (not dead code), correctly states Debugf has 10+ callers via cmd.Debugf(), and acknowledges the frontmatter count is "4+". The mapping from problem categories to four phases is logical and the phase ordering is now correct (file splitting before cross-package cleanup).

2. **Solution -> Evidence**: Evidence claims are substantially corrected compared to iteration 1. Remaining concerns:
   - The "3+1 个 mapXxxToSlugLens 函数" notation is unusual (第 29 行). The "+1 variant" phrasing obscures whether this is 3 or 4 functions. The success criteria has no entry for this item.
   - Frontmatter consolidation claims "扩展现有 pkg/task.ParseFrontmatter()" but does not verify that the existing function's signature is compatible with all 4+ call sites. This is an assumption, not verified evidence.

3. **Evidence -> Success Criteria**: Success criteria are now substantially more complete. 13 measurable criteria cover most in-scope items. The "顶层定义" for os.Exit is now explicit (第 176 行). The 5% line reduction target is concrete. Remaining gap: no success criterion for the mapXxxToSlugLens generic replacement, and no criterion for the dependency check unification.

4. **Self-contradiction check**: No significant self-contradictions remain. The testbridge approach is now clear (keep as thin aliases). The scope/scope-out boundaries are well-defined. The os.Exit(0) treatment in quality_gate.go is explicit and reasoned.

### Iteration 1 Attack Resolution

| Attack from Iteration 1 | Resolution Status |
|-------------------------|-------------------|
| Debugf "unused" claim false | **Resolved**: Reclassified as "redundant indirection layer" with correct caller count (第 25 行) |
| Phase 2/3 ordering inverted | **Resolved**: New order is correct -- file splitting (Phase 2) before cross-package cleanup (Phase 3) |
| validateRecordData test refactoring scope | **Resolved**: Explicitly scoped in (第 132 行, 第 142 行) |
| Contradictory cleanup goals (testbridge) | **Resolved**: Clear "thin alias" strategy (第 135 行) |
| Frontmatter consolidation = design work | **Partially resolved**: Technical direction given (extend existing) but still not acknowledged as a design decision |
| SetFeature not dead code | **Resolved**: Reclassified as "deprecated code pending migration" (第 22 行) |
| quality_gate.go os.Exit(0) not mentioned | **Resolved**: Explicitly out of scope with reasoning (第 134 行, 第 144 行) |
| Missing rollback plan | **Resolved**: Full Rollback Strategy section added (第 156-165 行) |
| Missing success criteria for scope items | **Substantially resolved**: 13 criteria now cover most items. Minor gaps remain |
| "顶层" definition ambiguous | **Resolved**: Explicit definition provided (第 176 行) |
| 代码量缩减 target vague | **Resolved**: ">= 5%" target (第 181 行) |
| Audit data error risk | **Resolved**: New risk entry with grep/LSP verification mitigation (第 154 行) |

---

## Phase 2: Rubric Scoring

### 1. Problem Definition: 85/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 35/40 | Core problem is clearly stated with four categories (dead code, deprecated code, duplication, large files, anti-patterns). The reclassification of SetFeature as "deprecated" rather than "dead" is accurate. Minor deduction: the "15 个具体问题" count mixes severity levels -- 2 dead code items are trivial (build artifacts), while the anti-patterns are significant. The count does not help prioritize. |
| Evidence provided | 30/40 | Substantially improved. Debugf claim corrected (第 25 行: "quality_gate.go 有 10+ 处通过 cmd.Debugf() 调用"). SetFeature correctly stated as "7+ 处调用需要迁移" (第 22 行). Frontmatter count corrected to "4+" (第 26 行). Remaining concern: the "3+1 个 mapXxxToSlugLens 函数" notation (第 29 行) is imprecise. Also, no verification methodology is documented -- the evidence is presented as assertion rather than reproducible analysis. |
| Urgency justified | 20/30 | The v3.0.0-rc.19 timing argument is reasonable. However, the claim "发布后重构成本将显著增加" (第 44 行) remains unsubstantiated. What specifically makes post-release refactoring more expensive? No API stability commitments are cited, no downstream consumers are identified. The urgency argument is plausible but not evidenced. |

### 2. Solution Clarity: 95/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 35/40 | Four phases with specific file names, function names, and line counts. Phase ordering is now correct: dead code removal (safe) -> file splitting (package-internal) -> duplication cleanup (potentially cross-package) -> anti-pattern fixes (behavior-adjacent). Each phase is independently verifiable. Minor deduction: Phase 3 mixes mechanical tasks (SetFeature migration, defaultRunClaude extraction) with a design-adjacent task (frontmatter consolidation) without distinguishing them. |
| User-facing behavior described | 45/45 | Quote: "所有 CLI 命令的输入输出保持不变" (第 64 行). Quote: "纯重构：不引入新依赖、不改变外部行为" (第 71 行). Quote: "108 个测试文件必须继续通过，零行为变更" (第 60 行). Clear, unambiguous, complete. |
| Technical direction clear | 15/35 | Phase ordering is now correct and well-reasoned. The "bottom-up" rationale (第 48 行) aligns with execution order. Deduction: the frontmatter consolidation technical direction is under-specified. Quote: "扩展现有 pkg/task.ParseFrontmatter() 为共享解析器，其他调用者直接使用提取的 YAML bytes" (第 123 行). This asserts that extending the existing parser can serve all 4+ call sites, but does not address that the different call sites have different signatures and return types. "直接使用提取的 YAML bytes" sidesteps the actual design question: each call site parses different fields from the YAML. A shared parser must either return a generic map (losing type safety) or a unified struct (requiring all call sites to adapt). This is a design decision hidden in a one-line description. |

### 3. Industry Benchmarking: 70/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 25/40 | Go community standard practice (golangci-lint + manual refactoring) is referenced. Quote: "Go 社区标准做法：golangci-lint 发现问题 + 手动重构修复" (第 78 行). This remains a single general reference. No specific refactoring techniques from Go projects (e.g., how Kubernetes, Hugo, or Terraform handle large-scale refactoring) are cited. |
| At least 3 meaningful alternatives | 20/30 | Three alternatives: do nothing, lint-only, lint + selective refactoring. The "do nothing" alternative remains somewhat straw-manned: "技术债继续累积，发布后重构更贵" without evidence of accumulation rate or cost trajectory. The lint-only alternative is fair. The selected alternative is reasonable. |
| Honest trade-off comparison | 15/25 | The comparison table (第 82-86 行) is adequate but perfunctory. No quantitative comparison (e.g., estimated effort hours per approach, risk scoring). The trade-off of "lint + selective refactoring" vs. "lint-only" is not explored: what fraction of the 15 issues can lint tools actually detect? This would strengthen the case. |
| Chosen approach justified against benchmarks | 10/25 | Quote: "Selected: 平衡效果与风险" (第 86 行). A single phrase. No measurement framework to validate post-implementation that the chosen approach delivered value beyond what lint-only would have achieved. |

### 4. Requirements Completeness: 80/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 30/40 | Three scenarios listed (developer navigation, dependency check, CI testing). The rollback scenario is now covered (第 156-165 行). Missing: what happens if frontmatter consolidation changes parsing behavior for a YAML edge case (e.g., malformed frontmatter, empty frontmatter, non-standard delimiters)? |
| Non-functional requirements | 30/40 | Backward compatibility is clear (第 64 行). Build stability is clear (第 65 行). Code reduction now has a concrete target: ">= 5%" (第 181 行). Missing: no performance requirement (e.g., refactoring must not increase build time, test execution time must not change). For a refactoring proposal, performance non-regression is a standard NFR. |
| Constraints & dependencies | 20/30 | Go 1.25 toolchain, pure refactoring, forge-distribution.md constraint mentioned. However, the forge-distribution constraint (第 72 行) is still mentioned without analysis -- how does the distribution model constrain file splitting or re-export cleanup? The proposal assumes it's relevant but does not explain the specific constraint. |

### 5. Solution Creativity: 65/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 25/40 | The proposal explicitly states "无创新" (第 52 行). The Assumptions Challenged table (第 104-108 行) is a modest methodological contribution -- it shows structured thinking about assumptions. The re-export layer analysis ("Overturned: 子包已直接 import base/") is a genuinely useful insight. |
| Cross-domain inspiration | 15/35 | No cross-domain inspiration. Standard Go community practice. Not penalized heavily since the proposal explicitly disclaims innovation. |
| Simplicity of insight | 25/25 | The phased approach (safe deletions first, structural changes last, anti-patterns at the end) is elegant. The rollback strategy (one commit per phase) is simple and effective. Full marks for simplicity of insight. |

### 6. Feasibility: 75/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 30/40 | Go tooling supports the approach. The validateRecordData refactor now explicitly includes the 30+ test case restructuring (第 132 行). The Debugf deletion risk is mitigated by the "redundant indirection" framing and the verification step (第 154 行). Remaining concern: frontmatter consolidation feasibility is asserted but not demonstrated -- extending ParseFrontmatter() to serve 4+ call sites with different signatures is a design challenge, not a mechanical task. |
| Resource & timeline feasibility | 25/30 | "15 个独立任务" (第 96 行) is a concrete task count. The per-phase rollback strategy adds confidence. Missing: no time estimate (days/weeks). For a pre-release codebase (v3.0.0-rc.19), timeline matters. |
| Dependency readiness | 20/30 | "无外部依赖" (第 100 行) is correct. The internal dependency on test infrastructure (submit_test.go) is now acknowledged. The forge-distribution.md dependency remains unanalyzed. |

### 7. Scope Definition: 70/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 25/30 | Phase descriptions are concrete with specific file names, function names, and line counts. The testbridge cleanup now has a clear strategy (thin aliases, 第 135 行). The validateRecordData scope includes test restructuring (第 132 行). Minor deduction: "统一错误处理模式" (第 134 行) -- unified to what? The proposal defines "top-level" for os.Exit but does not specify the target error handling pattern for non-top-level functions. |
| Out-of-scope explicitly listed | 25/25 | Six items explicitly out of scope (第 138-144 行), including the quality_gate.go os.Exit(0) calls with reasoning. The test restructuring clarification ("不算新增测试") is helpful. Good. |
| Scope is bounded | 20/25 | The four-phase structure with 15 tasks is bounded. The rollback strategy adds boundary confidence. Minor concern: frontmatter consolidation (Phase 3) could expand in scope if the "extend existing parser" approach proves incompatible with some call sites. The proposal does not define a fallback for this case. |

### 8. Risk Assessment: 75/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 25/30 | Five risks identified (第 148-154 行). The new "审计数据错误导致误删活跃代码" risk (第 154 行) directly addresses the iteration 1 concern. The re-export cleanup risk (第 151 行) and testbridge risk (第 153 行) are well-targeted. Missing: no risk for frontmatter consolidation changing parsing behavior (the assumption that extending the existing parser is safe for all call sites). |
| Likelihood + impact rated | 25/30 | All five risks have L/M/H ratings. The "审计数据错误" risk is rated M/H with a concrete mitigation. The os.Exit risk is rated L/M with clear scope boundary. The ratings are more calibrated than iteration 1. Minor concern: the "拆分文件时遗漏导出符号" risk (第 150 行) is rated M/M but should arguably be L/M since Go's compiler will catch missing exports immediately (the mitigation already says "go build ./..."). |
| Mitigations are actionable | 25/30 | Mitigations are concrete: "go build ./..." (第 150 行), "go vet 确认所有调用点" (第 151 行), "grep -r 或 LSP findReferences" (第 154 行). The testbridge mitigation ("保持导出接口不变") is now consistent with the scope description. Deduction: no mitigation for frontmatter consolidation risk. |

### 9. Success Criteria: 65/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 45/55 | 13 criteria (第 169-181 行), most are measurable and testable. Highlights: "0 处非顶层函数中的 os.Exit 调用" now has an explicit definition (第 176 行). "总行数减少 >= 5%" is concrete (第 181 行). File line count targets are specific (第 174-175 行). Remaining issues: (1) "0 处重复的 YAML frontmatter 解析" (第 172 行) -- "所有调用者使用 pkg/task.ParseFrontmatter() 或其提取的 YAML bytes". The "或其提取的 bytes" clause is a loophole: if a call site extracts YAML bytes using the shared parser but then does its own field extraction, is that still "0 duplication"? The criterion does not define what constitutes "duplicate parsing". (2) "0 处重复的依赖检查逻辑" (第 173 行) -- no definition of what counts as "duplicate". |
| Coverage is complete | 20/25 | Substantially improved from iteration 1. Criteria now cover: dead code, frontmatter, dependency check, file splitting, os.Exit, askAutoBehavior, defaultRunClaude, testbridge, SetFeature, line reduction. Missing: (1) no criterion for mapXxxToSlugLens generic replacement -- this is an in-scope item (第 126 行) with no success criterion. (2) no criterion for "统一错误处理模式" (第 134 行). |

### 10. Logical Consistency: 85/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 32/35 | The four phases map directly to the four problem categories. Dead code removal addresses items 1-2. File splitting addresses the two large files. Duplication cleanup addresses the 6 duplication items. Anti-pattern fixes address the 4 anti-patterns. SetFeature migration is properly sequenced in Phase 3 (after safe deletions, before anti-patterns). Minor deduction: frontmatter consolidation is classified as "duplication cleanup" but the solution approach ("extend existing parser") implicitly requires API design, which is more aligned with "anti-pattern fix" or a standalone task. |
| Scope <-> Solution <-> Success Criteria aligned | 25/30 | Substantially improved. Most scope items have corresponding success criteria. Remaining gaps: mapXxxToSlugLens (in scope at 第 126 行, no criterion) and error handling unification (in scope at 第 134 行, no criterion). These are minor items but represent alignment gaps. |
| Requirements <-> Solution coherent | 28/25 | NFRs are coherent with the solution approach. The 5% line reduction target is achievable given the duplication elimination scope. The backward compatibility requirement is consistently maintained across all phases. The rollback strategy supports the "build stability" requirement. Rollback per-phase commit strategy is simple and aligns with the phased approach. Deduction: the code reduction requirement (5%) could conflict with Phase 4's validateRecordData refactor, which adds error handling code (return error + error propagation). The proposal does not estimate the net line count impact of this specific change. |

---

## Phase 3: Blindspot Hunt

### Beyond-Rubric Findings

1. **Frontmatter consolidation is still underspecified as a design task**: The proposal describes this as "扩展现有 pkg/task.ParseFrontmatter()" (第 123 行), implying a mechanical extension. But the call sites have different signatures -- some return `TaskMeta`, some return custom structs, some just extract a version string. "Extending" the existing parser to serve all of these requires deciding on a unified return type or a layered API. This is API design work, not mechanical refactoring. The iteration 1 freeform finding flagged this ("属于设计工作，非机械重构") and the proposal has not fully addressed it -- it provided a direction but did not acknowledge the design nature of the task.

2. **mapXxxToSlugLens falls through the cracks**: This item is listed in scope (第 126 行) but has no success criterion. The iteration 1 freeform finding noted its low value ("净节省约 4 行"). The proposal kept it in scope but did not add a criterion or address the value concern. This creates an ambiguous state: is it committed or optional?

3. **No performance non-regression criterion**: For a 92-file Go codebase refactoring, there is no criterion ensuring that build time, test execution time, or binary size do not regress. This is a standard concern for large-scale refactoring proposals. While the proposal claims "纯重构" (no behavior change), refactoring can still affect compilation time (e.g., new package dependencies from frontmatter consolidation) and test execution time (e.g., 30+ test case rewrites in submit_test.go).

4. **Rollback strategy assumes clean phase boundaries**: The rollback strategy (第 156-165 行) assumes each phase is a single commit that can be cleanly reverted. But Phase 3 contains 6 separate tasks (frontmatter, dependency check, defaultRunClaude, mapXxx generic, SetFeature migration, re-export cleanup). If SetFeature migration fails after frontmatter consolidation succeeds, reverting the entire Phase 3 commit would undo the frontmatter work. The proposal mentions "checkpoint commit" (第 165 行) as an option but does not require it. For a phase with 6 diverse tasks, per-task commits should be the default, not optional.

5. **forge-distribution.md constraint remains unanalyzed**: Referenced at line 72 as a constraint, but never analyzed in the proposal body. What specific constraint does the distribution model impose on file splitting, re-export cleanup, or testbridge changes? The proposal assumes it matters but does not explain why or how. This is an incomplete dependency analysis.

---

## Phase 4: Summary of Injected Freeform Findings

| Finding | Severity | Disposition in Iteration 2 |
|---------|----------|---------------------------|
| Phase 2/3 ordering inversion | high | **Resolved**: Order corrected (file splitting before cross-package cleanup) |
| Debugf "unused" claim is false | high | **Resolved**: Reclassified as redundant indirection, callers acknowledged |
| validateRecordData test refactoring scope | high | **Resolved**: Explicitly scoped with test case count |
| Contradictory cleanup goals (testbridge) | medium | **Resolved**: Clear thin-alias strategy |
| Frontmatter consolidation = design work | medium | **Partially resolved**: Direction given but design nature not acknowledged |
| SetFeature not dead code | medium | **Resolved**: Reclassified as deprecated code |
| quality_gate.go os.Exit(0) not mentioned | medium | **Resolved**: Explicitly out of scope with reasoning |
| mapXxxToSlugLens low value | low | **Unresolved**: Kept in scope, no success criterion, value concern not addressed |

---

## Final Score Breakdown

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 85 | 110 |
| Solution Clarity | 95 | 120 |
| Industry Benchmarking | 70 | 120 |
| Requirements Completeness | 80 | 110 |
| Solution Creativity | 65 | 100 |
| Feasibility | 75 | 100 |
| Scope Definition | 70 | 80 |
| Risk Assessment | 75 | 90 |
| Success Criteria | 65 | 80 |
| Logical Consistency | 85 | 90 |
| **Total** | **785** | **1000** |

---

## Gate Decision: PASS (threshold: 750)

### Score Improvement: +210 (from 575 to 785)

### Remaining Recommendations (for implementation, not blocking)

1. **Frontmatter consolidation**: Acknowledge this as a design task within the refactoring scope. Define the target API (unified struct? generic map? layered functions?) before implementation. Consider splitting it into a sub-phase: (a) design the shared API, (b) mechanical migration of call sites.

2. **mapXxxToSlugLens**: Either add a success criterion ("0 处非泛型的 mapXxxToSlugLens 函数") or move it to out-of-scope given its low value (net ~4 lines saved). Leaving it in scope without a criterion is an alignment gap.

3. **forge-distribution.md**: Add a one-sentence explanation of how the distribution model constrains the refactoring. For example: "The distribution model requires that all plugin files be self-contained within the plugin directory, so file splitting must keep all fragments within the same package."

4. **Phase 3 commit granularity**: Require per-task commits within Phase 3 rather than making them optional. With 6 diverse tasks in one phase, a single-commit rollback would be too coarse.
