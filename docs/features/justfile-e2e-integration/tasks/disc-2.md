---
id: "disc-2"
title: "Fix: add justfile prerequisite check to run-e2e-tests SKILL.md"
priority: "P0"
dependencies: []
status: 
breaking: true
---

# disc-2: Fix: add justfile prerequisite check to run-e2e-tests SKILL.md

TC-013 fails because plugins/forge/skills/run-e2e-tests/SKILL.md does not reference a justfile existence check or /init-justfile. The test expects the skill to prompt the user to run /init-justfile when the Justfile is missing the required e2e-setup/e2e-verify recipes. Add a prerequisite check in the Prerequisites section of SKILL.md that verifies the Justfile exists and contains the e2e-setup recipe, and prompts the user to run /init-justfile if it is missing.
