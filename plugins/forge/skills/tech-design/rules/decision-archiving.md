# Decision Archiving Rules

Rules for archiving key decisions from the tech-design document into `docs/decisions/`.

## Step 7.1 — Check for candidates

Scan the approved tech-design document for entries marked as key decisions. If none exist, skip to Step 7.5.

## Step 7.2 — Display candidate list

Show the numbered list of key decisions with their type in parentheses:

```
The following decisions are marked as key decisions and recommended for archiving:

  [1] Adopt event-driven architecture (Architecture)
  [2] Use SQLite as local cache storage (Data Model)
  [3] Choose Vitest over Jest as test framework (Dependencies)

Enter numbers to archive (comma-separated), or all / none:
```

## Step 7.3 — Handle user input

- `none` → skip to Step 7.5
- `all` → archive every candidate
- comma-separated numbers (e.g. `1,3`) → archive only those entries
- `edit:<number>` → enter the edit sub-flow for that entry, then re-display the prompt

Invalid input (number not in candidate list): re-prompt with "Number X is not in the candidate list. Please re-enter."

## Step 7.4 — Write and update

For each selected entry:
1. Append a decision row to `docs/decisions/<type>.md` (see decision entry row format below).
2. Update `docs/decisions/manifest.md` (see manifest update protocol below).

### Decision Entry Row Format

Append to the end of `docs/decisions/<type>.md`:

```
| YYYY-MM-DD | <feature-slug> | <Decision, one sentence> | <Rationale, one sentence> | <feature-slug>/design/tech-design.md §<Section> |
```

Field constraints:
- `Date`: ISO 8601 (YYYY-MM-DD)
- `Feature`: feature slug, e.g. `feat-log-decisions`; use `-` if unknown
- `Decision`: single sentence, max 80 characters
- `Rationale`: single sentence, max 80 characters
- `Source`: `<feature-slug>/<file>.md §<Section>` or `manual`

Use `templates/decision-entry.md` for the template format.

### Manifest Update Protocol

Target file: `docs/decisions/manifest.md`

**Operation A — Categories table**

Find the row matching the decision type. Increment the `Decisions` count by 1. Set `Last Updated` to today's date (YYYY-MM-DD).

**Operation B — Recent Decisions table**

Insert a new row immediately below the table header (newest first). Keep a maximum of 10 rows; remove the oldest row if the count exceeds 10.

Row format:

```
| YYYY-MM-DD | <feature-slug> | <Type Name> | <Decision, one sentence> | <source> |
```

## Step 7.5 — Skip logic

If no key decisions exist in the tech-design document, silently skip the archiving step and proceed with the rest of the tech-design flow.

## Edit Sub-flow

Triggered when the user inputs `edit:<number>` during the candidate selection prompt.

1. Validate that `<number>` exists in the current candidate list. If not, re-prompt: "Number X is not in the candidate list. Please re-enter."
2. Display the current Decision and Rationale for that entry.
3. Ask: "Enter new Decision (press Enter to keep current):"
4. Ask: "Enter new Rationale (press Enter to keep current):"
5. Update the in-memory candidate entry with the new values.
6. Return to the candidate selection prompt (Step 7.3).

See `examples/ask-question.md` for question formatting and `examples/exploration.md` for context exploration commands.
