# Example: Asking Clarifying Questions

Use `AskUserQuestion` tool for ALL questions.

## Example Usage

### Identify target users (ask early — feeds User Stories)

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "Who are the primary users of this feature?",
    "header": "Target Users",
    "multiSelect": true,
    "options": [
      {"label": "Developer / Engineer", "description": "Writes code, runs CLI tools"},
      {"label": "AI Agent", "description": "Automated agent consuming APIs or CLI"},
      {"label": "Product Manager", "description": "Monitors progress via Web UI"},
      {"label": "End User", "description": "Uses the product directly"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```

### Identify primary goal

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "What is the primary goal of this feature?",
    "header": "Goal",
    "options": [
      {"label": "Improve performance", "description": "Make existing operations faster"},
      {"label": "Add new capability", "description": "Enable something not currently possible"},
      {"label": "Fix pain point", "description": "Address a specific user frustration"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```
