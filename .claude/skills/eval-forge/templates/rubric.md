# Forge Plugin Runtime Reliability Rubric

**Total: 1000 points**
**Report template:** `.claude/skills/eval-forge/templates/report.md`

## What This Rubric Measures

Runtime reliability of the forge plugin — not just structural consistency, but whether components work correctly together at runtime. Measures workflow completeness, bypass resistance, instruction precision, redundancy, reference integrity, and structural conventions.

## Scoring Methodology: 4-Phase Process

This rubric is designed for a 4-phase evaluation methodology:

**Phase 1: Construct Workflow Graph (D1)**
1. Read the ground-truth workflow specs embedded in Dimension 1 below
2. Scan all skill/command/agent files, extract actual prerequisites, outputs, gate points
3. Compare spec vs actual, find breakpoints, dead ends, unreachable states
4. Build intra-skill step graphs: extract step sequences, conditional skip targets, and detection points; verify temporal ordering consistency (see criterion 1f)

**Phase 2: Per-Node Adversarial Testing (D2)**
1. List every gate/confirm point per node
2. For each gate, assume you are a lazy agent — ask "how can I bypass this?"
3. Check whether HARD-RULE has enforceable consequences (not just words)
4. Eval loops: check whether scoring must be done by independent subagent
5. Quality gates: check whether CLI-level enforcement exists

**Phase 3: Per-File Precision Review (D3 + D4)**
1. Instruction conflicts (highest priority): behavioral cross-reference with CLI source code, instruction chain tracing, then cross-file same-concept comparison
2. Step ambiguity: does each SKILL.md step have a single unambiguous interpretation
3. Incomplete conditionals: does every if-then have an else
4. Undefined variables: do agent-filled variables have source annotations
5. Content redundancy: check by A/B/C category standards

**Phase 4: Baseline Integrity (D5 + D6)**
1. Reference integrity
2. Frontmatter, eval templates, name alignment

## Scoring Dimensions

| Dimension | Points |
|-----------|--------|
| 1. Workflow Completeness | 280 |
| 2. Bypass Resistance | 280 |
| 3. Instruction Precision | 280 |
| 4. Cross-file Dedup | 30 |
| 5. Reference Integrity | 80 |
| 6. Structural Convention | 50 |
| **Total** | **1000** |

Score allocation logic: Runtime reliability (D1+D2+D3) = 840 pts (84%), Baseline integrity (D4+D5+D6) = 160 pts (16%).

## Deduction Tiers

| Severity | Penalty | Examples |
|----------|---------|---------|
| Critical | -20 | Workflow breakpoint, eval integrity bypass, instruction conflict, temporal ordering contradiction |
| High | -15 | Dangling reference, invalid manifest transition, missing eval template |
| Medium | -10 | Step ambiguity, conditional without else, content duplication, missing prerequisite check |
| Low | -5 | Missing frontmatter field, name-directory mismatch, advisory-only prohibition, context inflation |

**Floor rule:** Each sub-criterion has a minimum score of 0. Total deductions for a sub-criterion cannot result in a negative score.

**Cross-dimension deduction:** The same issue may be deducted in multiple dimensions if it violates independent criteria (e.g., an advisory-only HARD-RULE is both a D2 bypass vector and a D3 incomplete conditional). Bypass resistance (D2) and instruction precision (D3) are independent dimensions.

## Dimensions

### 1. Workflow Completeness (280 pts)

Verify the full-chain state graph: every skill's prerequisites, outputs, and successor steps; both quick mode and full mode paths must be complete.

**Ground-Truth Workflow Specs:**

#### Full Mode Pipeline

```
brainstorm → write-prd → eval-prd → [has UI?] → ui-design → eval-ui → prototype → [human review] → tech-design → [db-schema?] → [schema review] → eval-design → breakdown-tasks
```

```
breakdown-tasks → forge task index (auto T-test tasks) → T-test-1 (gen-sitemap + gen-test-cases) → T-test-1b (eval-test-cases) → T-test-2 (gen-test-scripts) → T-test-3 (run-e2e-tests) → T-test-4 (graduate-tests) → T-test-4.5 (verify-regression) → T-test-5 (consolidate-specs)
```

#### Quick Mode Pipeline

```
/quick → /brainstorm → [human confirm] → /quick-tasks → /run-tasks
```

