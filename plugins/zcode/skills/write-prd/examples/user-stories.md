# Example: Writing User Stories

Derive stories directly from the target users identified in Background. One story per key workflow.

## Structure

```
### Story N: [Short title]
**As a** [user role]
**I want to** [specific action]
**So that** [concrete benefit]

**Acceptance Criteria:**
- Given [precondition]
- When [action taken]
- Then [expected outcome]
```

## Example (Agent Task Center)

Background identified two users: Developer (Web UI) and AI Agent (CLI).

```markdown
### Story 1: Developer monitors task progress
**As a** developer managing multiple projects
**I want to** view all tasks across projects on a single kanban board
**So that** I can spot blockers and track overall progress without switching contexts

**Acceptance Criteria:**
- Given I open the Web UI
- When I select a project and feature filter
- Then I see tasks grouped by status (pending / in_progress / completed / blocked)

### Story 2: AI Agent claims a task
**As an** AI agent (Claude Code, Codex, etc.)
**I want to** claim the next available task via CLI
**So that** I can start working without manual coordination or file conflicts

**Acceptance Criteria:**
- Given TASK_REMOTE_URL is set and a pending task exists
- When I run `task claim`
- Then the task status changes to in_progress and is assigned to my agent_id
- And no other agent can claim the same task concurrently
```

## Tips

- If the PRD only describes APIs or technical requirements, **still write stories** — ask "who calls this API and why?"
- Each story's AC becomes a row in the final Acceptance Criteria section
- Aim for 1–3 stories; more than 5 usually means scope is too broad
