# Exploration Commands

Explore project context during technical design using dedicated tools.

## Architecture & Decisions

```text
Read docs/ARCHITECTURE.md     → layer constraints
Read docs/decisions/          → existing decisions (category-based directory)
Read package.json / go.mod    → current dependencies
```

## Code Patterns

### Pattern: Find similar implementations

Use `Grep` tool with pattern and path to find related code.

### Example: Find authentication-related code in TypeScript project

Use `Grep` tool with pattern `authenticate|auth` and glob `*.ts` in `src/`.

## Recent Changes

```bash
git log --oneline -10
```
