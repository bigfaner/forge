---
date: "2026-05-30"
scorer: "adversarial-cto"
rubric: "proposal.md (1000 pts)"
target: 900
document: "proposal.md"
baseline_snapshot: "baseline-snapshot/proposal.md"
freeform_review: "freeform-review.md"
---

# Baseline Score: Forge CLI 代码库重组与规范建立

**Total: 647/1000**

---

## Dimension Scores

### 1. Problem Definition: 78/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 30/40 | Core problem is identifiable -- conventions gap leads to code quality issues. However, the problem conflates two distinct concerns: "lack of coding standards" (a documentation problem) and "accumulated technical debt" (a code problem). These have different root causes and different solutions. The title says "代码库重组与规范建立" which combines both, but the problem statement doesn't justify why they must be solved together rather than sequentially. Two readers could reasonably disagree on whether this is primarily a documentation initiative or a refactoring initiative. |
| Evidence provided | 28/40 | Evidence is specific and verifiable in principle, but contains a factual error: quote -- `"tests/results/raw-output.txt"` 在 `quality_gate.go` 中出现 7 次. The freeform reviewer verified only 2 occurrences in `quality_gate.go` itself; the remaining hits are in `quality_gate_test.go`. This is not a rounding error -- it misattributes test-file occurrences to production code, inflating the perceived severity in production. The dead code and package structure evidence is solid and concrete. |
| Urgency justified | 20/30 | Quote -- `v3.0.0 是唯一的大版本重构窗口。发布后 API 和包结构将趋于稳定，技术债修复成本指数增长。` The urgency claim rests on "指数增长" which is vague language without quantification. What does "指数增长" mean concretely? A 2x cost increase? 10x? The claim that this is the "last opportunity" is stated as fact without evidence -- why would minor restructuring not be possible post-release? The v3.0.0 window argument is reasonable but the urgency is overstated. |

### 2. Solution Clarity: 72/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 25/40 | The two-phase structure is clear at the top level, but Phase 2 is dangerously vague. Quote -- `Phase 2 — 代码重组与清理：以新规范为指导，全面重新设计 internal/cmd/ 和 pkg/ 两层包结构，同时彻底清除所有已识别的死代码和魔法值。不保留兼容层。` "全面重新设计" is not concrete. What is the target package structure? Which packages merge where? There is no mapping table from current state to target state. A reader cannot explain back what will actually be built -- only the general direction. |
| User-facing behavior described | 35/45 | The developer-facing behavior is described through 4 key scenarios (write new commands, code review, package restructuring, magic value cleanup). These are adequate but miss the critical developer experience of the restructuring process itself -- what does the developer experience during the transition? How do they handle in-flight PRs? |
| Technical direction clear | 12/35 | Quote -- `Go 的包重组主要是文件移动和 import 路径更新，工具链（gorename、IDE refactor）支持良好。` `gorename` is unmaintained since 2018 and does not support Go modules. Citing it as a technical enabler undermines credibility. The technical direction for HOW the restructuring will happen (which files move where, what the target dependency graph looks like) is absent. The only technical hint is "领域合并" without specifying which domains merge into which. |

### 3. Industry Benchmarking: 55/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 20/40 | Three references: `golang-standards/project-layout`, Go standard library's domain-merge strategy, and `goconst` linter. These are real and relevant. However, the references are shallow -- `golang-standards/project-layout` is cited without engaging with its well-known criticisms (the repo itself says it's NOT a standard). The proposal doesn't reference how mature Go projects like `golangci-lint`, `helm`, or `terraform` organize their `pkg/` layers. |
| At least 3 meaningful alternatives | 15/30 | Three alternatives are listed, but the "do nothing" alternative is a textbook straw man. Quote -- `Rejected: v3.0.0 是最后窗口`. The rejection is based on the unverified urgency claim, not on the merits of doing nothing vs. doing something. The second alternative "仅输出规范文档" is also rejected in a single phrase: `Rejected: 用户要求实际清理`. This dismisses a legitimate intermediate approach without analysis. Only the selected approach receives serious consideration. |
| Honest trade-off comparison | 10/25 | The cons column is perfunctory. The selected approach's con is simply "工作量较大" -- vague language without quantification. How much larger? 2x? 10x? The "do nothing" alternative lists real cons (technical debt growth) but its pros are "零工作量" which trivializes the alternative. A more honest pro would be "zero risk of regression." |
| Chosen approach justified against benchmarks | 10/25 | The justification is circular: the approach was chosen because it's the proposal itself. Quote -- `规范先行 + 代码重组 | 本方案 | 规范指导实践，审计可追溯 | 工作量较大 | **Selected: 两阶段确保方向正确**`. There is no argument for why this particular sequencing (conventions first, then code) is better than the industry-standard approach of incremental refactoring guided by linters (which would address the same problems without a big-bang restructuring). |

