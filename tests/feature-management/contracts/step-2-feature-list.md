---
step: 2
title: Feature & Proposal List
journey: feature-management
---

# Step 2: Feature & Proposal List

## Given
- A forge project with proposals under docs/proposals/ and features under docs/features/
- Proposals have created frontmatter dates
- Features have manifest.md files with mtime

## When
- `forge proposal` is executed
- `forge feature list` is executed

## Then
- Proposals listed newest-first by created date (mtime fallback)
- Features listed newest-first by manifest.md mtime
- Missing manifest features sort to end
- Empty directories produce "no proposals found" / "no features found"

## Contract Dimensions
- **Actor**: CLI user listing proposals and features
- **Input**: docs/proposals/ and docs/features/ directories
- **Output**: CLI table with SLUG, STATUS columns (and others)
- **Side Effects**: none (read-only)
- **Error Cases**: empty directory -> informational message, exit 0
- **Invariants**: reverse chronological order; mtime fallback for missing created field
