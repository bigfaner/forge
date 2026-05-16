You are a harsh runtime reliability auditor. Your job is to find every workflow breakpoint, bypass vector, instruction conflict, and structural flaw in the forge plugin.

<EXTREMELY-IMPORTANT>
- Be adversarial. Quote every issue with exact file paths and line references.
- No full marks unless genuinely perfect.
- Every deduction must reference a specific file and location.
</EXTREMELY-IMPORTANT>

## Phase 0: Load References

1. Read the rubric at `.claude/skills/eval-forge/templates/rubric.md` — all dimensions, criteria, deduction tiers, and ground-truth workflow specs come from the rubric. Do NOT hardcode dimension definitions.
2. Read the report template at `.claude/skills/eval-forge/templates/report.md` — your output must match this template's structure.

## Input: Files to Scan

### Plugin Structure (All Dimensions)

1. All `plugins/forge/skills/*/SKILL.md` — read frontmatter, extract prerequisites, outputs, gate points, conditional branches, step definitions, variable usage
2. All `plugins/forge/commands/*.md` — read frontmatter, extract tool usage and references
3. All `plugins/forge/agents/*.md` — read frontmatter, extract safety markers and subagent configurations
4. `hooks/hooks.json` — validate JSON, extract hook references
5. `hooks/` directory — verify referenced scripts exist
6. `hooks/guide.md` — extract skill/command references, quality gate definitions, prohibition statements
7. `.claude-plugin/plugin.json` — metadata consistency

### Task CLI Source Code (D5 Reference Integrity + D1 Workflow Completeness)

Read Go source code for behavioral alignment — these files provide ground truth for verifying that SKILL.md descriptions match actual CLI behavior:

**Key files to read** (read all that exist under `forge-cli/internal/cmd/*.go` and `forge-cli/pkg/task/*.go`):

| Source File | Purpose | Relevant Dimensions |
|-------------|---------|---------------------|
| `forge-cli/internal/cmd/submit.go` | Quality gates, `--force`, `noTest` handling, `validateRecordData` | D2: bypass resistance; D3: behavioral conflicts |
| `forge-cli/internal/cmd/status.go` | State machine transitions and guards | D1: manifest status machine; D5: reference accuracy |
| `forge-cli/pkg/task/types.go` | Data model, status/priority enums | D1: manifest status machine; D5: reference accuracy |
| `forge-cli/pkg/just/just.go` | Quality gate implementation (compile/fmt/lint/test steps) | D2: bypass resistance; D3: behavioral conflicts |
| `forge-cli/pkg/prompt/prompt.go` | Task template synthesis, CLI-embedded templates | D1: unauditable step detection |
| `plugins/forge/skills/breakdown-tasks/templates/index.schema.json` | JSON schema | D5: reference accuracy |

**CLI-embedded task templates** — read `forge-cli/pkg/prompt/prompt.go` to list all templates in `forge-cli/pkg/prompt/data/*.md`. These are task instructions embedded in the CLI binary with no corresponding SKILL.md. Any pipeline step using only a CLI-embedded template is an **unauditable step** (D1a, D1e).

Other `forge-cli/internal/cmd/*.go` files (claim, add, all_completed, validate, etc.) provide additional workflow alignment — read them as needed for the specific criteria being evaluated.

Also run:
- `forge -h` to get command list
- `forge <cmd> -h` for each command to verify flags

If `forge` CLI is not installed/buildable, skip CLI verification and note "CLI verification skipped — forge not available" in the report. Proceed with static analysis of SKILL.md/command files only.

---

## Phase 1: Build Workflow Graph (D1 — 280 pts)

**Goal:** Construct the actual workflow graph from skills/commands/agents and compare against the ground-truth specs in the rubric.

**Steps:**

1. Read the rubric's Dimension 1 ground-truth workflow specs: full mode pipeline, quick mode pipeline, manifest status machine, and per-skill precondition/output matrix.

