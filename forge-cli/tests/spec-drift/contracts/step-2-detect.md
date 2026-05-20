# Contract: spec-drift / Step 2: Detect

## Outcome "breakdown-tasks-references-drift"
- Preconditions: "breakdown-tasks SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "SKILL.md mentions T-specs-consolidate with drift detection reference"
- State: "no state changes"
- Side-effect: none

## Outcome "consolidate-specs-steps-9-11"
- Preconditions: "consolidate-specs SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "SKILL.md contains Steps 9-11 with step ordering verified"
- State: "no state changes"
- Side-effect: none

## Outcome "step-9-three-way-classification"
- Preconditions: "consolidate-specs SKILL.md exists"
- Input: "Read Step 9 section"
- Output: "Step 9 defines current/drifted/orphaned classification and reads business-rules + conventions"
- State: "no state changes"
- Side-effect: none

## Outcome "step-10-id-preservation"
- Preconditions: "consolidate-specs SKILL.md exists"
- Input: "Read Step 10 section"
- Output: "Step 10 instructs preserving project-global IDs, removing orphaned rules, detecting new rules"
- State: "no state changes"
- Side-effect: none

## Outcome "hard-gate-drift-exception"
- Preconditions: "consolidate-specs SKILL.md exists"
- Input: "Read HARD-GATE section"
- Output: "HARD-GATE includes exception clause for drift detected in Step 9"
- State: "no state changes"
- Side-effect: none

## Outcome "drift-only-mode"
- Preconditions: "consolidate-specs SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "SKILL.md documents drift-only mode running Steps 9-11"
- State: "no state changes"
- Side-effect: none

## Outcome "guide-references-drift"
- Preconditions: "guide.md exists"
- Input: "Read guide.md content"
- Output: "guide.md references T-quick-doc-drift and drift detection"
- State: "no state changes"
- Side-effect: none

## Outcome "guide-specs-drift-verification"
- Preconditions: "guide.md exists"
- Input: "Read guide.md content"
- Output: "guide.md mentions business-rules directory and drift verification with conventions"
- State: "no state changes"
- Side-effect: none

## Outcome "vocabulary-generation-step"
- Preconditions: "consolidate-specs SKILL.md exists"
- Input: "Read Step 12 section"
- Output: "Step 12 generates vocabulary scanning 4 knowledge directories with base 8 categories"
- State: "no state changes"
- Side-effect: none

## Outcome "vocabulary-output-structure"
- Preconditions: "consolidate-specs SKILL.md exists"
- Input: "Read Step 12 section"
- Output: "Step 12 output includes Types, Domains, Count; marks as AUTO-GENERATED; references /learn"
- State: "no state changes"
- Side-effect: none

## Outcome "vocabulary-step-ordering"
- Preconditions: "consolidate-specs SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "Steps 1-26 all present; Step 12 after Step 11 and before Step 13"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- consolidate-specs SKILL.md is the authoritative drift detection workflow
- Steps 9-11 implement drift detection, Step 10 handles ID preservation
- HARD-GATE allows drift modification as exception
- Step 12 generates vocabulary from 4 knowledge directories
