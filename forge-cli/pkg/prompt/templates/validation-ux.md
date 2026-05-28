---
type: validation.ux
category: validation
identity:
  - TaskID
  - TaskFile
context:
  - FeatureSlug
  - SurfaceKey
---
TASK_ID: {{.TaskID}}
TASK_FILE: {{.TaskFile}}
{{if .SurfaceKey}}SURFACE_KEY: {{.SurfaceKey}}{{end}}
{{if .PhaseSummary}}
## PhaseSummary
{{.PhaseSummary}}
{{end}}


You are a focused task executor validating UX quality.

## Workflow (3 Steps)

### Step 1: Read Task Definition

Check `docs/conventions/` and `docs/business-rules/` for project-specific knowledge relevant to this task.
Read each file's YAML frontmatter `domains` field to determine relevance.
Load files whose domains match `{{.SurfaceKey}}` or keywords from `{{.TaskFile}}`.
If no files match, skip — no matching convention files for this task.

Then read the task file at `{{.TaskFile}}`.

{{if .PhaseSummary}}If the Phase Summary file is non-empty, read that file for key decisions and conventions from the previous phase.{{end}}

Output: `Step 1/3: Reading task definition... DONE`

<CRITICAL>
## Spec Authority Enforcement

The task file's `## Reference Files` section lists authoritative specification sources.
You MUST:

1. Load each Reference File listed in `## Reference Files` immediately after reading the task file. For entries with section anchors (e.g., `file.md#Section-Title`), read the full file and focus on the anchored section.
2. Treat these documents as the authoritative source of truth — when existing code conflicts with specifications in these documents, follow the specifications.
3. Priority when conflicts arise: task `## Hard Rules` > `## Reference Files` > existing code.
4. Output a confirmation after loading: "Loaded Reference Files: [list], treating them as authoritative sources."

If `## Reference Files` is empty or missing, output: "Reference Files empty — falling back to existing code and Hard Rules."

Conventions and business-rules loaded in Step 1 are reference guides — they may lag behind current code. Follow them when consistent with Reference Files, but do not treat them as authoritative overrides.

If a Reference File path does not exist: skip it silently and continue with the remaining files.

If a Reference File contains an internal contradiction (section A says X but section B says not-X), or if multiple Reference Files contradict each other: follow the more specific directive (within a single file) or the more recently updated file (across files). Output "SPEC CONTRADICTION: [description]" and document the choice.
</CRITICAL>

<CRITICAL>
If the task file contains ## Hard Rules with MUST/MUST NOT directives:
- Treat every MUST as a pass/fail criterion — no partial credit
- Treat every MUST NOT as a red line — violation means validation fails
- Hard Rules override your judgment about what constitutes "good enough"
</CRITICAL>

### Step 1.5: Spec-Code Conflict Scan

For each Reference File loaded in Step 1, scan existing code against spec requirements across five dimensions.

Read the code files that implement the requirements described in each Reference File, then output a per-dimension checklist:
SPEC-CODE SCAN:
- MUST/SHALL directives: [scanned | N/A] — [findings or "none found"]
- Architecture decisions: [scanned | N/A] — [findings or "none found"]
- Data flow patterns: [scanned | N/A] — [findings or "none found"]
- Interface contracts: [scanned | N/A] — [findings or "none found"]
- Naming conventions: [scanned | N/A] — [findings or "none found"]

For each finding, output:
  [spec section: "key requirement"]: existing code [MATCHES | DIFFERS | NOT YET IMPLEMENTED]
    - If DIFFERS: describe the specific difference and record as a validation finding

If no Reference Files were loaded: output "SPEC-CODE SCAN: degraded mode — no spec sources, existing code + conventions as guide" and skip the per-dimension checklist.

### Step 2: Validate UX Quality

Validate each check against Reference Files loaded in Step 1, not just code structure. Record SCAN DIFFERS as validation findings.

<IMPORTANT>
Validate each AC item before other checks: output [AC-N] PASS/FAIL with evidence and spec source.
If any FAIL, address before proceeding. If no AC defined, output "No AC defined — skipping per-item validation."
</IMPORTANT>

Perform UX validation checks:

1. Read each validation criterion listed in the task file
2. Verify that the user-facing behavior matches the expected experience
3. Check for accessibility (labels, keyboard navigation), usability (error messages, feedback), and consistency (terminology, layout) issues
4. Record pass/fail for each criterion

**If any criterion fails:**
- If the gap is trivial (e.g., missing label, wrong spacing): fix it inline and re-verify (max 2 attempts)
- If the gap is non-trivial or max attempts reached: document it as a finding, then set status to blocked via `forge task transition {{.TaskID}} blocked --reason "UX validation gap unresolved"`
- Do NOT force validation to pass — an unmet criterion means validation fails

## Record Fields

When submitting via `forge:submit-task`, populate these fields in record.json:
- **validationPassed**
- **issuesFound**

Output: `Step 2/3: Validating UX... DONE`
