# AskUserQuestion Example

Example of using AskUserQuestion for clarifying technical decisions.

```xml
<function_calls>
<invoke name="AskUserQuestion">
<parameter name="questions">[
  {
    "question": "How should we handle connection failures?",
    "header": "Error Handling",
    "options": [
      {"label": "Retry with backoff", "description": "Automatic retry with exponential backoff"},
      {"label": "Fail fast", "description": "Return error immediately to caller"},
      {"label": "Circuit breaker", "description": "Stop attempts after repeated failures"}
    ]
  }
]
</parameter>
</invoke>
</function_calls>
```

## Question Categories

| Category | Example Question |
|----------|------------------|
| Approach | "Should we use X or Y for this?" |
| Priority | "Performance vs. simplicity for this feature?" |
| Pattern | "Follow existing pattern A or introduce new pattern B?" |
| Scope | "Include edge case X in MVP or defer?" |
