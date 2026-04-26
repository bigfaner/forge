# Decision Logging — Shared Archiving Protocol

This file defines the shared decision archiving logic used by the `tech-design` skill and the `/record-decision` command. Read this file and follow the steps described below when archiving decisions.

---

## 1. Type Mapping

| Number | Type Name          | File Path                  |
|--------|--------------------|----------------------------|
| 1      | Architecture       | architecture.md            |
| 2      | Interface          | interface.md               |
| 3      | Data Model         | data-model.md              |
| 4      | Dependencies       | dependencies.md            |
| 5      | Error Handling     | error-handling.md          |
| 6      | Testing            | testing.md                 |
| 7      | Security           | security.md                |
| 8      | Local Dev & Deployment | local-dev-deployment.md |

All type files live under `docs/decisions/`.

---

## 2. tech-design Archiving Steps

Triggered after the user approves the tech-design document.

### Step 2.1 — Check for candidates

Scan the approved tech-design document for entries marked as key decisions. If none exist, skip to Step 2.5.

### Step 2.2 — Display candidate list

Show the numbered list of key decisions with their type in parentheses:

```
The following decisions are marked as key decisions and recommended for archiving:

  [1] Adopt event-driven architecture (Architecture)
  [2] Use SQLite as local cache storage (Data Model)
  [3] Choose Vitest over Jest as test framework (Dependencies)

Enter numbers to archive (comma-separated), or all / none:
```

### Step 2.3 — Handle user input

- `none` → skip to Step 2.5
- `all` → archive every candidate
- comma-separated numbers (e.g. `1,3`) → archive only those entries
- `edit:<number>` → enter the edit sub-flow (see Section 4) for that entry, then re-display the prompt

Invalid input (number not in candidate list): re-prompt with "Number X is not in the candidate list. Please re-enter."

### Step 2.4 — Write and update

For each selected entry:
1. Append a decision row to `docs/decisions/<type>.md` (see Section 6 for row format).
2. Update `docs/decisions/manifest.md` (see Section 7).

### Step 2.5 — Skip logic

If no key decisions exist in the tech-design document, silently skip the archiving step and proceed with the rest of the tech-design flow.

---

## 3. record-decision 4-Round Interaction

### Round 1 — Decision type

Display the type list and ask for a number:

```
Select decision type:

  1. Architecture
  2. Interface
  3. Data Model
  4. Dependencies
  5. Error Handling
  6. Testing
  7. Security
  8. Local Dev & Deployment

Enter number (1-8):
```

If the user enters a value outside 1-8, re-prompt: "Please enter a number between 1 and 8."

### Round 2 — Decision description

```
Enter decision description (one sentence, max 80 characters):
```

### Round 3 — Rationale

```
Enter decision rationale (one sentence, max 80 characters):
```

### Round 4 — Associated feature

```
Enter associated feature slug (e.g. feat-log-decisions), or press Enter to skip:
```

If skipped, set Feature to `-`.

### Auto-filled fields

- `Date`: today's date in ISO 8601 format (YYYY-MM-DD)
- `Source`: `<feature-slug>/tech-design.md` if invoked from tech-design flow; `manual` if invoked directly via `/record-decision`

After Round 4, write the entry (Section 6) and update the manifest (Section 7).

---

## 4. edit Sub-flow

Triggered when the user inputs `edit:<number>` during the tech-design candidate selection prompt.

### Steps

1. Validate that `<number>` exists in the current candidate list. If not, re-prompt: "Number X is not in the candidate list. Please re-enter."
2. Display the current Decision and Rationale for that entry.
3. Ask: "Enter new Decision (press Enter to keep current):"
4. Ask: "Enter new Rationale (press Enter to keep current):"
5. Update the in-memory candidate entry with the new values.
6. Return to the candidate selection prompt (Step 2.3).

---

## 5. Skip Logic

- **tech-design flow**: if the approved tech-design document contains zero key decisions, skip the entire archiving step without prompting the user.
- **record-decision flow**: no skip logic; the skill always runs the 4-round interaction.

---

## 6. Decision Entry Row Format

Append to the end of `docs/decisions/<type>.md`:

```
| YYYY-MM-DD | <feature-slug> | <Decision, one sentence> | <Rationale, one sentence> | <feature-slug>/tech-design.md §<Section> |
```

Field constraints:
- `Date`: ISO 8601 (YYYY-MM-DD)
- `Feature`: feature slug, e.g. `feat-log-decisions`; use `-` if unknown
- `Decision`: single sentence, max 80 characters
- `Rationale`: single sentence, max 80 characters
- `Source`: `<feature-slug>/<file>.md §<Section>` or `manual`

---

## 7. Manifest Update Protocol

Target file: `docs/decisions/manifest.md`

### Operation A — Categories table

Find the row matching the decision type. Increment the `Decisions` count by 1. Set `Last Updated` to today's date (YYYY-MM-DD).

### Operation B — Recent Decisions table

Insert a new row immediately below the table header (newest first). Keep a maximum of 10 rows; remove the oldest row if the count exceeds 10.

Row format:

```
| YYYY-MM-DD | <feature-slug> | <Type Name> | <Decision, one sentence> | <source> |
```

---

## 8. Error Handling

| Scenario | Handling |
|----------|----------|
| `docs/decisions/` directory does not exist | Auto-create the directory plus all 8 type files and `manifest.md` from their initial templates before archiving |
| `manifest.md` is missing | Rebuild it from the manifest template before archiving |
| Invalid type number in record-decision Round 1 | Re-prompt: "Please enter a number between 1 and 8." |
| `edit:<number>` references a non-existent candidate number | Re-prompt: "Number X is not in the candidate list. Please re-enter." |
| Type file header row is missing (file corrupted or empty) | Prepend the standard header before appending the new row: `# <Type Name> Decisions\n\n| Date | Feature | Decision | Rationale | Source |\n|------|---------|----------|-----------|--------|` |

---

## 9. Type File Initial State (reference)

Each type file should have this structure when first created:

```markdown
# <Type Name> Decisions

| Date | Feature | Decision | Rationale | Source |
|------|---------|----------|-----------|--------|
```

Replace `<Type Name>` with the name from the type mapping table (e.g. `Architecture`, `Interface`, etc.).