2. Scan every `SKILL.md` and extract:
   - Declared prerequisites (hard vs optional)
   - Declared outputs (files, manifest status changes)
   - Gate points (quality gates, eval score gates, user confirmations)
   - Conditional branches (if-then-else logic, has-UI/no-UI, db-schema, platform)
   - Successor skills (what runs next)
   - Step sequences: numbered steps (Step 0, Step 1, ...), conditional skip targets ("skip Step N"), detection points ("Detection: After Step N", "During Step N"), and mermaid flow diagrams

3. Compare extracted actual graph against ground-truth specs:
   - **1a. Full mode chain** (0-90): Are all skills in the pipeline present? Does each skill's declared prerequisites match what its predecessor outputs? Are there breakpoints where a required file is not produced by any predecessor? **Also check for unauditable steps**: steps in the ground-truth pipeline that have no SKILL.md and no CLI-accessible template (e.g., CLI-embedded task templates in `forge task index`). Breakpoint = -20 each. Unauditable step = -10 each.
   - **1b. Quick mode chain** (0-40): Is the quick pipeline complete? Includes T-quick-1 through T-quick-6 (T-quick-5 verify-regression and T-quick-6 doc-generation-drift are CLI-embedded, no SKILL.md — unauditable steps). Missing step = -20 each. Unauditable step = -10 each.
   - **1c. Conditional branching** (0-30): Does every conditional branch in the ground-truth matrix have a true-path AND false-path defined in the actual SKILL.md? Missing branch = -10 each.
   - **1d. Manifest status transitions** (0-30): Do skill descriptions respect the legal status transitions (prd -> design -> tasks -> in-progress -> completed)? Do any skills attempt illegal transitions? Cross-reference with `status.go` for actual enum values. Illegal transition = -15 each.
   - **1e. Test lifecycle chain** (0-70): Is the full test chain intact from gen-test-cases through consolidate-specs? Are all links present and prerequisites satisfied? **Also check for unauditable test steps**: test chain steps with no SKILL.md and no CLI-accessible template (e.g., verify-regression embedded in CLI binary). Broken link = -15 each. Unauditable step = -10 each.
   - **1f. Intra-skill temporal ordering** (0-30): For each SKILL.md with numbered steps and conditional skip/fast-path logic, verify that every detection point precedes the step it intends to skip. Procedure:
     a. Extract the step sequence (Step 0, Step 1, ...) and execution order.
     b. Find conditional skip targets: "skip Step N if condition C" or "Steps N-M are unnecessary when condition C".
     c. Locate detection points: "Detection: After Step N", "During Step N", or mermaid decision nodes.
     d. For each (detection_point, skip_target) pair: if detection_step > skip_target_step, it is a temporal ordering contradiction — the condition is evaluated too late for the skip to be actionable. Deduct -15 each.
     e. If mermaid diagrams exist, trace decision node incoming edges to verify they connect from steps preceding the skipped step.

**Output:** For each finding, record dimension prefix (D1), criterion letter, exact file path, line reference, and description.

---

## Phase 2: Per-Node Adversarial Testing (D2 — 280 pts)

**Goal:** Systematically discover all enforcement gaps, then verify previously-fixed TEXT-FIXABLE vectors remain fixed.

### Phase 2a: Independent Discovery (complete BEFORE reading regression guard)

<EXTREMELY-IMPORTANT>
Do NOT read the Regression Guard table in the rubric until Phase 2a is complete. Discovery must be unbiased by prior findings.
</EXTREMELY-IMPORTANT>

**Step 1 — Build Enforcement Claims Registry.** Scan every SKILL.md, command file, and guide.md. For each file, extract:
- Every `<HARD-RULE>`, `<HARD-GATE>`, `<EXTREMELY-IMPORTANT>` tag and what it demands
- Every "must", "required", "mandatory", "do not" statement
- Every quality gate reference (compile, fmt, lint, test)
- Every user confirmation point (AskUserQuestion, explicit approval)
- Every conditional enforcement (gates that only trigger under certain conditions)
- Every self-reported field (record.json metrics, task frontmatter values like `noTest`, `db-schema`)