### 4. Requirements Completeness: 65/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 25/40 | Four scenarios cover the happy paths: new commands, code review, restructuring, cleanup. Missing edge cases: (1) What happens when a convention document conflicts with existing code that cannot be changed? (2) What about files that serve dual purposes (e.g., `quality_gate.go` is both a command and a library)? (3) How are in-flight PRs handled during restructuring? Error scenarios are absent -- what if a package merge creates circular dependencies? |
| Non-functional requirements | 22/40 | Three NFRs are listed but have gaps. Quote -- `向后兼容：此为 v3.0.0 内部重构，不影响已发布 API（二进制尚未正式发布）`. This claim is not verified -- the proposal does not check whether other Go modules in the monorepo import `forge-cli/pkg/` via `go.mod` replace directives. The "构建稳定性" NFR says every commit must pass tests, but this is aspirational rather than a mechanism. Missing NFRs: performance impact (does package restructuring affect compile times?), rollback plan (how to revert if Phase 2 goes wrong?), documentation completeness (how to verify conventions are complete?). |
| Constraints & dependencies | 18/30 | Four constraints listed. Missing: (1) test-bridge pattern constraint -- the proposal doesn't acknowledge that alias functions in `claim.go` and `testbridge.go` are test infrastructure, not dead code; (2) the `quality_gate.go` file at 1067 lines is a constraint on "just move it" restructuring; (3) no mention of CI/CD pipeline constraints (do build scripts reference specific paths?). |

### 5. Solution Creativity: 30/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 10/40 | The proposal explicitly disclaims novelty. Quote -- `此方案并非创新，而是工程实践的标准操作——在重构窗口期建立规范并执行。` This is honest but means the novelty score is low by definition. The approach is a direct application of well-known Go community practices. No differentiation from standard linter + convention approaches. |
| Cross-domain inspiration | 10/35 | Two sources of inspiration cited: Go standard library package philosophy and `goconst` linter. Both are from the same domain (Go ecosystem). No cross-domain borrowing -- e.g., how Rust's `cargo clippy` manages lint categories, or how TypeScript projects use `eslint --fix` for automated convention enforcement. |
| Simplicity of insight | 10/25 | The insight is simple but perhaps too simple. "Write conventions, then follow them" is not an insight -- it's the default expectation. The lack of any automation or tooling integration (e.g., `goconst` as a CI gate, custom linter rules) means the "insight" relies entirely on human discipline. |

### 6. Feasibility: 52/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 18/40 | Quote -- `Go 的包重组主要是文件移动和 import 路径更新，工具链（gorename、IDE refactor）支持良好。` The cited tool `gorename` is deprecated and incompatible with Go modules. This is a factual error in the feasibility assessment. More importantly, the proposal ignores the test-bridge problem: deleting alias functions that are referenced by test files will break `go test` before Phase 2 even gets to restructuring. The claim that everything is feasible because "it's just file moves" understates the complexity of maintaining test compilation while restructuring. |
| Resource & timeline feasibility | 20/30 | Single person, 4-7 days total. The timeline is plausible for Phase 1 (writing 4 convention docs + extending 2) but Phase 2 is underestimating. Restructuring `pkg/` (19 packages to <=12) involves dependency analysis, file moves, import updates, test adjustments, and verification. Doing this thoroughly in 3-5 days assumes no complications -- a risky assumption for a restructuring touching 19 packages. |
| Dependency readiness | 14/30 | Quote -- `无外部依赖。所有涉及的包都是 forge-cli 内部包。` This is technically true but ignores internal dependencies: (1) test-bridge files create dependencies between test code and alias functions; (2) the proposal doesn't check for cross-module imports within the monorepo; (3) `golangci-lint` configuration may reference specific package paths. |

