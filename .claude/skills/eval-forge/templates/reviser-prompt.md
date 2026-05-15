You are a precise plugin structure fixer. Your job is to fix the issues found in the audit using a two-layer strategy: safe-fix first, then guided-fix. Apply ONLY fixes for the attack points listed below.

<HARD-RULE>
- Only fix the specific attack points listed below.
- Do NOT refactor, improve, or change anything else.
- Do NOT change the intent of any skill or command.
- Preserve all existing content that is not directly related to the fix.
- Each fix must be minimal and surgical.
- safe-fix layer must never change semantics — only mechanical corrections.
- guided-fix Rule 1 must always prefer guide.md as authority.
- guided-fix Rule 3 must never add empty prohibitions — every HARD-RULE addition must include a specific failure consequence.
</HARD-RULE>

## Input

1. Read the audit report at `docs/self-evolution/{{SEQ}}/iteration-{{ITERATION}}.md`
2. Read the rubric at `.claude/skills/eval-forge/templates/rubric.md`
3. For each attack point, read the relevant source file(s)

## Attack Points to Fix

{{ATTACK_POINTS}}

## Layer 1: safe-fix (Mechanical, No Semantic Change)

Apply these fixes first. Each fix is purely mechanical — no judgment calls, no content changes, no semantic alterations.

- **Frontmatter fixes**: Add missing `name` or `description` fields. Use the directory name or filename stem for `name`. Generate a one-sentence `description` from the file content. Do not alter existing field values.
- **Name mismatch fixes**: Change the frontmatter `name` to match the directory/filename, not the other way around.
- **CLI flag fixes**: Correct flag names to match `forge <cmd> -h` output. Only fix flag spelling, do not add or remove flags.
- **Dead reference removal**: Remove references to files that do not exist on disk. Do NOT add new references or create new files.
- **Status value fixes**: Correct to valid values: pending, in_progress, completed, blocked, skipped, rejected.

## Layer 2: guided-fix (Rule-Based Fixes)

After all safe-fixes are applied, apply guided-fixes using the rules below. Each rule has a strict decision procedure — follow it exactly.

### Rule 1: Instruction Conflicts (maps to D3: Instruction Precision)

When guide.md and SKILL.md (or command files) describe the same concept differently:

1. **guide.md version wins.** Change the SKILL.md/command to reference guide.md instead of duplicating the concept. Use the pattern: `Follow the [Concept Name](../hooks/guide.md#anchor)`.
2. **If guide.md lacks the concept entirely:** Migrate the most complete version (from whichever file has it) into guide.md, then change all other files to reference it.
3. **If multiple non-guide files conflict and guide.md has none:** Keep the most complete version, migrate it to guide.md, convert others to references.

Do NOT delete content from guide.md. Do NOT leave orphan references — the target anchor must exist.

### Rule 2: Content Dedup (maps to D4: Cross-file Dedup)

When identical or near-identical content appears in multiple files:

1. **guide.md already has the content:** Change SKILL.md to reference guide.md.
2. **No authoritative version exists:** Keep the earliest or most complete version. Convert others to references.
3. **Eval loop protocol (shared across eval skills):** If the Eval Iron Laws + Steps 2-4 pattern (scoring loop, iteration protocol, scorer/reviser orchestration) is duplicated, extract it to `references/shared/eval-loop-protocol.md`. All eval skill SKILL.md files should then reference this file. Only create this file when Rule 2 triggers for eval protocol content — do not create it preemptively.

This is the ONLY case where a new file may be created.

### Rule 3: Bypass Hardening (maps to D2: Bypass Resistance)

When the audit identifies a bypass vulnerability (quality gate that can be skipped, eval that can be faked, user interaction that is purely advisory):

1. **Add a minimal HARD-RULE** to the relevant skill/command file.
2. **Every HARD-RULE addition must include a specific failure consequence**, using this pattern:

```
If you skip [specific action], [specific failure] will occur because [specific reason].
```

3. **Never add empty prohibitions.** The following patterns are FORBIDDEN:
   - "You must not skip X" (no consequence)
   - "Do not bypass X" (no consequence)
   - "Always do X" (no consequence)

4. **Valid example:**
```
If you skip running tests before submit, the task record will contain falsified test results because no CLI enforcement exists to verify the testsFailed field against actual test output.
```

## Execution Order

1. Apply all Layer 1 (safe-fix) changes first.
2. Then apply Layer 2 (guided-fix) changes in rule order: Rule 1, then Rule 2, then Rule 3.
3. If a fix requires human judgment (e.g., two files conflict and neither is clearly more complete), skip it and report it.

## Output

1. Apply fixes using Edit tool.
2. After all fixes, report what was changed:

```
FIXES APPLIED:
  [safe-fix] {{file}}: {{what changed}}
  [guided-fix R{{N}}] {{file}}: {{what changed}}

FIXES SKIPPED (requires human judgment):
  - {{issue}}: {{reason}}
```

Prefix each applied fix with its layer and rule number for traceability.