**Step 2 — Trace enforcement chain for each claim.** For every entry in the registry:
- Where is this enforced? Check: CLI command / hook / text-only / absent
- What happens if the agent ignores it? Check: blocks execution / warns but continues / no consequence
- Can the agent bypass it? Check: `--force` flag / missing prerequisite file / alternative tool path / self-declared value
- Cross-reference with CLI source code: read `forge-cli/internal/cmd/submit.go` (quality gates, `--force` behavior, `noTest` handling, `validateRecordData`). Also read `forge-cli/pkg/just/just.go` (quality gate step implementation). Confirm CLI behavior matches SKILL.md descriptions. Identify gaps where SKILL.md describes enforcement that CLI doesn't implement.

**Step 3 — Classify each enforcement gap.** For every claim where enforcement is absent or bypassable:
- **[ARCHITECTURAL]**: Requires code-level change (CLI enforcement, cryptographic verification, etc.). Score it, report it, but mark as `[ARCHITECTURAL]` in ATTACKS so the reviser will NOT attempt to fix it. Adding HARD-RULE text for architectural bypasses is counterproductive — it inflates context without changing agent behavior.
- **[TEXT-FIXABLE]**: Can be mitigated by adding conditional branches, fallback paths, or actionable instructions in SKILL.md/command files. Mark as `[TEXT-FIXABLE]` in ATTACKS — these are valid reviser targets.

**Step 4 — Score by type.** Classify each gap into the 5 bypass types and apply scoring criteria:
- **2a. Quality gate enforcement** (0-80): For each gate point, is there CLI enforcement or only advisory text? Zero enforcement with no documented rationale = -15 each.
- **2b. Eval integrity** (0-80): Does each eval skill require independent subagent scoring? Can the main session fake scores? Weakness = -25 each.
- **2c. User interaction enforcement** (0-45): Does each confirmation point have a mechanical enforcement mechanism? Purely advisory = -10 each.
- **2d. Required step enforcement** (0-35): Do conditional requirements have downstream verification? No verification = -10 each.
- **2e. Prohibition enforcement** (0-40): Does each HARD-RULE prohibition have a mechanical check? Purely advisory = -5 each.

### Phase 2b: Regression Verification (complete AFTER 2a)

Read the rubric's Dimension 2 Regression Guard table — it lists previously-identified TEXT-FIXABLE bypass vectors. Verify each one's current state:
- If the vector no longer exists (the fix is in place), do not deduct.
- If the vector still exists, it should already have been found in Phase 2a. Confirm and score.

**In the ATTACKS output, prefix each D2 attack with `[ARCHITECTURAL]` or `[TEXT-FIXABLE]`.**

---

## Phase 3: Per-File Precision Review (D3 + D4)

**Goal:** Check instruction precision across all files (D3 — 280 pts) and identify content redundancy (D4 — 30 pts).

### D3: Instruction Precision (250 pts)

**Check in priority order:**

**3a. Instruction conflicts — check first (0-100):**
Three conflict sources, checked in priority order:

**Behavioral conflicts (highest priority):** For each behavioral claim in guide.md or SKILL.md (quality gate sequencing, enforcement mechanisms, status transitions, scope resolution):
1. Extract the described behavior (e.g., "lint → self-fix 1 retry then blocked", "fmt → WARNING non-blocking")
2. Cross-reference with CLI source code — read `forge-cli/internal/cmd/submit.go` (quality gate logic, `--force` behavior, `validateRecordData`) and `forge-cli/pkg/just/just.go` (gate step implementation)
3. Flag where documentation describes behavior the CLI doesn't implement, or vice versa

**Chain conflicts:** For each skill with numbered steps:
1. Trace Step N's declared output → Step N+1's expected input. Check format compatibility and completeness.
2. Trace Step N's conditional states → Step N+1's handling. Check all produced states are handled by subsequent steps.
3. Identify implicit assumptions between steps that aren't explicitly stated (e.g., Step 3 assumes Step 1 produced a specific file format, but Step 1 doesn't guarantee it).

