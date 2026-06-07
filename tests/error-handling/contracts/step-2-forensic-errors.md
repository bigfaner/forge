# Contract: error-handling / Step 2: Forensic Errors

## Outcome "forensic-search-no-results"
- Preconditions: "keyword that matches no session transcripts"
- Input: `forge forensic search --keyword <unique-keyword> --last 1`
- Output: "exit code 0, output is '[]' (empty JSON array)"
- State: "no state changes"
- Side-effect: none

## Outcome "forensic-search-missing-records"
- Preconditions: "history.jsonl file missing or inaccessible"
- Input: `forge forensic search`
- Output: "appropriate error indicating records directory not found"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- forensic search always returns valid JSON
- exit code 0 for empty results, non-zero for infrastructure errors
