# forge Skills Registry

## Skills

| Skill | Description |
|-------|-------------|
| brainstorm | Use when a user has a vague idea or feature request and needs to explore it before formalizing into a PRD. Outputs a structured proposal document. |
| write-prd | Use when user provides requirements or feature requests that need to be formalized into a structured PRD document through collaborative dialogue. |
| eval-prd | Evaluate a PRD document with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents. |
| tech-design | Use after PRD (and UI design if applicable) is finalized to create technical design with architecture and implementation details. |
| ui-design | Use after PRD ui-functions are defined to create UI design specifications. Auto Eval UI after design, then generates HTML prototype. |
| eval-design | Evaluate a tech design document with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents. |
| breakdown-tasks | Use when the technical design is finalized to break down into executable tasks. Creates task files based on technical design. |
| record-task | Use after completing a task to create its execution record and update task status. |
| git-commit | Use when creating git commits. Ensures commit messages follow Conventional Commits format. |
| eval-harness | Harness health evaluation with 100-point rubric scoring. Based on OpenAI's harness engineering practices. Produces scored report with actionable P0/P1/P2 priorities. |
| eval-proposal | Evaluate a proposal document with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents. Specify target score and max iterations. |
| eval-ui | Evaluate a UI design document with 100-point scoring from four stakeholder perspectives (User/Designer/Developer/PM), then run adversarial iterations until target score is met. |
| improve-harness | Dynamically implement harness improvements from eval-harness report. Reads P0/P1/P2 priorities and fixes each finding. |
| learn-lesson | Use when you have solved an error or discovered a useful pattern. Extracts reusable knowledge from the current session. |
| gen-test-cases | Generate structured test cases from PRD acceptance criteria. Classifies by type (UI/API/CLI) with full traceability to PRD sections. |
| eval-test-cases | Evaluate test-cases.md for downstream executability with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents. |
| gen-test-scripts | Generate executable TypeScript e2e test scripts from test cases. Uses @playwright/test for all tests (no node:test or node:assert). Playwright for UI, fetch for API, child_process for CLI. |
| run-e2e-tests | Execute e2e test scripts and generate a results report. Runs UI tests via Playwright, API tests via fetch, CLI tests via child_process. Produces evidence-backed pass/fail report. |
| graduate-tests | Migrate feature test scripts to the regression suite (tests/e2e/). Agent-driven: reads scripts, analyzes content, decides classification, splits/merges as needed, rewrites imports, creates graduation marker. |
| consolidate-specs | Extract business rules and tech specs from feature docs into preview files, detect overlaps with existing knowledge, user confirms before integrating to project-level dirs. |

## Commands

Invoked via `/command-name`. Commands live in `commands/` directory.

| Command | Description |
|---------|-------------|
| gen-sitemap | Auto-generate and maintain sitemap.json for a web app. Uses agent-browser to explore routes, capture accessibility tree, and discover dynamic states. Preserves element IDs across runs. |
| execute-task | Execute single task with focused TDD workflow. |
| fix-bug | Systematically fix a bug using TDD workflow — reproduce, write failing tests, fix, verify. Ensures the bug is captured by tests before any code changes. |
| run-tasks | Autonomous task dispatcher that continuously claims tasks and dispatches to subagents. |
| init-forge | Build and install the task-cli tool. |
| init-justfile | Scaffold a Justfile with standard forge targets for the current project. |
| record-decision | Record an architecture/technical decision to docs/decisions/ at any stage of development. |
| extract-design-md | Analyze a web app's visual style and generate a DESIGN.md for use with ui-design skill. |
| simplify-skill | Refactor skill files by extracting templates/examples to separate files. |