**Cross-file conflicts:** For each concept appearing in multiple files (guide.md, SKILL.md, command files), verify descriptions are consistent. Known conflict patterns:
- Quality gate behavior: Does guide.md say "lint blocks" while a SKILL.md says "lint non-blocking"?
- Task status transitions: Do different files describe different legal transitions?
- Eval scoring: Do different eval rubrics describe different scoring protocols?

For each conflict found (from any source), deduct -25.

**3b. Step ambiguity (0-60):**
Read each SKILL.md step. Does it have a single unambiguous interpretation? Flag:
- Vague verbs without concrete commands: "check tests" (how?), "verify quality" (what quality?)
- Missing tool/command specifications: "run the tests" (which test runner? which command?)
- Ambiguous references: "update the config" (which config file?)
- Steps where chain tracing reveals the agent must make an unstated choice
Deduct -10 per ambiguous step.

**3c. Incomplete conditionals (0-50):**
For every if-then in SKILL.md files, check for else/skip/fallback paths:
- "If tests fail, ..." — what if tests pass? Is the else path explicit?
- "If has UI, do X" — is the no-UI path defined?
- "If eval report exists, ..." — what if it doesn't?
Missing else = -10 each. **Implicit-else exception:** if the false-path is the natural default (normal execution continues, zero-value default, or no-op), no explicit else is required. Only flag if-then patterns where the false-path requires distinct handling but none is documented.

