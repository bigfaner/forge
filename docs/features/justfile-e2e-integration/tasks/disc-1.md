---
id: "disc-1"
title: "Fix: add e2e-setup and e2e-verify recipes to Justfile"
priority: "P0"
dependencies: []
status: 
breaking: true
---

# disc-1: Fix: add e2e-setup and e2e-verify recipes to Justfile

TC-003, TC-004, TC-010, TC-011, TC-012, TC-020 all fail because the project Justfile is missing the e2e-setup and e2e-verify recipes. The init-justfile command documents these recipes but they have not been added to the actual Justfile at the project root. Add both recipes: e2e-setup (idempotent npm install + playwright install chromium, exits 1 if tests/e2e/package.json missing) and e2e-verify --feature <slug> (scans tests/e2e/<slug>/*.spec.ts for // VERIFY: markers, exits 1 with file:line output if found, exits 0 with OK message if clean).
