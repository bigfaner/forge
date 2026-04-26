Upgrade the forge plugin version number. This bumps the version in both plugin manifest and marketplace entry, then commits.

Steps:

1. Read both version files to get the current version:
   - `plugins/forge/.claude-plugin/plugin.json`
   - `.claude-plugin/marketplace.json`

2. Ask the user for the new version number (suggest semver bump based on changes since last version: patch for fixes, minor for new features/skills, major for breaking changes).

3. Update the `version` field in both files:
   - `plugins/forge/.claude-plugin/plugin.json` line 3
   - `.claude-plugin/marketplace.json` line 14

4. Commit with message: `chore(forge): bump plugin version to <new-version>`
