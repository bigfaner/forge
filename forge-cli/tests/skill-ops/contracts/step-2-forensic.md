# Contract: skill-ops / Step 2: Forensic

## Outcome "forensic-search-sessions"
- Preconditions: "history.jsonl with recorded sessions"
- Input: `forge forensic search --last 5`
- Output: "exit code 0, session output contains sessionId field"
- State: "no state changes"
- Side-effect: none

## Outcome "forensic-extract-evidence"
- Preconditions: "valid session JSONL file path"
- Input: `forge forensic extract <path>`
- Output: "evidence summary with structured output"
- State: "no state changes"
- Side-effect: none

## Outcome "forensic-subagents-list"
- Preconditions: "session directory with subagent transcripts"
- Input: `forge forensic subagents <session-id>`
- Output: "list of subagent transcripts"
- State: "no state changes"
- Side-effect: none

## Outcome "forensic-extract-nonexistent"
- Preconditions: "no file at given path"
- Input: `forge forensic extract /nonexistent/path.jsonl`
- Output: "exit code 1, output mentions cannot/not found/no such"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- forge binary path consistent across all steps
- all commands use built binary, not system-installed
