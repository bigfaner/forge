---
id: "1"
title: "Sort lesson list by file modification time in reverse chronological order"
priority: "P2"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "enhancement"
mainSession: false
---

# 1: Sort lesson list by file modification time in reverse chronological order

## Description
`forge lesson` lists lessons in filesystem (alphabetical) order. Users cannot quickly find recently added lessons. Add sorting by file modification time (newest first) to the `Discover()` function's return value.

## Reference Files
- `docs/proposals/lesson-reverse-chronological/proposal.md` — Source proposal
- `forge-cli/pkg/lesson/lesson.go` — Discover() function to add sorting
- `forge-cli/internal/cmd/lesson.go` — runLessonList() consumer (no changes expected)

## Acceptance Criteria
- [ ] `forge lesson` output is sorted by file modification time, newest first
- [ ] Lessons without valid modification times sort to the end of the list
- [ ] `forge lesson <name>` detail view still works correctly
- [ ] New sorting logic has unit test coverage

## Hard Rules
- Use Go standard library `sort.Slice` only — no new dependencies
- Do not change the lesson file format or introduce new fields

## Implementation Notes
- Add `sort.Slice` call in `Discover()` after building the lessons slice, sorting by `os.FileInfo.ModTime()` descending
- The `FileInfo` is already available during iteration (from `os.ReadDir` + `os.Stat`). Consider capturing ModTime during the existing loop to avoid a second stat call.
- Key risk: file modification time != creation time. Acceptable for git-managed projects per proposal assessment.
