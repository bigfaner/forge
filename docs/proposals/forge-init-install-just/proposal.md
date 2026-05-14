---
created: 2026-05-14
author: faner
status: Draft
---

# Proposal: forge init Auto-Install just

## Problem

`just` (casey/just) is a hard dependency of forge's quality gate pipeline, but users must install it manually before `forge init` can function. This creates a two-step bootstrap friction: install just yourself, then run `forge init`. If just is missing, every quality gate step silently fails or errors out.

### Evidence

- `forge init` currently appends `claude`/`claude-c` recipes to the justfile, but never verifies just is installed.
- The quality gate (`just compile → just fmt → just lint → just test`) fails with a confusing shell error if just is absent.
- `/init-justfile` skill is AI-driven and unstable — it depends on Claude context for template selection, placeholder substitution, and dry-run verification.

### Urgency

`just` is the single point of failure for forge's entire quality gate. New users who skip installing it get a broken experience with no clear guidance. Solving this in `forge init` — the natural onboarding entry point — eliminates the most common setup failure.

## Proposed Solution

Add a `ensureJust()` step to `forge init` that:

1. **Detects** whether `just` is already in PATH and meets the minimum version (>= 1.50.0, required for `[arg]` attribute support).
2. **Prompts** the user for confirmation if just is missing or too old.
3. **Installs** via system package manager (brew/cargo/scoop/choco) first, falling back to an embedded binary extracted to `~/.forge/bin/`.
4. **Reports** the result alongside existing init summary (CREATED/INSTALLED/SKIPPED/FAILED).

### Innovation Highlights

This is a straightforward "batteries included" approach — embedding a tool binary in the CLI for offline fallback is common in the industry (e.g., rustup embeds cargo, bun embeds itself). The dual strategy (package manager + embedded fallback) balances freshness with reliability.

## Requirements Analysis

### Key Scenarios

- **Happy path**: User runs `forge init` on a fresh machine → just detected as missing → user confirms → installed via brew/cargo/scoop → `just --version` succeeds → init continues.
- **Fallback path**: Package manager unavailable or fails → embedded binary extracted to `~/.forge/bin/` → PATH updated → init continues.
- **Already installed**: `just --version` returns >= 1.50.0 → step reports SKIPPED → init continues.
- **Outdated version**: `just --version` returns < 1.50.0 → prompt to upgrade → user declines → WARNING (non-blocking), init continues.
- **Skip requested**: User passes `--skip-just` → step skipped entirely → init continues.
- **CI environment**: (Out of scope for this iteration — future `--yes` flag).

### Non-Functional Requirements

- **Binary size**: Embedded just binary adds ~4 MB per platform. Compile-time embedding (build tags) ensures only the target platform's binary is included.
- **Performance**: Detection is a single `just --version` subprocess call. Installation is I/O-bound (package manager or file extraction).
- **Security**: Embedded binary is checksum-verified at build time. No runtime downloads from untrusted sources.

### Constraints & Dependencies

- just >= 1.50.0 required for `[arg]` attribute support used in e2e recipes.
- 6 platform binaries needed: linux/mac/windows × amd64/arm64.
- `~/.forge/bin/` must be added to PATH (or user shell profile) for the fallback to work.
- Build pipeline must support per-platform binary embedding via Go build tags or conditional `go:embed`.

## Alternatives & Industry Benchmarking

### Industry Solutions

- **rustup**: Embeds and manages Rust toolchain components. Handles PATH setup automatically.
- **fnm (Fast Node Manager)**: Single binary with platform-specific releases. Users download manually or via package managers.
- **Protocol Buffers compiler (protoc)**: Many tools embed it or download it at build time (e.g., `buf` embeds it).

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero development cost | Users hit cryptic errors; broken first experience | Rejected: core UX failure |
| Print install guide | fnm-style | Zero binary size increase | User still must act; friction remains | Rejected: doesn't solve the problem |
| Runtime download | buf-style | Always latest version; small binary | Requires network; download may fail; security concerns | Rejected: offline reliability matters |
| **System pkg + embedded fallback** | rustup-style | Works offline; always succeeds; user choice preserved | +4 MB binary size; version may lag | **Selected: best reliability/ergonomics trade-off** |

## Feasibility Assessment

### Technical Feasibility

Go supports `go:embed` for binary embedding and build tags for conditional compilation. The existing `forge init` architecture already uses the step pattern with `initAction` status reporting. Adding a new step is straightforward.

### Resource & Timeline

Single developer, estimated 3-5 tasks:
1. Download and embed just binaries for 6 platforms
2. Implement `ensureJust()` with detection and package manager dispatch
3. Implement embedded binary extraction fallback
4. Wire into `forge init` step sequence
5. Tests (unit + e2e)

### Dependency Readiness

- just releases are available on GitHub for all 6 target platforms.
- System package manager commands are well-documented and stable.
- No external API or service dependency.

## Scope

### In Scope

- `ensureJust()` step in `forge init` with detect → confirm → install flow
- System package manager support: brew (macOS), cargo (cross-platform), scoop/choco (Windows)
- Embedded just binary fallback: compile-time embed for current platform, extract to `~/.forge/bin/`
- Version check: >= 1.50.0, with upgrade prompt for outdated versions
- `--skip-just` CLI flag to skip the step
- 6 platform binaries: linux/mac/windows × amd64/arm64

### Out of Scope

- justfile generation logic (remains in `/init-justfile` skill)
- `just` auto-update mechanism
- CI/CD silent install mode (`--yes` flag)
- Installation of other tools (node, go, etc.)
- PATH persistence across shell sessions (document manual step if needed)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Embedded binary significantly inflates forge CLI size beyond ~4 MB | L | M | Compile-time per-platform embedding; monitor binary size in CI |
| `~/.forge/bin/` not in PATH after fallback install | M | H | Print clear instructions; consider appending to shell profile (.bashrc/.zshrc) |
| just 1.50.0 becomes incompatible with future justfile features | L | M | Pin minimum version; upgrade embedded binary with forge releases |
| System package manager installs outdated just (< 1.50.0) | M | M | Post-install version check; suggest embedded fallback if too old |
| Antivirus flags embedded binary on Windows | L | M | Code-sign the embedded binary; document the expected hash |

## Success Criteria

- [ ] `forge init` on a machine without just successfully installs it (either via package manager or fallback)
- [ ] `forge init --skip-just` skips just installation entirely
- [ ] `forge init` on a machine with just >= 1.50.0 reports SKIPPED and continues
- [ ] `forge init` on a machine with just < 1.50.0 warns about version and prompts upgrade
- [ ] forge CLI binary size increase is <= 5 MB per platform
- [ ] Works on all 6 target platforms (linux/mac/windows × amd64/arm64)

## Next Steps

- Proceed to `/write-prd` to formalize requirements
