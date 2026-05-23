Simplify and clean up code for the {{FEATURE_SLUG}} feature.

## Discovery Strategy
1. Run `git diff --name-only main...HEAD` to identify files changed by this feature
2. Focus cleanup on changed files only
3. The skill resolves scope: git diff > feature context > user-specified paths

Do NOT clean up files outside this feature's scope.