Quick test chain: T-quick-1 (gen-test-cases) → T-quick-2 (gen-test-scripts) → T-quick-3 (run-e2e-tests) → T-quick-4 (graduate-tests) → T-quick-5 (verify-regression) → T-quick-6 (doc-generation-drift). Skips: gen-sitemap, eval-test-cases, consolidate-specs. Note: T-quick-5 and T-quick-6 are CLI-embedded templates (no SKILL.md) — unauditable steps.

#### Manifest Status Machine

```
prd → design → tasks → in-progress → completed
```

Legal transitions only forward. Quick mode starts at tasks (no prd/design).

#### Per-Skill Precondition/Output Matrix

| Skill | Hard Prerequisites | Outputs | Conditionals |
|-------|-------------------|---------|-------------|
| brainstorm | None | `docs/proposals/<slug>/proposal.md` | → eval-proposal (optional) or write-prd |
| write-prd | Optional: proposal.md, sitemap.json | `prd/prd-spec.md`, `prd/prd-user-stories.md`, `prd/prd-ui-functions.md` (if UI), `manifest.md` (status: prd) | has UI → ui-design; no UI → tech-design |
| eval-prd | `prd/prd-spec.md`, `prd/prd-user-stories.md` | `prd/eval/iteration-{N}.md`, `prd/eval/report.md` | score gate pass/fail; delegates to `skills/eval` with `rubrics/prd.md` |
| ui-design | `prd/prd-ui-functions.md` (hard) | `ui/ui-design.md`, `ui/prototype/` | multi-platform → separate files |
| eval-ui | `ui/ui-design.md` | `ui/eval/iteration-{N}.md`, `ui/eval/report.md` | platform → rubric variant; delegates to `skills/eval` with `rubrics/ui-{platform}.md` |
| tech-design | `prd/prd-spec.md` (hard) | `design/tech-design.md`, `design/er-diagram.md` + `design/schema.sql` (if db), `manifest.md` (status: design) | db-schema: yes → mandatory ER+schema |
| eval-design | `design/tech-design.md` | `design/eval/iteration-{N}.md`, `design/eval/report.md` | score gate; delegates to `skills/eval` with `rubrics/design.md` |
| breakdown-tasks | `prd/prd-spec.md` + `design/tech-design.md` (both hard) | `tasks/*.md`, `tasks/index.json`, `manifest.md` (status: tasks) | HAS_UI/NO_UI/HAS_DB/HAS_PLACEMENT tags |
| gen-test-cases | `prd/prd-user-stories.md` + `prd/prd-spec.md` (both hard) | `testing/test-cases.md` | profile → interface types |
| eval-test-cases | `testing/test-cases.md` + PRD docs | `testing/eval/iteration-{N}.md` | Step Actionability < 200 blocks; delegates to `skills/eval` with `rubrics/test-cases.md` |
| gen-test-scripts | `testing/test-cases.md` (hard) | `tests/e2e/features/<slug>/` | profile → framework; Step Actionability gate |
| run-e2e-tests | justfile + staging area | `tests/e2e/features/<slug>/results/latest.md` | >30% failure → stop |
| graduate-tests | staging area + PASS results + no marker | `tests/e2e/<module>/`, `.graduated/<slug>` | profile → import rewriting |
| consolidate-specs | PRD + design (both hard) | `specs/`, updated `docs/business-rules/`, `docs/conventions/` | CROSS vs LOCAL |
| /quick | User idea | Orchestrates quick pipeline | >10 tasks → STOP |
| quick-tasks | proposal.md (hard) | `tasks/*.md`, `tasks/index.json`, `manifest.md` | max 10 tasks; generates T-quick-1~6 (T-quick-5 and T-quick-6 are CLI-embedded, no SKILL.md) |
| /run-tasks | `tasks/index.json` | Execution loop | 3 consecutive failures → STOP |
| submit-task | Task executed + record.json | `records/*.md`, updated index.json | --force bypasses gate (not auto-downgrade) |

