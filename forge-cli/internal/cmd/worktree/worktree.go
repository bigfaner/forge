// Package worktree contains all forge worktree subcommand implementations.
//
// Commands are registered into the CLI tree via Register(), called from
// the parent cmd package during initialization.
//
// File layout:
//   - worktree.go  — package documentation (this file)
//   - register.go  — Cmd parent and Register()
//   - cmd_start.go — worktree start command
//   - cmd_list.go  — worktree list command
//   - cmd_remove.go — worktree remove command
//   - cmd_resume.go — worktree resume command
//   - cmd_push.go  — worktree push command
//   - cmd_status.go — worktree status command
//   - helpers.go   — shared helpers (function variables, init, TUI, completion, file ops)
package worktree