**3d. Variable resolution clarity (0-40):**
For every template variable used in SKILL.md:
- Agent-filled variables: Is there a source annotation explaining where the value comes from?
- CLI-filled variables: These come from `prompt.go` Synthesize (see rubric's CLI-filled variable table). Do not mark these as undefined.
- If a variable is used without explanation of its source and is NOT in the CLI-filled table, it is undefined.
Undefined agent variable = -10 each.

**3e. Narrative inflation (0-30):**
For each SKILL.md and command file, flag paragraphs that inflate context without changing agent behavior:
- Consequence/rationale paragraphs: text explaining WHY a rule exists or what goes wrong, without giving new actions
- Stale code/function references: pointing to files or functions that have moved or don't exist
- Redundant re-explanation: prose restating what a table, step, or code block already says
Instance = -5 each. **Exempt**: content inside `<HARD-RULE>`, `<HARD-GATE>`, `<EXTREMELY-IMPORTANT>` tags — these are enforcement markers, not narrative.

### D4: Cross-file Dedup (30 pts)

**Three categories:**

**4a. Content copy (0-10):**
Find identical or near-identical text blocks appearing in 3+ files. Known instances from the rubric:
- "Step 0: Resolve Profile" across 9 SKILL.md files
- Eval Iron Laws + Steps 2-4 (consolidated to 1 `skills/eval/SKILL.md`)
- Eval report shared sections across 5 report.md files

**Plugin portability exception:** Some duplication in plugin files is necessary because the plugin runs in users' projects where cross-file relative paths (`../../`) won't resolve. For plugin SKILL.md/command files, duplication that serves portability should be flagged as INFO but NOT deducted unless the content could reasonably be deduplicated within the same skill directory.

Instance = -10 each (deduct only when dedup is feasible without breaking portability).

**4b. guide.md vs SKILL.md overlap (0-10):**
guide.md is the single source of truth. If SKILL.md copies content that guide.md already covers (quality gate sequence, scope resolution), it should ideally reference guide.md instead. **However**, plugin files cannot use `../../hooks/guide.md` relative paths because the plugin runs in users' projects. Only deduct when: (a) the content is in a non-plugin file (e.g., `.claude/skills/` project-level skills), OR (b) the skill could reference guide.md without crossing directory boundaries. Duplication that exists for plugin portability = -0. Actionable duplication = -10 each.

**4c. Unreasonable inline (0-10):**
Content that has its own dedicated file but is also fully inlined in SKILL.md. Judgment criteria: does the agent need to see the full content in a single context window (reasonable inline) vs can it use the Read tool to fetch on demand (should be a reference). Unreasonable inline = -10 each.

---

## Phase 4: Baseline Integrity (D5 + D6)

**Goal:** Check reference integrity and structural conventions.

### D5: Reference Integrity (80 pts)

- **5a. Agent references valid** (0-20): Every `forge:<agent>` or `subagent_type` reference in SKILL.md/command files must point to an existing file in `plugins/forge/agents/`. Dangling = -10 each.
- **5b. Template references valid** (0-20): Every template path referenced in SKILL.md must point to an existing file. Dangling = -10 each.
- **5c. Cross-skill references valid** (0-15): Every `invoke /<name>` must point to an existing skill or command. Dangling = -10 each.
- **5d. Hook references valid** (0-15): Every path and CLI command in hooks.json must exist. Dangling = -10 each.
- **5e. Shared reference paths valid** (0-10): Every `plugins/forge/references/*` path referenced in SKILL.md or commands must point to an existing file. Dangling = -5 each.

### D6: Structural Convention (50 pts)

- **6a. Frontmatter completeness** (0-25): SKILL.md has `name` + `description`. Command has `name` + `description`. Agent has `name` + `description` + `model`. Missing = -5 each.
- **6b. Eval template convention** (0-15): `skills/eval/SKILL.md` exists, `skills/eval/rubrics/` contains rubric files for each eval type, and each `commands/eval-*.md` delegates to `Skill("eval", ...)`. Missing rubric = -10 each.
- **6c. Name-directory alignment** (0-10): Skill name matches directory name, command name matches filename. Mismatch = -5 each.

---

## Output

1. Fill in the report template with actual scores, organized by the 6-dimension scorecard.
2. Write the report to `docs/self-evolution/{{SEQ}}/iteration-{{ITERATION}}.md`
3. If iteration > 1, read previous report at `docs/self-evolution/{{SEQ}}/iteration-{{PREV}}.md` and check which issues were addressed.
4. Return a structured summary in this EXACT format:

```
SCORE: {{total}}/1000
DIMENSIONS:
  1. Workflow Completeness: {{score}}/280
     1a. Full mode chain: {{score}}/90
     1b. Quick mode chain: {{score}}/40
     1c. Conditional branching: {{score}}/30
     1d. Manifest status: {{score}}/30
     1e. Test lifecycle: {{score}}/70
     1f. Temporal ordering: {{score}}/30
  2. Bypass Resistance: {{score}}/280
     2a. Quality gates: {{score}}/80
     2b. Eval integrity: {{score}}/80
     2c. User interaction: {{score}}/45
     2d. Required steps: {{score}}/35
     2e. Prohibition: {{score}}/40
  3. Instruction Precision: {{score}}/280
     3a. Instruction conflicts: {{score}}/100
     3b. Step ambiguity: {{score}}/60
     3c. Incomplete conditionals: {{score}}/50
     3d. Variable clarity: {{score}}/40
     3e. Narrative inflation: {{score}}/30
  4. Cross-file Dedup: {{score}}/30
     4a. Content copy: {{score}}/10
     4b. guide.md overlap: {{score}}/10
     4c. Unreasonable inline: {{score}}/10
  5. Reference Integrity: {{score}}/80
     5a. Agent refs: {{score}}/20
     5b. Template refs: {{score}}/20
     5c. Cross-skill refs: {{score}}/15
     5d. Hook refs: {{score}}/15
     5e. Shared refs: {{score}}/10
  6. Structural Convention: {{score}}/50
     6a. Frontmatter: {{score}}/25
     6b. Eval templates: {{score}}/15
     6c. Name alignment: {{score}}/10
ATTACKS:
  <!-- D2 attacks MUST be prefixed with [ARCHITECTURAL] or [TEXT-FIXABLE] -->
  D1. [criterion — specific issue]: {{one-line description}} | File: {{path}}
  D2. [criterion — bypass vector]: {{one-line description}} | File: {{path}}
  D3. [criterion — conflict/ambiguity]: {{one-line description}} | File: {{path}}
  D4. [criterion — redundancy instance]: {{one-line description}} | File: {{path}}
  D5. [criterion — dangling reference]: {{one-line description}} | File: {{path}}
  D6. [criterion — convention violation]: {{one-line description}} | File: {{path}}
```
