# Commit Message Templates

## Interactive Mode

```
chore(specs): drift auto-fix -- 2 updated, 1 removed, 1 added

Updated:
  - BIZ-auth-001: align with renamed validateToken -> verifySession
  - TECH-api-003: reflect new rate limit threshold (100 -> 200)

Removed:
  - TECH-api-002: corresponding legacy proxy module removed

Added:
  - TECH-error-006: implicit error wrapping convention (user-confirmed)
```

## Non-Interactive Mode

```
chore(specs): [auto-specs] auto-integrate -- 3 added + drift auto-fix

Added:
  - BIZ-auth-002: session token validation rule
  - TECH-api-004: rate limiting convention
  - TECH-error-006: error wrapping pattern

Updated:
  - BIZ-auth-001: align with renamed validateToken -> verifySession

Overlap warnings:
  - docs/conventions/error-handling.md and docs/conventions/error-reporting.md share 66% domains (kept separate)

Review: git diff HEAD~1 | Revert: git revert HEAD
```
