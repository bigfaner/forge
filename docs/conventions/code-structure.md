---
title: "Code Structure Conventions"
domains: [nesting, indentation, early-return, flat-control-flow, readability]
---

# Code Structure Conventions

### TECH-code-structure-001: Prefer Flat Control Flow Over Deep Nesting

**Requirement**: Avoid deep nesting (3+ levels). Use early returns, guard clauses, and extracted helper functions to keep the main logic at minimal indentation. Each tab level adds cognitive load — flatten aggressively.
**Scope**: [CROSS]
**Source**: /learn entry 2026-05-18

**Pattern to avoid** (4 levels):
```
func process() {
    if cond1 {
        if cond2 {
            if cond3 {
                if cond4 {
                    // actual logic buried here
                }
            }
        }
    }
}
```

**Preferred pattern** (1 level, guard clauses):
```
func process() {
    if !cond1 { return ... }
    if !cond2 { return ... }
    if !cond3 { return ... }
    if !cond4 { return ... }
    // actual logic at minimal indentation
}
```

**For switch/case**: prefer flat switch bodies. Avoid nested if/else inside case arms — extract to helpers when logic exceeds 5-10 lines.
