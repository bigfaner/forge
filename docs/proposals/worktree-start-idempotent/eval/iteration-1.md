# Proposal Evaluation Report: Make `forge worktree start` Idempotent & Rename `copy-files` to `includes`

**Iteration**: 1 of 1
**Date**: 2026-06-09
**Evaluator Role**: CTO / Adversary

---

## Total Score: 826 / 1000

---

## Per-Dimension Scores

### D1. Problem Definition: 98 / 110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 36/40 | Two problems stated, both concrete. Minor ambiguity: the title couples two changes (idempotent start + rename) that are logically independent, but each is individually clear. Deduction: the second problem (naming inconsistency) lacks a clear "what breaks without this" articulation — it's a style issue framed with the same weight as a functional gap. |
| Evidence provided | 32/40 | Four evidence points given. Points 1-3 are specific and verifiable against actual code (confirmed: `cmd_start.go` line 108-110 does print error + hint for existing directory). Point 4 is speculative ("if future support...") — not current evidence, it's a hypothetical. Deduction: -8 for speculative evidence mixed with concrete evidence. |
| Urgency justified | 30/30 | "Cost of delay" articulation is clear and specific: "每次需要在新会话中继续某个 feature 工作时都会遇到" — high frequency pain point with concrete workaround cost (manual cd or resume context noise). |

**Attack Points:**
1. **Evidence item 4 is speculative, not evidence**: "如果未来支持 glob 模式或目录，'copy-files' 不再准确" — this is a forward-looking rationale, not evidence of a current problem. It should be separated from the evidence section.
2. **Two problems in one proposal**: The idempotent start and the config rename are orthogonal changes bundled together. The proposal does not argue why they must be done simultaneously rather than as separate smaller PRs.

---

