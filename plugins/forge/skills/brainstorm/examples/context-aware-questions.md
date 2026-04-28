# Example: Context-Aware Questions

When analysis finds concrete facts, derive questions that reference them. Below are worked examples showing the full flow: analysis → finding → targeted question with recommendation.

## Example 1: Overlapping Feature Found

**User's idea:** "I want to add a skill that automatically reviews code quality."

**Analysis results:**
- `Grep "review"` → found `skills/eval-design/SKILL.md` with design review capability
- `Grep "quality"` → found `skills/simplify/SKILL.md` that checks code quality
- No existing proposals for code review skill

**Finding:** Two existing skills partially cover this — `eval-design` reviews design docs, `simplify` checks code quality.

**Derived question:**

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "I found two existing skills that touch on quality: `eval-design` reviews design docs, and `simplify` checks code quality. Is your idea covering a gap between these, or something different?",
    "header": "Existing Overlap",
    "multiSelect": false,
    "options": [
      {"label": "Gap between them (Recommended)", "description": "Neither fully covers automated code review as you envision it"},
      {"label": "Extending simplify", "description": "Build on the existing simplify skill with deeper analysis"},
      {"label": "Different scope", "description": "Your idea targets a different aspect of code review entirely"},
      {"label": "Didn't know about those", "description": "Need to understand what they do before deciding"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```

**Contrast with generic template:** A template would ask "What problem are you trying to solve?" — the user already told you. The context-aware version immediately shows you've done homework and advances the conversation.

## Example 2: Related Recent Work

**User's idea:** "I want to restructure the docs directory."

**Analysis results:**
- `Bash git log --oneline -20` → commits `3fb6d85` and `da96aa5` about directory structure redesign
- `Glob docs/superpowers/specs/*directory*` → found `2026-04-09-directory-structure-redesign.md`
- `Glob docs/superpowers/plans/*directory*` → found `2026-04-11-directory-structure-redesign.md`

**Finding:** There's already active work on directory restructuring with both a spec and a plan.

**Derived question:**

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "I see active work on directory restructuring — there's a spec (2026-04-09) and a plan (2026-04-11), plus recent commits aligning docs to the new structure. Is your idea extending that redesign, or addressing a different aspect?",
    "header": "Related Work",
    "multiSelect": false,
    "options": [
      {"label": "Continuing that work (Recommended)", "description": "Building on the existing redesign spec and plan"},
      {"label": "Different directories", "description": "Your restructuring targets areas not covered by the existing plan"},
      {"label": "Revising the approach", "description": "You want to change the direction of the existing redesign"},
      {"label": "Unaware of it", "description": "Let me review the existing spec/plan first"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```

## Example 3: User Workflow Overlap

**User's idea:** "I want to add a batch import feature for task records."

**Analysis results:**
- `Grep "import"` → found existing single-record creation flow in task CLI
- `Grep "batch"` → no batch operation exists
- `Glob docs/features/**/prd-spec.md` → found a PRD mentioning bulk operations as out-of-scope

**Finding:** Single-record creation exists; a previous PRD explicitly excluded batch operations, suggesting this was considered but deferred.

**Derived question:**

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "I found that a previous PRD explicitly excluded batch operations (listed as out-of-scope). Is your idea addressing the same need that was deferred, or a different batch scenario?",
    "header": "Deferred Need",
    "multiSelect": false,
    "options": [
      {"label": "Same deferred need (Recommended)", "description": "The deferred batch import is now a priority"},
      {"label": "Different scenario", "description": "Your batch import targets a different workflow or user role"},
      {"label": "Broader than before", "description": "Covers the deferred case plus additional scenarios"},
      {"label": "Unaware of it", "description": "Let me review that PRD first"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```

**Note:** This question focuses on business context (why it was deferred, whether the need has changed) rather than technical details (file format, import mechanism).

## Example 4: Greenfield / No Findings

**User's idea:** "I want to add a changelog generation skill."

**Analysis results:**
- `Grep "changelog"` → no results
- `Glob docs/proposals/*` → no changelog proposals
- `Bash git log --oneline -20` → no changelog-related commits

**Finding:** Nothing related found. This is genuinely greenfield.

**Fallback question (derived from "no findings"):**

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "I couldn't find any existing changelog tooling in the codebase. Is this for generating changelogs from git history, from task records, or from PRD/user-story changes?",
    "header": "Source",
    "multiSelect": false,
    "options": [
      {"label": "Git history (Recommended)", "description": "Parse commits, tags, and merge messages"},
      {"label": "Task records", "description": "Aggregate from docs/features/*/tasks/"},
      {"label": "PRD diff tracking", "description": "Track changes to PRD specs over time"},
      {"label": "Manual entries", "description": "Provide a template for manual changelog entries"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```

**Note:** Even with no findings, the question is specific to the idea (changelog source) rather than generic ("what problem?"). The "no findings" result itself tells you something — the user is likely starting from scratch and needs help scoping.

## Summary Pattern

```
Analysis finding → Specific, referenced question with recommendation
No finding → Idea-specific scoping question with recommendation (still better than generic template)
```

Always prefer questions that prove you've looked at the codebase. Users trust and engage more when they see the skill has done its homework. For each question, mark your recommended answer as the first option with `(Recommended)` — this reduces user burden and shows you've thought through the analysis.
