Release the forge-cli binary. This bumps the version in `forge-cli/scripts/version.txt`, commits, tags, and pushes — GitHub Actions will auto-trigger from the tag to build and publish the release.

Steps:

1. Read the current version from `forge-cli/scripts/version.txt`.

2. Ask the user for the new version number (suggest semver bump based on changes since last version: patch for fixes, minor for new features, major for breaking changes). Note: CLI version numbers are independent from Plugin version numbers — they follow their own semver schedule starting at 5.x.x.

3. Update `forge-cli/scripts/version.txt` with the new version number.

4. Execute the following git operations in order:
   - `git add forge-cli/scripts/version.txt`
   - `git commit -m "chore(forge-cli): bump version to <new-version>"`
   - `git tag forge-cli/v<new-version>`
   - `git push origin HEAD forge-cli/v<new-version>`

5. Inform the user that GitHub Actions will automatically detect the `forge-cli/v<new-version>` tag push and run the release pipeline to build binaries and publish them to GitHub Releases.
