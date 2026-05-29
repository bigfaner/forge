---
type: doc.fix
category: doc
identity:
  - ID
  - Title
  - Priority
  - EstimatedTime
  - Description
  - SourceTaskID
context:
  - SourceFiles
  - TestResults
---
---
id: "{{.ID}}"
title: "{{.Title}}"
priority: "P0"
estimated_time: "{{.EstimatedTime}}"
dependencies: []
status: pending
breaking: false
type: "doc.fix"
---

# {{.Title}}

## Root Cause

{{.Description}}

## Reference Files

- Source: {{.SourceFiles}}
- Error details: {{.TestResults}}

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

When this task is recorded as completed via `task record`, the source task {{.SourceTaskID}} is automatically restored to pending if all its dependencies are completed.
