# Contract: test-generation / Step 2: Gen Scripts

## Outcome "ts-morph-devDependency"
- Preconditions: "gen-test-scripts templates package.json exists"
- Input: "Read package.json from gen-test-scripts templates"
- Output: "devDependencies contains ts-morph with valid semver range"
- State: "no state changes"
- Side-effect: none

## Outcome "skill-md-contains-step-45"
- Preconditions: "gen-test-scripts SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "SKILL.md contains Step 4.5 heading and validate-specs references"
- State: "no state changes"
- Side-effect: none

## Outcome "step-actionability-gate"
- Preconditions: "gen-test-scripts SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "SKILL.md has Prerequisites section with Step Actionability gate aborting when score < 20"
- State: "no state changes"
- Side-effect: none

## Outcome "type-filter-argument"
- Preconditions: "gen-test-scripts SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "SKILL.md documents --type filter with capability, invalid type handling, and type-only processing"
- State: "no state changes"
- Side-effect: none

## Outcome "type-filter-step-skip"
- Preconditions: "gen-test-scripts SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "SKILL.md describes shared infra always runs, Fact Table and locator skip for non-matching types"
- State: "no state changes"
- Side-effect: none

## Outcome "element-field-required"
- Preconditions: "gen-contracts SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "SKILL.md states Element is required, defines Element field, defines sitemap-missing fallback"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- gen-test-scripts SKILL.md is the authoritative script generation skill
- gen-contracts SKILL.md is the authoritative contract generation skill
- Profile capabilities define valid --type filter values
