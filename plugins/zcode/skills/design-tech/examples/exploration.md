# Exploration Commands

Commands for exploring project context during technical design.

## Architecture & Decisions

```bash
# Check architecture constraints
cat docs/ARCHITECTURE.md

# Check existing decisions
cat docs/DECISIONS.md

# Check dependencies
cat package.json  # or go.mod, Cargo.toml, etc.
```

## Code Patterns

### Pattern: Find similar implementations
```bash
find <path> -name "*.<ext>" | xargs grep -l "<pattern>"
```

### Example: Find authentication-related code in TypeScript project
```bash
find src -name "*.ts" | xargs grep -l "authenticate\|auth"
```

## Recent Changes

```bash
git log --oneline -10
```
