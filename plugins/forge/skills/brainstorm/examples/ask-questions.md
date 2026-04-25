# Example: Generic Question Templates (Fallback)

These are generic templates. Use them **only when analysis yields no useful findings** (e.g., greenfield project, empty repo). Prefer context-aware questions — see `context-aware-questions.md`.

## Phase 1: Problem Exploration

### Clarify the problem

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "What problem are you trying to solve?",
    "header": "Problem",
    "multiSelect": false,
    "options": [
      {"label": "Pain point / friction", "description": "Users struggle with something that exists today"},
      {"label": "Missing capability", "description": "Something users need but can't do today"},
      {"label": "Performance / quality", "description": "Something works but not well enough"},
      {"label": "New opportunity", "description": "A new possibility enabled by tech or market change"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```

### Identify who is affected

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "Who experiences this problem most acutely?",
    "header": "Target Users",
    "multiSelect": true,
    "options": [
      {"label": "Developer / Engineer", "description": "Writes code, runs CLI tools"},
      {"label": "AI Agent", "description": "Automated agent consuming APIs or CLI"},
      {"label": "Product Manager", "description": "Monitors progress via dashboards"},
      {"label": "End User", "description": "Uses the final product directly"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```

### Understand current workaround

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "How are people working around this problem today?",
    "header": "Workaround",
    "multiSelect": false,
    "options": [
      {"label": "Manual process", "description": "Doing things by hand, copy-paste, etc."},
      {"label": "Partial automation", "description": "Some scripting or tooling, but incomplete"},
      {"label": "Third-party tool", "description": "Using an external service or plugin"},
      {"label": "No workaround", "description": "Problem is just accepted or ignored"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```

## Phase 2: Solution Exploration

### Validate solution direction

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "What does success look like? How would you know this is solved?",
    "header": "Success",
    "multiSelect": false,
    "options": [
      {"label": "Measurable metric", "description": "Speed, accuracy, throughput improvement"},
      {"label": "User satisfaction", "description": "Fewer complaints, better feedback"},
      {"label": "New workflows enabled", "description": "Users can do things they couldn't before"},
      {"label": "Reduced complexity", "description": "Fewer steps, less cognitive load"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```

### Scope boundary

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "What is the MVP scope? What would be nice-to-have but not essential?",
    "header": "Scope",
    "multiSelect": false,
    "options": [
      {"label": "Core only", "description": "Just the minimum to validate the idea"},
      {"label": "Core + polish", "description": "Core functionality with good UX"},
      {"label": "Full vision", "description": "Complete feature set from the start"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```

## Phase 3: Challenge Assumptions

### Pressure-test the idea

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "What if we did nothing — would this problem resolve itself or get worse?",
    "header": "Urgency",
    "multiSelect": false,
    "options": [
      {"label": "Gets worse", "description": "Problem grows without intervention"},
      {"label": "Stays the same", "description": "Status quo is painful but stable"},
      {"label": "Gets resolved", "description": "Other changes may address this naturally"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```

### Simpler alternative

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "Is there a simpler way to achieve 80% of the value with 20% of the effort?",
    "header": "Simpler Alt",
    "multiSelect": false,
    "options": [
      {"label": "Yes, there's a shortcut", "description": "A simpler approach could work"},
      {"label": "Partial, but insufficient", "description": "Simpler helps but doesn't fully solve it"},
      {"label": "No, full solution needed", "description": "The problem requires a complete solution"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```

## Tips

- **These are fallbacks** — always prefer context-aware questions derived from codebase analysis.
- **Adapt questions to context** — these are templates, not a script.
- **Skip questions that are already answered** — if the user already stated the problem clearly, move on.
- **Dig deeper when answers are vague** — if the user picks "Pain point" but can't articulate it, ask "Can you describe a specific moment when this pain was most noticeable?"
- **Never ask more than one question at a time** — respect the user's cognitive load.