### 7. Scope Definition: 55/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 22/30 | 11 items listed, most are concrete deliverables. However, items 7-8 are vague: `重组 internal/cmd/ 包结构` and `重组 pkg/ 层` don't specify the target state. Without a mapping table (current package -> target package), these are areas rather than deliverables. Items 1-6 (convention docs) are concrete and well-defined. |
| Out-of-scope explicitly listed | 18/25 | Seven out-of-scope items listed. Notable omission: file-level refactoring (splitting large files like `quality_gate.go` at 1067 lines). The proposal moves files but doesn't commit to splitting oversized ones. Also missing: `go.mod` / `go.sum` changes as a scope item. |
| Scope is bounded | 15/25 | The scope is bounded by "v3.0.0 pre-release" but the actual execution boundary is unclear. Phase 2 says "全面重新设计" which could expand indefinitely if the target structure is disputed. No explicit timebox or feature freeze for Phase 2. |

### 8. Risk Assessment: 50/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 16/30 | Four risks listed. Missing risks: (1) test-bridge deletion breaking test compilation; (2) convention docs being ignored or becoming stale; (3) monorepo cross-module import breakage; (4) the Phase 1/Phase 2 sequencing risk -- conventions "extracted from current patterns" may not provide actionable guidance for "target state restructuring." |
| Likelihood + impact rated | 16/30 | Ratings are honest -- not all high-impact/high-likelihood. However, Risk 2 ("规范过于理想化") has likelihood M and impact L, which seems understated. If conventions are wrong, Phase 2 executes against wrong guidance, which is high impact. Risk 4 ("破坏 golangci-lint 配置") at L/L seems reasonable. |
| Mitigations are actionable | 18/30 | Mitigations range from actionable to vague. Risk 1 mitigation ("每步重组后立即 go build + go test 验证") is actionable. Risk 2 mitigation ("规范基于现有代码模式提炼，而非凭空设计") is not a mitigation -- it's a design choice. A mitigation would be "Phase 1 output includes a deviation analysis showing which existing code violates the conventions, with specific remediation steps." Risk 3 mitigation ("每个包内通过文件名和注释区分子职责") is a design guideline, not a mitigation for the risk of merged packages becoming incoherent. |

### 9. Success Criteria: 48/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 22/30 | SC-1 through SC-4 are excellent -- they use concrete grep commands with expected zero results. SC-5 through SC-7 are clear. SC-8 (6 convention files) is countable. SC-9 (pkg count <= 12) is measurable. SC-10 (build + test pass) is standard. Good overall. Deductions: SC-5 is ambiguous about what counts as a "command file" vs. infrastructure file (root.go, output.go, surfaces.go). SC-7 says "所有别名函数" but doesn't define which functions are aliases vs. test bridges. |
| Coverage is complete | 11/25 | Significant gaps: (1) No SC for convention document quality (what makes a convention doc "complete"?). (2) No SC for file size health (quality_gate.go at 1067 lines is not addressed). (3) No SC for golangci-lint passing. (4) No SC for test-bridge cleanup. (5) No SC verifying that the target package structure has no circular dependencies. The 10 SCs cover magic values, dead code, package count, and build health -- but miss structural quality. |
| SC internal consistency | 15/25 | SC-7 ("删除所有别名函数") conflicts with the implicit requirement that `go test` must pass (SC-10). As the freeform review identified, alias functions like `checkExistingTaskState`, `getTaskPhase`, `compareVersionIDs` in `claim.go` are test infrastructure referenced by test files. Deleting them breaks test compilation. SC-5 ("零个顶层命令文件") may conflict with the practical need for `root.go` and shared infrastructure files at the `internal/cmd/` root. The `consistency_check_result: pass` at the bottom of the proposal is not credible given these contradictions. |