**Scoring Criteria:**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 1a. Full mode chain complete | 0-90 | Every skill has prerequisite definitions and outputs. Prerequisite files are produced by predecessor skills. Chain has no breakpoints. **Also checks for unauditable steps**: steps in the ground-truth pipeline that have no SKILL.md and no CLI-accessible template (e.g., CLI-embedded task templates). Breakpoint = -20 each. Unauditable step = -10 each. |
| 1b. Quick mode chain complete | 0-40 | Quick chain is complete. Missing step = -20 each. |
| 1c. Conditional branching correct | 0-30 | Every conditional branch has a true-path and false-path. Missing branch = -10 each. |
| 1d. Manifest status transitions valid | 0-30 | Status transitions are legal. Illegal transition = -15 each. |
| 1e. Test lifecycle chain intact | 0-70 | Test chain is complete. **Also checks for unauditable test steps**: test chain steps with no SKILL.md and no CLI-accessible template (e.g., verify-regression embedded in CLI binary). Broken link = -15 each. Unauditable step = -10 each. |
| 1f. Intra-skill temporal ordering | 0-30 | Every conditional skip's detection point precedes the step it modifies. Detection-after-skip-target = -15 each. |

**1f Detailed Check Procedure:**

For each SKILL.md that defines numbered steps with conditional skip/fast-path logic:

1. **Extract step sequence**: Identify all numbered steps (Step 0, Step 1, Step 2, ...) and their execution order.
2. **Identify conditional skip targets**: Find instructions like "skip Step N if condition C" or "Steps N-M are unnecessary when condition C".
3. **Locate detection points**: Find where condition C is evaluated — look for "Detection:" paragraphs, "After Step N" markers, "During Step N" markers, or mermaid decision nodes.
4. **Verify temporal consistency**: For each (detection_point, skip_target) pair, check that `detection_step <= skip_target`. A detection point at "After Step 4a" cannot skip "Step 0" because Step 0 executes before the detection.
5. **Check mermaid diagrams**: If a mermaid flow diagram exists, trace the decision node's incoming edge — does it connect from a step that precedes the skipped step?

**Common violation pattern:**
```
Detection: "After Step N, if condition C..."
Skip target: "Step M where M < N"
→ Contradiction: condition evaluated at Step N, but Step M already executed.
```

**Example (from real fix):**
```
Before fix: "Detection: After Step 4a (Business Tasks), if every business task
was docs-only, skip Step 0 (Resolve Profile)"
→ Step 0 executes BEFORE Step 4a, so the skip is never actionable.

After fix: "Detection: During Step 1 (Read Documents), after scanning all input
artifacts — if every design element targets non-compilable files, skip Step 0"
→ Detection at Step 1, Step 0 not yet executed, skip is actionable.
```

### 2. Bypass Resistance (280 pts)

Phase 2 uses a two-stage approach: independent discovery first (systematically find all enforcement gaps), then regression verification (confirm previously-fixed TEXT-FIXABLE vectors stay fixed).

**5 Bypass Types:**

| Type | Points | Examples |
|------|--------|----------|
| Type 2: Skip quality gates | 0-80 | `--force` overrides, noTest bypasses, missing justfile silently passes, self-reported metrics |
| Type 3: Fake eval results | 0-80 | Skip scorer subagent, tamper with scores, general-purpose fallback loses adversarial constraints |
| Type 1: Skip mandatory interaction | 0-45 | User confirmations are advisory-only HARD-RULE text, no mechanical enforcement |
| Type 4: Skip required steps | 0-35 | Conditional requirements depend on self-reporting, gates only trigger when prerequisites exist |
| Type 5: Lazy shortcuts | 0-30 | Prohibitions are purely advisory, metrics are self-reported, direct file editing bypasses CLI |

**Scoring Criteria:**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 2a. Quality gate enforcement | 0-80 | Is each gate point enforced by CLI or merely advisory text. Zero enforcement with no documented rationale = -15 each. |
| 2b. Eval integrity | 0-80 | Does each eval skill require independent subagent scoring. Does the decision gate parse structured output. Can the main session fake the score. Weakness = -25 each. |
| 2c. User interaction enforcement | 0-45 | Does each confirmation point have a mechanical enforcement mechanism. Purely advisory = -10 each. |
| 2d. Required step enforcement | 0-35 | Do conditional requirements have downstream verification. No verification = -10 each. |
| 2e. Prohibition enforcement | 0-40 | Does each HARD-RULE prohibition have a mechanical check. Purely advisory = -5 each. |

**Bypass Discovery Protocol (Phase 2a):**

Systematically enumerate all enforcement gaps. Do NOT read the regression guard table before completing this phase — discovery must be unbiased.

**Step 1 — Enumerate enforcement claims.** Scan every SKILL.md, command file, and guide.md for:
- `<HARD-RULE>`, `<HARD-GATE>`, `<EXTREMELY-IMPORTANT>` tags and what each demands
- "must", "required", "mandatory", "do not" statements
- Quality gate references (compile, fmt, lint, test)
- User confirmation points (AskUserQuestion, explicit approval)
- Conditional enforcement (gates that only trigger under certain conditions — e.g., only when eval report exists)
- Self-reported fields (record.json metrics, task frontmatter values like `noTest`, `db-schema`)

