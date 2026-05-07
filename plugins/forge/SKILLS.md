# forge Skills Registry

## Skills

| Skill | Description |
|-------|-------------|
| brainstorm | Explore vague ideas through collaborative dialogue and produce a structured proposal document before formalizing into a PRD |
| write-prd | Formalize requirements into a structured PRD through collaborative dialogue |
| eval-prd | Evaluate PRD quality against standards |
| tech-design | Create technical design after PRD is finalized |
| ui-design | Create UI design specifications with style selection and HTML prototype generation |
| eval-design | Evaluate a tech design with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents |
| breakdown-tasks | Break down tech-design into executable tasks |
| record-task | Record task execution result and update task status |
| git-commit | Create git commits with Conventional Commits format |
| eval-harness | Evaluate harness health |
| eval-proposal | Evaluate a proposal with 100-point scoring and adversarial iterations until target score is met |
| eval-ui | Evaluate UI design with four-perspective scoring |
| improve-harness | Improve harness based on evaluation report |
| learn-lesson | Extract reusable knowledge from current session |
| gen-test-cases | Generate test cases from PRD acceptance criteria |
| gen-test-scripts | Generate executable e2e test scripts |
| run-e2e-tests | Execute e2e test scripts and generate results report |
| graduate-tests | Migrate feature test scripts to the regression suite (tests/e2e/) |
| consolidate-specs | Extract business rules and tech specs from feature docs into preview files, detect overlaps with existing knowledge, user confirms before integrating to project-level dirs. |

## Commands

Invoked via `/command-name`. Commands live in `commands/` directory.

| Command | Description |
|---------|-------------|
| gen-sitemap | Auto-generate and maintain sitemap.json for a web app. Uses agent-browser to explore routes, capture accessibility tree, and discover dynamic states |
| execute-task | Execute single task with focused TDD workflow |
| fix-bug | Systematically fix a bug using TDD workflow — reproduce, write failing tests, fix, verify |
| run-tasks | Autonomous task dispatcher that continuously claims tasks and dispatches to subagents |
| init-forge | Build and install the task-cli tool |
| init-justfile | Scaffold a Justfile with standard forge targets for the current project |
| record-decision | Record an architecture/technical decision to docs/decisions/ at any stage of development |
| extract-design-md | Analyze a web app's visual style and generate a DESIGN.md for use with ui-design skill |
| simplify-skill | Refactor skill files by extracting templates/examples to separate files |
