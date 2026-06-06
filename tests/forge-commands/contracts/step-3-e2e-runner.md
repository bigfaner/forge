# Contract: forge-commands / Step 3: E2E Runner

## Outcome "e2e-run-nonexistent-feature"
- Preconditions: "forge CLI binary built from current source"
- Input: `forge e2e run --feature nonexistent-feature`
- Output: "exit code 1, output mentions feature not found"
- State: "no state changes"
- Side-effect: none

## Outcome "e2e-run-with-profile"
- Preconditions: ".forge/config.yaml with valid profile and feature test data"
- Input: `forge e2e run`
- Output: "test suite execution results"
- State: "no state changes"
- Side-effect: none

## Outcome "e2e-run-no-profile"
- Preconditions: ".forge/config.yaml with no profile field"
- Input: `forge e2e run`
- Output: "error indicating no profile configured"
- State: "no state changes"
- Side-effect: none

## Outcome "e2e-run-unknown-profile"
- Preconditions: ".forge/config.yaml with unknown profile value"
- Input: `forge e2e run`
- Output: "error listing available profiles"
- State: "no state changes"
- Side-effect: none

## Outcome "forge-init-without-just"
- Preconditions: "clean temp directory, just may or may not be in PATH"
- Input: `forge init --project-root <tmpdir>`
- Output: "init summary includes 'just installation' step"
- State: "project artifacts created (.forge/, CLAUDE.md, .gitignore, justfile)"
- Side-effect: filesystem modifications

## Outcome "forge-init-skip-just"
- Preconditions: "clean temp directory"
- Input: `forge init --project-root <tmpdir> --skip-just`
- Output: "init summary shows 'just installation' as SKIPPED with 'skipped via --skip-just flag'"
- State: "all non-just artifacts created"
- Side-effect: filesystem modifications

## Outcome "forge-init-custom-project-root"
- Preconditions: "clean temp directory"
- Input: `forge init --project-root <tmpdir> --skip-just`
- Output: "all artifacts (.forge/, CLAUDE.md, .gitignore, justfile) created in custom directory"
- State: "project initialized in specified root"
- Side-effect: filesystem modifications

## Outcome "ensure-result-mapping"
- Preconditions: "none"
- Input: "Verify EnsureResult constants"
- Output: "INSTALLED, SKIPPED, FAILED status string representations correct"
- State: "no state changes"
- Side-effect: none

## Outcome "detect-just"
- Preconditions: "just installed in PATH"
- Input: "just.DetectJust()"
- Output: "found=true, non-empty path and version string matching semver pattern"
- State: "no state changes"
- Side-effect: none

## Outcome "parse-just-version"
- Preconditions: "none"
- Input: "just.ParseJustVersion with valid/invalid inputs"
- Output: "correct version extraction or error"
- State: "no state changes"
- Side-effect: none

## Outcome "is-minimum-version"
- Preconditions: "none"
- Input: "just.IsMinimumVersion with various version pairs"
- Output: "correct comparison results"
- State: "no state changes"
- Side-effect: none

## Outcome "embedded-binary"
- Preconditions: "none"
- Input: "just.EmbeddedBinaryFunc / ExtractEmbeddedBinaryFunc"
- Output: "valid extraction or proper failure status"
- State: "possible file creation in ~/.forge/bin/"
- Side-effect: filesystem modifications for extraction

## Journey Invariants
- forge binary path consistent across all steps
- all commands use built binary, not system-installed
- temp directories cleaned up after each test
