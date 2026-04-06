---
name: learn-lesson
description: Use when you have solved an error or discovered a useful pattern. Extracts reusable knowledge from the current session.
---

# Learn Lesson

## Overview

从当前会话中提取可复用的知识点，记录到 `docs/lessons/`。

**核心原则**：记录"下次遇到类似问题可以怎么处理"，不是"我做了什么"。

## When to Use

**Trigger conditions:**
- Solved a non-trivial error with reusable insights
- Discovered a pattern/technique worth remembering
- User explicitly requests `/learn-lesson`

**Skip when:**
- Trivial typo fixes
- One-off environment issues
- Information already documented elsewhere

## Workflow

```
1. Identify lesson → 2. Classify category → 3. Write doc → 4. Commit
```

## Step 1: Identify Lesson

- 遇到的问题（症状）
- 根本原因
- 解决方案
- **可复用的知识点**

## Step 2: Classify Category

| Category | Prefix | Example |
|----------|--------|---------|
| Debugging | `debug-` | `debug-race-condition.md` |
| Architecture | `arch-` | `arch-dependency-direction.md` |
| Tooling | `tool-` | `tool-go-test-coverage.md` |
| Pattern | `pattern-` | `pattern-error-wrapping.md` |
| Gotcha | `gotcha-` | `gotcha-context-cancellation.md` |

## Step 3: Write Document

Output: `docs/lessons/<category-prefix><slug>.md`

Template:
```markdown
# <Title>

## Problem
<!-- 症状描述 -->

## Root Cause
<!-- 根本原因 -->

## Solution
<!-- 解决方案 -->

## Key Takeaway
<!-- 可复用的知识点 -->
```

## Step 4: Commit

```bash
git add docs/lessons/<filename>.md
git commit -m "docs(lessons): add <title>"
```

## Common Mistakes

| Mistake | Correction |
|---------|------------|
| Recording "what I did" | Focus on "what to do next time" |
| Too specific | Generalize to reusable pattern |
| Missing root cause | Always include why |