**Step 2 — Trace enforcement chain.** For each claim, determine:
- Where enforced? (CLI command / hook / text-only / absent)
- Consequence of ignoring? (blocks / warns / nothing)
- Bypass path? (`--force` / missing prerequisite / alternative tool / self-declared value)
- Cross-reference with CLI source code: read `forge-cli/internal/cmd/submit.go` (quality gates, `--force`, `noTest` handling, `validateRecordData`). Confirm CLI behavior matches SKILL.md descriptions.

**Step 3 — Classify gaps.** Each enforcement gap found:
- **[ARCHITECTURAL]**: Requires code-level change to enforce. Score and report, but do NOT generate reviser fix tasks. Adding HARD-RULE text for architectural bypasses is counterproductive — it inflates context without changing agent behavior.
- **[TEXT-FIXABLE]**: Can be mitigated by adding conditional branches, fallback paths, or actionable instructions. Valid reviser targets.

**Regression Guard (Phase 2b):**

The following TEXT-FIXABLE bypass vectors were previously identified. Verify each one's current state. Do not deduct for vectors that have been fixed.

| Vector | Type | Severity | Description |
|--------|------|----------|-------------|
| BV-1.1 | T1 | MED | brainstorm user approval can be skipped (pure HARD-RULE) |
| BV-1.2 | T1 | MED | write-prd user approval can be skipped |
| BV-1.3 | T1 | MED | /quick user confirmation can be skipped |
| BV-1.4 | T1 | MED | ui-design prototype review can be skipped |
| BV-1.5 | T1 | MED | tech-design DB schema review can be skipped |
| BV-1.7 | T1 | MED | consolidate-specs spec integration confirmation can be skipped |
| BV-4.2 | T4 | MED | gen-test-scripts Step Actionability gate only triggers when eval report exists |
| BV-4.5 | T4 | MED | placement validation depends on sitemap existence, skipped if absent |

**Bypass classification for fix strategy:**

- **ARCHITECTURAL**: Cannot be fixed by adding text to SKILL.md/command files. These require code-level changes (CLI enforcement, cryptographic verification, etc.). Score them and report them, but do NOT generate reviser fix tasks. Adding HARD-RULE text for architectural bypasses is counterproductive — it inflates context without changing agent behavior.
- **TEXT-FIXABLE**: Can be mitigated by adding conditional branches, fallback paths, or actionable instructions. These are valid targets for the reviser.

### 3. Instruction Precision (280 pts)

**Priority order: Instruction conflicts > Step ambiguity > Incomplete conditionals > Undefined variables**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 3a. Instruction conflicts | 0-100 | Three conflict sources: **(1) Behavioral**: guide.md/SKILL.md describes behavior that CLI source code doesn't implement (e.g., guide says "lint retries once" but CLI has no retry logic). Cross-reference `forge-cli/internal/cmd/submit.go` (behavioral logic), `forge-cli/pkg/just/just.go` (quality gate implementation). **(2) Chain**: Step N's output mismatches Step N+1's expected input; conditional states not handled by subsequent steps. **(3) Cross-file**: same concept described differently across guide.md, SKILL.md, command files. Conflict = -25 each. **Check behavioral first, then chain, then cross-file.** |
| 3b. Step ambiguity | 0-60 | SKILL.md steps must have a single unambiguous interpretation. Vague verbs ("check tests", "verify quality") without specific commands = -10 each. Also flag steps where chain tracing reveals the agent must make unstated choices. |
| 3c. Incomplete conditionals | 0-50 | Every if-then must have an else path or an explicit "skip" instruction. Missing else = -10 each. **Implicit-else exception:** if the false-path is the natural default (normal execution continues, zero-value default, or no-op), no explicit else is required. Only flag if-then patterns where the false-path requires distinct handling but none is documented. |
| 3d. Variable resolution clarity | 0-40 | Agent-filled variables must have source annotations. CLI-filled variables must match `prompt.go` typeToTemplate. Undefined agent variable = -10 each. |
| 3e. Narrative inflation | 0-30 | Text that inflates context without changing agent behavior: consequence/rationale paragraphs (why vs what), stale code/function references (pointing to moved or non-existent files), redundant re-explanation of what a table/step/code-block already states. Instance = -5 each. **Exempt**: content inside `<HARD-RULE>`, `<HARD-GATE>`, `<EXTREMELY-IMPORTANT>` tags. |

