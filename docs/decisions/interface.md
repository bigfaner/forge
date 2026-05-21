# Interface Decisions

| Date | Feature | Decision | Rationale | Source |
|------|---------|----------|-----------|--------|
| 2026-05-21 | cli-created-field-and-display | CLI 列表排序统一使用 frontmatter `created` 字段（YYYY-MM-DD）降序，mtime 作为 fallback | git clone/pull 会将文件 mtime 重置为 clone 时间而非原始创建时间，导致 mtime 排序不可靠。Proposal 已有 `created` 字段验证此方案可行。 | proposal: cli-created-field-and-display |
