---
step: 3
title: Feature & Proposal Status Display
journey: feature-management
---

# Step 3: Feature & Proposal Status Display

## Given
- Proposals with various status values (Approved, Completed, draft)
- Features with manifest.md containing status field

## When
- `forge proposal` lists proposals with STATUS column
- `forge feature status <slug>` queries a feature's status

## Then
- Proposal list displays frontmatter status correctly (Approved, Completed)
- Feature status command shows manifest status in structured block format

## Contract Dimensions
- **Actor**: CLI user checking status of proposals and features
- **Input**: proposal.md frontmatter, feature manifest.md status field
- **Output**: STATUS column in table (proposal) or structured --- block (feature status)
- **Side Effects**: none (read-only)
- **Error Cases**: N/A
- **Invariants**: status values preserved exactly from frontmatter
