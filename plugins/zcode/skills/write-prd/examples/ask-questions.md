# Example: Asking Clarifying Questions

Use `AskUserQuestion` tool for ALL questions.

## Example Usage

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