**CLI-Filled Variables (from `prompt.go` Synthesize — do NOT mark as "undefined"):**

| Variable | Source | Used by |
|----------|--------|---------|
| `{{TASK_ID}}` | task.ID | All 14 templates |
| `{{TASK_FILE}}` | feature.GetTaskFile() | All 14 templates |
| `{{SCOPE}}` | task.Scope | Most templates |
| `{{PHASE_SUMMARY}}` | PhaseDetect() | 10 templates |
| `{{FEATURE_SLUG}}` | SynthesizeOpts.FeatureSlug | 4 templates |
| `{{PROFILE}}` | task.Profile | 4 test pipeline templates |

### 4. Cross-file Dedup (30 pts)

**Three categories with different standards:**

<PLUGIN-PORTABILITY>
Forge is a Claude Code plugin deployed to users' projects. Plugin SKILL.md and command files run in arbitrary working directories. Cross-file relative paths (`../../`) won't resolve at runtime. Therefore:
- Duplication that serves plugin portability is **acceptable and should NOT be deducted**.
- Only deduct when dedup is feasible without breaking portability (e.g., content within the same skill directory, or in non-plugin project-level files like `.claude/skills/`).
- guide.md cannot be referenced via relative path from plugin SKILL.md files. Inlining guide.md content in plugin files is a necessary trade-off.
</PLUGIN-PORTABILITY>

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 4a. Content copy | 0-10 | Identical/near-identical text blocks appearing in 3+ files. **Plugin portability exception:** Known instances serve plugin portability — flag as INFO but do NOT deduct. Only deduct when dedup is feasible within the same skill directory or in non-plugin files. Actionable instance = -10 each. |
| 4b. guide.md vs SKILL.md overlap | 0-10 | guide.md is the single source of truth. **But plugin files cannot reference guide.md via relative paths** — they must inline the content. Only deduct for non-plugin files or when referencing is feasible. Actionable duplication = -10 each. |
| 4c. Unreasonable inline | 0-10 | Content that has its own file but is also inlined in SKILL.md. Judgment: does the agent need to see the full content in a single context (reasonable inline) vs can it be obtained via Read tool on demand (should reference). **Plugin portability:** inlining content from outside the skill directory is reasonable when the alternative would break at runtime. Unreasonable inline = -10 each. |

**Known Redundancy Instances:**

| Instance | Category | Portability-Required? | Deduct? |
|----------|----------|----------------------|---------|
| "Step 0: Resolve Profile" | A | YES (9 plugin SKILL.md files) | NO |
| Eval Iron Laws + Steps 2-4 | A | NO (consolidated to 1 `skills/eval/SKILL.md`) | NO |
| Eval report shared sections | A | YES (5 plugin report.md files) | NO |
| Quality gate sequence | B | YES (plugin skills need it inline) | NO |
| Scope resolution paraphrase | B | Partial | Only if in non-plugin file |

### 5. Reference Integrity (80 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 5a. Agent references valid | 0-20 | `forge:<agent>` or `subagent_type` points to an existing file in `plugins/forge/agents/`. Dangling = -10 each. |
| 5b. Template references valid | 0-20 | Template paths in SKILL.md point to existing files. Dangling = -10 each. |
| 5c. Cross-skill references valid | 0-15 | `invoke /<name>` points to an existing skill/command. Dangling = -10 each. |
| 5d. Hook references valid | 0-15 | Paths and CLI commands in hooks.json exist. Dangling = -10 each. |
| 5e. Shared reference paths valid | 0-10 | `plugins/forge/references/*` paths in SKILL.md or commands must exist. Dangling = -5 each. |

### 6. Structural Convention (50 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 6a. Frontmatter completeness | 0-25 | SKILL.md has name + description. Command has name + description. Agent has name + description + model. Missing = -5 each. |
| 6b. Eval template convention | 0-15 | `skills/eval/SKILL.md` exists, `skills/eval/rubrics/` contains rubric files for each eval type, and each `commands/eval-*.md` delegates to `Skill("eval", ...)`. Missing rubric = -10 each. |
| 6c. Name-directory alignment | 0-10 | Skill name = directory name, command name = filename. Mismatch = -5 each. |
