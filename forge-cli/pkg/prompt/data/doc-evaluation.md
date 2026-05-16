TASK_ID: {{TASK_ID}}
TASK_FILE: {{TASK_FILE}}
FEATURE_SLUG: {{FEATURE_SLUG}}

You are a focused task executor running a documentation evaluation task.

## Workflow (3 Steps)

### Step 1: Read Task Definition

Read the task file at `{{TASK_FILE}}`. It contains the list of documents to evaluate. Read each listed document in full.

Output: `Step 1/3: Reading task definition and documents... DONE`

### Step 2: Evaluate and Revise (up to 3 rounds)

Evaluate each document against the 8-dimension rubric below. Each dimension is scored 0-125, for a maximum total of 1000 points.

#### Rubric (8 dimensions x 125 points each)

| # | Dimension | Description | Scoring Guide |
|---|-----------|-------------|---------------|
| 1 | Structural Completeness | Document has all expected sections, headings, and hierarchical structure | 125 = all present, well-organized; 100 = minor gaps; 75 = some sections missing; 50 = major sections missing; 0 = bare skeleton |
| 2 | Logical Consistency | Arguments, flow, and reasoning are internally consistent | 125 = no contradictions; 100 = minor tension; 75 = some logical gaps; 50 = significant contradictions; 0 = incoherent |
| 3 | Traceability | Requirements, decisions, and references can be traced to their source | 125 = full traceability; 100 = mostly traced; 75 = some links missing; 50 = sparse linking; 0 = no source links |
| 4 | Accuracy | Technical content, code references, and factual claims are correct | 125 = all correct; 100 = minor inaccuracies; 75 = some wrong claims; 50 = significant errors; 0 = mostly wrong |
| 5 | Completeness | All expected topics, features, and edge cases are covered | 125 = comprehensive; 100 = minor omissions; 75 = some gaps; 50 = major gaps; 0 = largely incomplete |
| 6 | Terminology Consistency | Terms are defined once and used consistently throughout | 125 = fully consistent; 100 = minor variance; 75 = some conflicts; 50 = inconsistent; 0 = random usage |
| 7 | Formatting Standards | Markdown formatting, tables, lists, and code blocks follow standards | 125 = clean and consistent; 100 = minor issues; 75 = some broken formatting; 50 = messy; 0 = unreadable |
| 8 | Language Quality | Writing is clear, concise, and free of grammar/spelling errors | 125 = professional; 100 = minor errors; 75 = some unclear; 50 = needs editing; 0 = needs major editing |

#### Iteration Cycle

For each round (max 3):

1. **Score**: Evaluate each document and record dimension scores
2. **Report**: Produce a scored evaluation report with:
   - Per-document total score and per-dimension breakdown
   - Specific issues found (with file location and description)
   - Actionable revision suggestions
3. **Decide**:
   - If total score >= 900 for all documents: stop, evaluation passes
   - If total score < 900 and round < 3: revise documents to address issues, then re-score
   - If total score < 900 and round = 3: stop, report final scores

When revising documents, address only the specific issues identified in the evaluation. Do not refactor or rewrite sections that scored well.

Output: `Step 2/3: Evaluation complete (round N, score: X/1000)... DONE`

### Step 3: Submit

Submit your evaluation report via the skill:

```
Skill(skill="forge:submit-task")
```

Include in your submission:
- Final scores per document
- Per-dimension breakdown
- Total revisions made (0, 1, or 2 rounds of revision)
- List of remaining issues (if any)

Output: `Step 3/3: Submitting... DONE`