### D2. Solution Clarity: 110 / 120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | The two-branch behavior (not exists → create + launch; exists → skip + launch) is crisp and unambiguous. Deduction: the "直接替换，不保留任何旧字段兼容逻辑" decision is stated without explaining what happens to users who have existing `copy-files` in their config.yaml — their config simply breaks? |
| User-facing behavior described | 42/45 | Key scenarios cover the major paths. The "entering existing worktree" vs "created new worktree" output distinction is clear. Deduction: the `--interactive` path scenario (#6) says "正常进入" but doesn't specify what the user sees — same output as explicit slug? Different? |
| Technical direction clear | 30/35 | File paths are named (`cmd_start.go`, `forgeconfig/`). "可复用 cmd_resume.go 中的 worktree 验证逻辑" is a clear direction. Deduction: the proposal says "修改「已存在则报错」分支" but the actual code structure (line 107-110) is a single `if` block before any git operations — the modification point is clear, but the proposal does not mention the git worktree validation that `cmd_resume.go` performs (symlink evaluation, .git check). Just saying "可复用" without specifying what exactly is reused is vague. |

**Attack Points:**
3. **No migration path for existing config**: "直接替换，不保留任何旧字段兼容逻辑" — quote from Solution section. Users with `copy-files` in `.forge/config.yaml` will silently lose their config or get parsing errors. No migration guidance.
4. **`--interactive` scenario underspecified**: Scenario #6 says "正常进入，因为 slug 来自交互选择" — but the interactive flow runs `listUnfinishedItems` + `promptSelection`. If the worktree for that slug already exists, what does the user see? Does the interactive prompt filter to only show items with existing worktrees? The proposal doesn't say.

---

### D3. Industry Benchmarking: 92 / 120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 30/40 | `kubectl apply`, `terraform plan`, `mkdir -p` are cited. These are appropriate analogies for idempotent operations. Deduction: the industry references only cover the idempotency aspect. The config renaming (copy-files → includes) has zero industry benchmarking — no examples of similar renames in CLI tools, no citation of how other tools handle config field migrations. |
| At least 3 meaningful alternatives | 22/30 | Four alternatives listed including "do nothing". However: (1) "新增 `open` 子命令" is a straw man — described in just 6 words with no real analysis. (2) "给 `resume` 加 `--fresh` flag" is also minimal. (3) None of the alternatives is an industry-validated solution in the sense of "another tool does X" — they are all self-invented. Deduction: -8 for two weak alternatives. |
| Honest trade-off comparison | 20/25 | The comparison table is present with pros/cons. Deduction: the "Cons" column for the selected approach says "轻微改变 start 的现有语义" — this is understated. Any behavioral change in a CLI command can break scripting integrations, and the Risk section itself rates this as L likelihood / M impact. The trade-off analysis contradicts the risk assessment on severity. |
| Chosen approach justified against benchmarks | 20/25 | "最符合「最小惊讶原则」" is the justification. `mkdir -p` analogy is apt. Deduction: the justification for why `start` is better than `open` is "认知负担大于收益" — but this is asserted, not demonstrated. No user research or cognitive load analysis is provided. |

**Attack Points:**
5. **Straw-man alternatives**: "新增 `open` 子命令" with pros "职责分离清晰" and cons "增加命令数量" — this is a legitimate alternative dismissed in one sentence without exploring whether it avoids the behavioral change risk entirely. It's presented to be rejected.
6. **Config rename has zero benchmarking**: The entire Industry Benchmarking section only discusses idempotency. The copy-files → includes rename (50% of the proposal's scope) has no industry comparison whatsoever.
7. **Trade-off severity mismatch**: Comparison table says "轻微改变 start 的现有语义" but Risk table says L likelihood / M impact for script dependency breakage. These are inconsistent — "轻微" contradicts "M impact".

---

### D4. Requirements Completeness: 92 / 110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 34/40 | Six scenarios listed covering happy paths and edge cases. Deduction: Missing scenario — what happens if the worktree directory exists but is corrupted (e.g., `.git` file missing or corrupted)? `cmd_resume.go` checks for this (line 69-71), but the proposal's "worktree 已存在" path does not mention validation. Also missing: what if two processes call `start` simultaneously for the same slug? |
| Non-functional requirements | 32/40 | Backward compatibility, performance, and observability are covered. Deduction: the "向后兼容" NFR says "start 在 worktree 不存在时的行为完全不变" — but this doesn't address the config rename backward compatibility at all. `copy-files` users are not "backward compatible." Also, security NFR is absent — launching `claude --dangerously-skip-permissions` in a worktree that might have been tampered with is a consideration. |
| Constraints & dependencies | 26/30 | Three dependencies listed with file paths. Deduction: the proposal mentions "依赖现有的 worktree 验证逻辑（resume 已有实现可复用）" but does not note that `cmd_resume.go` performs symlink resolution (`filepath.EvalSymlinks`) that `cmd_start.go` does not — this is a dependency detail that matters for implementation. |

**Attack Points:**
8. **Missing corrupted worktree scenario**: The proposal assumes worktree existence implies worktree validity. But a directory could exist without being a valid git worktree (partial creation failure, manual directory creation, etc.). This is a gap between the proposal's scenarios and the actual validation in `cmd_resume.go`.
9. **Config rename breaks backward compatibility NFR**: The NFR section claims "向后兼容" but the config rename is explicitly non-backward-compatible ("直接替换，不保留任何旧字段兼容逻辑"). This is a direct contradiction within the document.
10. **Missing concurrency scenario**: No mention of race condition when two `start` calls hit the same slug simultaneously — the "check exists → create" pattern in the current code is not atomic.

---

### D5. Solution Creativity: 72 / 100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 26/40 | The proposal explicitly states "不是业界首创" and cites `mkdir -p` as the inspiration. This is a straightforward application of an existing pattern. The innovation is recognizing the semantic shift from "create worktree" to "start working" — which is more insight than creativity. |
| Cross-domain inspiration | 24/35 | `mkdir -p`, `touch`, `kubectl apply` are cited — these are all within the CLI idempotency domain. No cross-domain inspiration (e.g., from IDE behavior, from web UX patterns, from database migration patterns). |
| Simplicity of insight | 22/25 | "核心洞察是 `start` 的语义应聚焦于「启动新会话」而非「创建 worktree」" is genuinely elegant. The insight is simple and correct — the command name "start" naturally maps to "start working" rather than "create infrastructure." |

**Attack Points:**
11. **No creative solution for the config rename**: The rename is purely mechanical — "copy-files" to "includes" because Claude Code uses ".worktreeinclude". No creative solution for how to handle the transition, how to detect stale configs, or how to make the naming future-proof.

---

### D6. Feasibility: 85 / 100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 35/40 | Confirmed: `cmd_start.go` lines 107-110 contain the exact "error if exists" block that needs to change. The resume command's validation logic is directly reusable. Deduction: the proposal says "修改 cmd_start.go 中的「已存在则报错」分支" — but the actual implementation requires more than just flipping an error to a skip. It needs to also skip the entire git branch creation + worktree add flow (lines 134-196), which is a significant section of the function. The proposal understates the code path complexity. |
| Resource & timeline | 28/30 | "8 个源文件 + 对应测试文件，约 120 行变更" — verified: 8 Go source files contain CopyFiles/copy-files references. However, the proposal lists `errors.go` as one of the 8 files, but no `errors.go` exists in the worktree command directory. The actual 8th file would be either `config_auto.go` or `register.go`. Deduction: -2 for listing a non-existent file. "预计 2 小时内完成" is reasonable. |
| Dependency readiness | 22/30 | "无外部依赖" is stated. But this is incomplete: the `forgeconfig` package has `config_auto.go` (reflection-based config access) that uses field names. Renaming `CopyFiles` to `Includes` in the struct requires updating the reflection-based path resolution as well. This internal dependency is not mentioned. |

**Attack Points:**
12. **`errors.go` listed but does not exist in worktree cmd directory**: The proposal says the changes involve `errors.go`（错误消息更新）, but `ls forge-cli/internal/cmd/worktree/` shows no `errors.go`. Error messages are likely in `cmd_start.go` itself (line 108-109) or in `base/errors.go`. This is a factual inaccuracy in the resource assessment.
13. **Reflection-based config access not mentioned as dependency**: `forgeconfig/config_reflect.go` uses reflection to resolve `worktree.copy-files` via field names. Renaming the struct field requires updating reflection paths, but this file is not listed in the scope.

---

### D7. Scope Definition: 73 / 80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 27/30 | Five concrete in-scope items, each describing a specific behavior change. Deduction: item "输出区分性日志" is slightly vague — "区分性" without specifying the exact output format. |
| Out-of-scope explicitly listed | 22/25 | Six out-of-scope items. Deduction: "修改 worktree 的创建/验证逻辑本身" is out of scope, but the proposal requires adding validation logic (checking if the existing directory is a valid git worktree) — this is arguably new validation logic, not just reusing existing. |
| Scope is bounded | 24/25 | Clear bounded scope with 2-hour timeline. The two changes (idempotent + rename) are both small and well-bounded. |

**Attack Points:**
14. **Out-of-scope contradicts in-scope**: "修改 worktree 的创建/验证逻辑本身" is out of scope, but the proposal requires the existing-directory path to validate that the directory is actually a valid git worktree (not just a stale directory). This is new validation logic.

---

### D8. Risk Assessment: 70 / 90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 24/30 | Three risks listed. Deduction: Missing risk — config rename breaks existing `.forge/config.yaml` files silently. When a user upgrades forge and their config has `copy-files`, the YAML key won't match `includes`, and their files simply won't be copied. This is the highest-impact risk and it's completely absent. Also missing: the risk of the `--no-launch` path now succeeding (exit 0) where it previously errored — scripts checking exit codes will see a behavioral change. |
| Likelihood + impact rated | 22/30 | Ratings are present. Deduction: "用户脚本依赖「已存在则报错」的行为" is rated L/M — but the current behavior has been the only behavior since the command existed. The likelihood of scripts depending on it could easily be M, not L. The rating seems optimistic rather than honest. |
| Mitigations are actionable | 24/30 | Mitigations are present and generally actionable. Deduction: "在 release notes 中标注行为变更" is the mitigation for script breakage — but release notes are a passive mitigation. No active mitigation (e.g., deprecation warning period, `--strict` flag) is considered. |

**Attack Points:**
15. **Missing the highest-impact risk: config silent breakage**: The proposal mandates "代码中不存在任何 copy-files / CopyFiles 兼容逻辑" (Success Criteria). This means users with existing `copy-files` configs get zero migration — their config key is silently ignored, their files aren't copied. No risk entry for this.
16. **Script breakage likelihood underrated**: Current behavior is the only behavior that has ever existed. Scripts built against this CLI have no alternative. Likelihood should be M, not L.

---

### D9. Success Criteria: 62 / 80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 22/30 | Most criteria use specific keywords ("entering existing worktree", exit code 0). Deduction: "worktree 不存在时的行为与当前完全一致（回归测试通过）" is vague — "完全一致" is not measurable without specifying what "一致" means. Which specific behaviors must be preserved? Also, "配置项正常工作" is vague — what does "正常工作" mean precisely? |
| Coverage is complete | 20/25 | SC covers most in-scope items. Deduction: The in-scope item "当 worktree 已存在时忽略 --source-branch，输出提示信息" has a SC entry but it only says "被忽略并输出 warning" — it doesn't specify the warning content. The in-scope item "输出区分性日志" maps to the SC about stderr keywords, but the specific keyword format ("entering existing worktree" vs "created new worktree") is only in the SC, not in the in-scope description — a minor disconnect. |
| SC internal consistency | 20/25 | SC entries are generally consistent with each other. Deduction: The SC says "代码中不存在任何 copy-files / CopyFiles 兼容逻辑" which means users with existing configs are broken. But another SC says "worktree 不存在时的行为与当前完全一致" — these two are in tension if "完全一致" includes config behavior. For a user who has `copy-files` in config and creates a new worktree after upgrade, the behavior is NOT consistent because the config key changed. This is an internal contradiction. |

**Attack Points:**
17. **"完全一致" contradicts config rename**: SC says both "行为与当前完全一致" and "不存在任何 copy-files / CopyFiles 兼容逻辑". If the config key changes without backward compat, then the behavior for users with existing configs is NOT consistent — their files won't be copied. These two criteria cannot both be satisfied for existing users.
18. **"正常工作" is not measurable**: SC "worktree.includes 配置项正常工作" — "正常工作" could mean anything. Does it mean: files are copied? Files listed in `includes` exist in worktree after creation? The YAML key is recognized? This is too vague for a testable criterion.

---

### D10. Logical Consistency: 72 / 90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 30/35 | The idempotent start directly solves problem #1 (no way to enter existing worktree with fresh session). The config rename addresses problem #2 (naming inconsistency). Deduction: the config rename is bundled with the idempotent start, but these are independent changes. The logical coupling is not argued — why must they ship together? |
| Scope <-> Solution <-> SC aligned | 22/30 | Generally aligned. Deduction: Scope says "修改 forge worktree start <slug> 在 worktree 已存在时的行为" but the SC also covers `--no-launch` and `--interactive` paths which are not explicitly in the "In Scope" list. The scope section lists "输出区分性日志" but doesn't mention the specific stderr keyword format that the SC mandates. |
| Requirements <-> Solution coherent | 20/25 | Requirements map to solution. Deduction: The NFR says "向后兼容" but the solution explicitly breaks backward compatibility for the config field. This is a coherence gap. |

**Attack Points:**
19. **Two independent changes bundled without justification**: The idempotent start and the config rename are orthogonal. The proposal does not argue why they must be in the same change. They could be separate PRs with separate reviews, reducing risk.
20. **NFR "向后兼容" contradicts solution's "直接替换"**: The NFR section claims backward compatibility, but the config rename is explicitly non-backward-compatible. This is the most significant logical inconsistency in the document.

---

## Blindspot Analysis

### [blindspot] Issues the rubric misses:

1. **Operational rollout plan**: The rubric does not evaluate whether the proposal includes a rollout strategy (deprecation notices, versioning, phased rollout). This proposal changes CLI behavior with no transition period.

2. **User communication plan**: Beyond "release notes," there's no plan for how existing users learn about the behavioral change and config rename. The rubric's Risk Assessment dimension doesn't penalize lack of user communication strategy.

3. **Testing strategy for behavioral change**: The proposal doesn't describe how regression tests will be structured for a behavioral change — existing tests assert the current "error on exists" behavior. Those tests need to be updated, and new tests for the idempotent path need to be written. The rubric doesn't evaluate test planning at the proposal stage.

4. **Coupling justification**: The rubric does not penalize proposals that bundle independent changes without justification. The two changes in this proposal (idempotent start + config rename) are logically independent and should arguably be separate proposals.

5. **Config migration completeness**: The rubric's Feasibility dimension doesn't specifically evaluate whether config/schema migrations are complete. The proposal ignores the need to handle existing user configs.

---

## Bias Detection Report

**Annotated regions** (marked with `<!-- pre-revised: {severity} -->`):
- Line 62 (medium): Observability NFR paragraph
- Line 68 (high): Config structure modification constraint
- Line 94-95 (medium): Resource & timeline paragraph
- Line 141 (medium): Success criteria for stderr keywords
- Line 147-149 (medium): Success criteria for --no-launch and --interactive

**Annotated regions**: 4 attack points / 5 paragraphs = density 0.80
- Attack #3 (config migration — from revised line 68 context, pre-revised: high)
- Attack #13 (reflection dependency — from revised line 68 context)
- Attack #17 (SC contradiction — from revised lines 141, 147-149)
- Attack #18 (vague SC — from revised line 147-149)

**Unannotated regions**: 16 attack points / 28 paragraphs = density 0.57

**Ratio (annotated/unannotated)**: 1.40

**Interpretation**: Annotated regions show slightly higher attack density (1.40x). This suggests a mild bias toward scrutinizing revised sections more heavily. However, the attacks on annotated regions are substantive (SC contradictions, missing dependencies), not nitpicks. The pre-revision improved the document in those areas but introduced new issues (SC internal contradictions between revised entries and unrevised entries).

**Conflict-with-pre-revision tags**: None detected in this evaluation.

---

## Summary of All Attacks

1. [D1] Evidence item 4 is speculative — "如果未来支持 glob 模式或目录" — must separate hypotheticals from current evidence
2. [D1] Two orthogonal problems bundled — no justification for simultaneous delivery
3. [D2/Annotated] No migration path for existing config — "直接替换，不保留任何旧字段兼容逻辑" — must address user impact
4. [D2] `--interactive` scenario underspecified — "正常进入，因为 slug 来自交互选择" — must describe full user experience
5. [D3] Straw-man alternative: "新增 `open` 子命令" — dismissed in one sentence without real analysis
6. [D3] Config rename has zero industry benchmarking — only idempotency is benchmarked
7. [D3] Trade-off severity mismatch — comparison says "轻微" but risk table says M impact
8. [D4] Missing corrupted worktree scenario — existence does not imply validity
9. [D4] Config rename contradicts backward compatibility NFR — "向后兼容" vs "直接替换"
10. [D4] Missing concurrency scenario — check-then-act is not atomic
11. [D5] No creative solution for config rename transition
12. [D6/Annotated] `errors.go` listed but does not exist in worktree directory — factual inaccuracy
13. [D6/Annotated] Reflection-based config access (`config_reflect.go`) not mentioned as dependency
14. [D7] Out-of-scope contradicts in-scope — validation logic is both out of and in scope
15. [D8] Missing highest-impact risk: existing configs silently broken by rename
16. [D8] Script breakage likelihood underrated (L should be M)
17. [D9/Annotated] "完全一致" SC contradicts "no compat logic" SC — both cannot be true for existing users
18. [D9/Annotated] "正常工作" is not measurable — SC must specify exact expected behavior
19. [D10] Two independent changes bundled without justification
20. [D10] NFR "向后兼容" contradicts solution "直接替换" — most significant logical inconsistency
