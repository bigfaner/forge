# forge Skills Registry

## Skills

| Skill | Description |
|-------|-------------|
| brainstorm | Explore vague ideas before formalizing into a PRD |
| write-prd | Formalize requirements into a structured PRD |
| eval-prd | Evaluate PRD quality against standards |
| tech-design | Create technical design after PRD is finalized |
| ui-design | Create UI design specifications for UI features |
| eval-design | Evaluate tech-design.md quality |
| breakdown-tasks | Break down design into executable tasks |
| record-task | Record task execution result |
| git-commit | Create git commits with Conventional Commits format |
| eval-harness | Evaluate harness health |
| eval-proposal | Evaluate a proposal document with scoring |
| eval-ui | Evaluate UI design with four-perspective scoring |
| improve-harness | Improve harness based on evaluation report |
| learn-lesson | Extract reusable knowledge from current session |
| gen-test-cases | Generate test cases from PRD acceptance criteria |
| gen-test-scripts | Generate executable e2e test scripts |
| run-e2e-tests | Execute e2e test scripts and generate results report |
| graduate-tests | Migrate feature test scripts to the regression suite (tests/e2e/) |
| consolidate-specs | Extract business rules and tech specs from feature docs, user confirms before integrating to project-level dirs |

## Commands

Invoked via `/command-name`. Commands live in `commands/` directory.

| Command | Description |
|---------|-------------|
| gen-sitemap | Generate page element map (sitemap.json) for test pipeline and ui-design |
| execute-task | Execute a single task with full quality gate verification |
| fix-bug | Fix a bug with minimal changes and quality gate verification |
| run-tasks | Dispatch and execute all pending tasks in a feature |
| init-forge | Initialize forge directory structure in a new project |
| init-justfile | Generate project-specific justfile with standard vocabulary |
| record-decision | Record a technical decision with context and rationale |
| extract-design-md | Extract design document from conversation context |
| simplify-skill | Simplify and optimize a skill definition |
