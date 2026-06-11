---
id: "doc-fix-1"
title: "Fix: TYPE not listed as extractable field in claim output docs"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: false
type: "doc.fix"
---

# Fix: TYPE not listed as extractable field in claim output docs

## Root Cause

doc.review AC failure: TYPE field output by forge task claim but not documented in extractable fields list of execute-task.md and run-tasks.md

## Reference Files

- Source: plugins/forge/commands/execute-task.md,plugins/forge/commands/run-tasks.md
- Error details: 

## Content Fix Guidance

When fixing documentation failures, observe these boundaries:

**Scope:**
- Fix only the markdown/content issues identified in the root cause
- Do not modify source code files — this is a documentation-only fix
- Do not run code quality gates (lint, compile, test) — they are irrelevant for doc fixes

**Correct workflow:**
1. Read the failing document and understand the reported issue
2. Identify the specific content problem (broken links, missing sections, incorrect terminology, formatting errors)
3. Apply the minimal fix to resolve the issue
4. Verify the document renders correctly and internal references are valid

When this task is recorded as completed via `task record`, the source task T-review-doc is automatically restored to pending if all its dependencies are completed.