### 10. Logical Consistency: 70/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 28/35 | The solution (conventions + restructuring) does address the stated problems (missing conventions, magic values, dead code, unprincipled packages). The mapping is clear. Minor gap: the problem says "无法指导新代码的编写" but the solution focuses heavily on cleaning existing code, with less emphasis on how conventions will be enforced for future code. |
| Scope <-> Solution <-> Success Criteria aligned | 22/30 | In-Scope items 1-6 map to SC-8. In-Scope items 7-8 map to SC-5, SC-9. In-Scope items 9-11 map to SC-1 through SC-4, SC-6, SC-7. However, In-Scope item 10 ("删除所有别名函数") includes test bridges that are not truly dead code, creating a misalignment between scope and solution intent. SC-9 sets a quantitative target (<=12 packages) without a corresponding scope item specifying which packages merge where. |
| Requirements <-> Solution coherent | 20/25 | Key scenarios map to convention docs and restructuring actions. The constraint "`pkg/types/` 作为 leaf package" is respected in the approach. Gap: the NFR "向后兼容" claims no impact on published API, but the scope includes deleting exported functions (alias functions) which are part of the public API surface, even if the binary isn't released. |

---

## Blindspot Hunt

### [blindspot-1] Test-bridge conflation with dead code
Quote -- `删除所有死代码：deprecated Scope 字段、别名函数、兼容层、构建产物（.out 文件）`. The proposal classifies ALL alias functions as dead code. The freeform reviewer identified `internal/cmd/task/testbridge.go` (125 lines) and `claim.go`'s `var` aliases as active test infrastructure. Deleting them will break test compilation. The proposal has no analysis of which aliases are truly dead vs. which are test bridges. This is not a risk -- it is a design error that will block execution of Phase 2.

### [blindspot-2] Phase 1 produces descriptive docs, Phase 2 needs prescriptive targets
Quote -- `分析现有代码库模式，扩展 docs/conventions/ 下的规范文件`. If Phase 1 "analyzes existing patterns" it will document the current messy state. Phase 2 needs a target state that by definition does not exist yet. The proposal has no mechanism for bridging this gap. Phase 1 must produce TARGET STATE specifications with deviation analysis, not descriptions of current patterns. This sequencing flaw could render Phase 1 output useless for Phase 2.

### [blindspot-3] No rollback plan
The proposal has zero mention of what happens if Phase 2 goes wrong. Quote -- `不保留兼容层`. If the restructuring introduces subtle runtime bugs (not just compilation errors), there is no rollback path. For a proposal that touches 19 packages and deletes dead code, the absence of a rollback strategy is a significant blindspot. The mitigation "every commit passes tests" is a prevention mechanism, not a rollback plan.

### [blindspot-4] Factual error in evidence section
Quote -- `"tests/results/raw-output.txt"` 在 `quality_gate.go` 中出现 7 次. Verified as only 2 occurrences in `quality_gate.go`. The remaining occurrences are in `quality_gate_test.go`. This inflates the perceived production-code severity and misleads effort estimation. Test file occurrences should be counted separately.

### [blindspot-5] No automation or enforcement mechanism
The proposal writes convention documents but includes no mechanism to enforce them. Without CI integration (e.g., `goconst` as a lint gate, custom `golangci-lint` rules, pre-commit hooks), conventions are aspirational. The proposal's own evidence shows the codebase already HAS some conventions that are violated. Writing more conventions without enforcement will repeat the same pattern.

### [blindspot-6] Large files moved but not split
`quality_gate.go` (1067 lines), `config.go` (1272 lines), and `pipeline.go` (1097 lines) are identified in the freeform review as exceeding healthy file sizes. The proposal's scope says "重组 internal/cmd/ 包结构：顶层散落的命令文件子包化" which would move `quality_gate.go` to a subdirectory without splitting it. This is structural theater -- relocating problems rather than solving them.
