---
feature: "Forge Architecture Simplification"
---

# User Stories: Forge Architecture Simplification

## Story 1: Safe Concurrent Task Operations

**As a** Forge CLI Developer
**I want to** run multiple forge task commands concurrently without data corruption risk
**So that** I can safely run parallel agents without fear of index.json or state.json corruption

**Acceptance Criteria:**
- Given two concurrent `forge task claim` invocations
- When both attempt to claim the same task
- Then exactly one succeeds and the other receives a lock conflict error (not a corrupted index)
- Given `forge task index` is running while `forge task submit` writes to the same index
- When both complete
- Then index.json is valid JSON with both changes applied (no truncation)
- Given index.json is malformed (truncated JSON)
- When `forge task claim` is called
- Then an AIError is returned with Hint "index.json may be corrupted, run forge task index to rebuild" (not a panic)
- Given lock acquisition times out (5s)
- When `forge task submit` attempts to write
- Then an AIError is returned with Action "retry in a few seconds"

---

## Story 2: Correct State Machine Enforcement

**As a** Forge Plugin Developer
**I want to** trust that task state transitions follow consistent rules regardless of which command triggers them
**So that** my skills and agents can rely on predictable task lifecycle behavior

**Acceptance Criteria:**
- Given a task with status `completed`
- When `forge task submit` is called on it
- Then the command returns an error ("task already completed, create a subtask if re-work needed") instead of silently overwriting（不可逆）
- Given a task with status `rejected`
- When `forge task reopen <id>` is called
- Then the task transitions to `pending`
- Given a task with status `completed`
- When `forge task reopen <id>` is called
- Then the command returns an error ("task already completed, create a subtask if re-work needed")
- Given `forge task status <id>` is called with a status argument
- When the command runs
- Then it returns an error ("task status is read-only. Use forge task submit to complete a task.")

---

## Story 3: Actionable Error Messages

**As a** Forge End User
**I want to** receive structured, actionable error messages from all forge commands
**So that** I can quickly understand what went wrong and how to fix it

**Acceptance Criteria:**
- Given `forge worktree remove` fails due to uncommitted changes
- When the error is displayed
- Then the output includes Code, Message, Cause, Hint, and Action fields (not raw `fmt.Errorf` text)
- Given `forge task submit` encounters a lock conflict
- When the error is displayed
- Then the output uses AIError format with Hint "retry in a few seconds" (not raw stderr + exit 1)
- Given any forge command fails
- When the user checks the output
- Then the error format is consistent across all commands (no mixed fmt.Errorf/AIError)

---

## Story 4: Eval Pipeline Never Loses Work

**As a** Forge Plugin Developer
**I want to** run eval pipelines safely knowing that failed iterations won't destroy my original documents
**So that** I can iterate aggressively without risk of data loss

**Acceptance Criteria:**
- Given eval reaches max iterations without meeting target score
- When the final report is generated
- Then original documents are restored from the Step 1 backup (reviser changes are rolled back)
- Given the eval scorer produces malformed output
- When the orchestrator attempts to extract the score
- Then the pipeline halts with an error message (not a crash or silent ignore)
- Given the reviser is modifying documents
- When the reviser receives its prompt
- Then it can see the same project context (conventions, business rules) as the scorer

---

## Story 5: Complete Configuration Control

**As a** Forge End User
**I want to** view and modify all forge configuration fields via CLI commands
**So that** I don't need to manually edit YAML files for common configuration changes

**Acceptance Criteria:**
- Given `forge config set auto.cleanCode true` is executed
- When `forge config get auto.cleanCode` is called
- Then it returns `true`
- Given `forge config get auto.e2eTest` is called
- When the config file has this field set
- Then the current value is displayed (currently only auto.gitPush is queryable)
- Given `forge config init` is run
- When the wizard completes
- Then all 4 auto fields (e2eTest, consolidateSpecs, cleanCode, gitPush) are configured (currently the bufio path skips auto entirely)
- Given config.yaml is empty or missing
- When `forge config get auto.gitPush` is called
- Then a meaningful error is returned (not a panic or empty output)

---

## Story 6: Quality Gate Actually Works

**As a** Forge End User
**I want to** the quality gate fix-task mechanism to correctly track, cap, and restore tasks
**So that** the auto-restore pipeline functions as intended

**Acceptance Criteria:**
- Given a quality gate creates a fix-task for step "2.1"
- When the fix-task is created
- Then its SourceTaskID is the actual blocked task ID (not `"quality-gate:2.1"` sentinel)
- Given 3 fix-tasks exist for a step but all are completed
- When the quality gate evaluates whether to create another fix-task
- Then it allows creation (cap counts active only, not lifetime)
- Given quality gate is invoked with no feature configured
- When the command runs
- Then it returns exit code 1 with an error message (not silent exit 0)

---

## Story 7: Reopen Rejected or Skipped Tasks

**As a** Forge CLI Developer
**I want to** reopen tasks that were rejected or skipped
**So that** I can retry tasks when the rejection was premature or the skip was a mistake

**Acceptance Criteria:**
- Given a task with status `rejected`
- When `forge task reopen <id>` is called
- Then the task transitions to `pending`
- Given a task with status `skipped`
- When `forge task reopen <id>` is called
- Then the task transitions to `pending`
- Given a task with status `completed`
- When `forge task reopen <id>` is called
- Then the command returns an error ("task already completed, create a subtask if re-work needed")
- Given a task with status `in_progress`
- When `forge task reopen <id>` is called
- Then the command returns an error ("task is not rejected or skipped")
