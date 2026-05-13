---
id: "T-test-1b"
title: "Evaluate e2e Test Cases"
priority: "P1"
estimated_time: "30min"
dependencies: ["T-test-1"]
type: "test-pipeline.eval-cases"
scope: "all"
noTest: true
mainSession: true
---

# Evaluate e2e Test Cases

Execute this test pipeline task.

## Main Session Instructions

1. Invoke the `/eval-test-cases` skill for the `forge-cli-v3` feature to evaluate the test cases document.
2. Record the result using the `/record-task` skill.
