# pkg/ Dependency Graph — Factual Baseline

> Auto-generated from `go list -json ./pkg/...` in forge-cli/ module.
> This document serves as the factual basis for all subsequent specification and package restructuring.

## 1. Complete Import Table

Each row shows one pkg/ subpackage and the other pkg/ subpackages it imports (standard library and third-party imports are excluded).

| Package | Internal pkg/ Imports | Import Count |
|---------|----------------------|-------------|
| `pkg/facttable` | (none) | 0 |
| `pkg/feature` | `pkg/git`, `pkg/index`, `pkg/types` | 3 |
| `pkg/forgeconfig` | `pkg/types` | 1 |
| `pkg/git` | (none) | 0 |
| `pkg/index` | (none) | 0 |
| `pkg/infocmd` | (none) | 0 |
| `pkg/just` | (none) | 0 |
| `pkg/lesson` | `pkg/infocmd` | 1 |
| `pkg/project` | (none) | 0 |
| `pkg/prompt` | `pkg/feature`, `pkg/forgeconfig`, `pkg/task` | 3 |
| `pkg/proposal` | `pkg/feature`, `pkg/infocmd` | 2 |
| `pkg/research` | `pkg/infocmd` | 1 |
| `pkg/serverprobe` | `pkg/feature`, `pkg/just` | 2 |
| `pkg/task` | `pkg/forgeconfig`, `pkg/index`, `pkg/infocmd`, `pkg/types` | 4 |
| `pkg/testrunner` | `pkg/just` | 1 |
| `pkg/types` | (none) | 0 |
| `pkg/version` | (none) | 0 |

**Total: 17 subpackages**

## 2. Three-Tier Classification

Classification rules:
- **Leaf**: zero internal forge-cli pkg/ dependencies
- **Infrastructure**: only depends on `pkg/types/` (pure data definitions)
- **Domain**: depends on other pkg/ subpackages beyond types

### Leaf (zero internal dependencies)

| Package | Notes |
|---------|-------|
| `pkg/facttable` | Standalone; no internal imports |
| `pkg/git` | Standalone; wraps git CLI operations |
| `pkg/index` | Standalone; index file management |
| `pkg/infocmd` | Standalone; info/command metadata |
| `pkg/just` | Standalone; just task runner integration |
| `pkg/project` | Standalone; project utilities |
| `pkg/types` | Zero imports (pure type definitions) |
| `pkg/version` | Zero imports (version constant) |

**Count: 8 packages**

### Infrastructure (only depends on pkg/types/)

| Package | Imports from pkg/ | Justification |
|---------|-------------------|---------------|
| `pkg/forgeconfig` | `pkg/types` | Config parsing; uses type definitions |

**Count: 1 package**

### Domain (depends on other pkg/ subpackages)

| Package | Imports from pkg/ | Classification Detail |
|---------|-------------------|----------------------|
| `pkg/feature` | `pkg/git`, `pkg/index`, `pkg/types` | Mixes leaf + types |
| `pkg/lesson` | `pkg/infocmd` | Single leaf dependency |
| `pkg/prompt` | `pkg/feature`, `pkg/forgeconfig`, `pkg/task` | Full domain: imports domain packages |
| `pkg/proposal` | `pkg/feature`, `pkg/infocmd` | Imports domain + leaf |
| `pkg/research` | `pkg/infocmd` | Single leaf dependency |
| `pkg/serverprobe` | `pkg/feature`, `pkg/just` | Imports domain + leaf |
| `pkg/task` | `pkg/forgeconfig`, `pkg/index`, `pkg/infocmd`, `pkg/types` | Mixes infrastructure + leaf + types |
| `pkg/testrunner` | `pkg/just` | Single leaf dependency |

**Count: 8 packages**

## 3. Fan-In Analysis (Imported-By Count)

How many other pkg/ packages import each package:

| Package | Imported By Count | Imported By |
|---------|:-:|-------------|
| `pkg/infocmd` | 4 | `pkg/lesson`, `pkg/proposal`, `pkg/research`, `pkg/task` |
| `pkg/feature` | 3 | `pkg/prompt`, `pkg/proposal`, `pkg/serverprobe` |
| `pkg/types` | 3 | `pkg/feature`, `pkg/forgeconfig`, `pkg/task` |
| `pkg/forgeconfig` | 2 | `pkg/prompt`, `pkg/task` |
| `pkg/index` | 2 | `pkg/feature`, `pkg/task` |
| `pkg/just` | 2 | `pkg/serverprobe`, `pkg/testrunner` |
| `pkg/git` | 1 | `pkg/feature` |
| `pkg/task` | 1 | `pkg/prompt` |

## 4. Horizontal Dependencies (Domain-to-Domain)

Horizontal dependencies exist when a domain package imports another domain package. These are the key coupling points to monitor during restructuring.

| Source | Target | Type |
|--------|--------|------|
| `pkg/prompt` | `pkg/feature` | domain -> domain |
| `pkg/prompt` | `pkg/task` | domain -> domain |
| `pkg/proposal` | `pkg/feature` | domain -> domain |
| `pkg/serverprobe` | `pkg/feature` | domain -> domain |
| `pkg/task` | `pkg/forgeconfig` | domain -> infrastructure |

**Key observation**: `pkg/feature` is the most heavily imported domain package (3 consumers), forming a central hub in the dependency graph. `pkg/prompt` has the deepest dependency chain, importing both `pkg/feature` and `pkg/task` (which itself imports 4 packages).

## 5. Bidirectional Coupling

No bidirectional coupling detected. For every pair (A, B) where A imports B, B does not import A. The dependency graph is a strict DAG.

## 6. Dependency Graph (Textual)

```
pkg/types            [leaf, fan-in: 3]
pkg/version          [leaf, fan-in: 0]
pkg/git              [leaf, fan-in: 1]
pkg/index            [leaf, fan-in: 2]
pkg/infocmd          [leaf, fan-in: 4]
pkg/just             [leaf, fan-in: 2]
pkg/project          [leaf, fan-in: 0]
pkg/facttable        [leaf, fan-in: 0]

pkg/forgeconfig      [infrastructure]
  └── pkg/types

pkg/feature          [domain]
  ├── pkg/git
  ├── pkg/index
  └── pkg/types

pkg/lesson           [domain]
  └── pkg/infocmd

pkg/research         [domain]
  └── pkg/infocmd

pkg/testrunner       [domain]
  └── pkg/just

pkg/proposal         [domain]
  ├── pkg/feature
  └── pkg/infocmd

pkg/serverprobe      [domain]
  ├── pkg/feature
  └── pkg/just

pkg/task             [domain]
  ├── pkg/forgeconfig
  ├── pkg/index
  ├── pkg/infocmd
  └── pkg/types

pkg/prompt           [domain, deepest chain]
  ├── pkg/feature
  │     ├── pkg/git
  │     ├── pkg/index
  │     └── pkg/types
  ├── pkg/forgeconfig
  │     └── pkg/types
  └── pkg/task
        ├── pkg/forgeconfig
        ├── pkg/index
        ├── pkg/infocmd
        └── pkg/types
```

## 7. Summary Statistics

| Metric | Value |
|--------|-------|
| Total pkg/ subpackages | 17 |
| Leaf packages | 8 (47%) |
| Infrastructure packages | 1 (6%) |
| Domain packages | 8 (47%) |
| Max dependency depth | 3 (pkg/prompt -> pkg/task -> pkg/forgeconfig -> pkg/types) |
| Horizontal dependencies | 4 |
| Bidirectional couplings | 0 |
| Highest fan-in | pkg/infocmd (4 importers) |
| Packages with zero importers | 5 (facttable, project, version, lesson, research) |

## 8. Risk Indicators for Restructuring

| Risk | Package | Detail |
|------|---------|--------|
| High fan-in | `pkg/infocmd` | Imported by 4 packages; changes have wide blast radius |
| High fan-in | `pkg/feature` | Imported by 3 domain packages; central hub |
| Deep chain | `pkg/prompt` | 3 levels deep; most complex dependency tree |
| Wide deps | `pkg/task` | Imports 4 internal packages; high coupling |
| Orphan | `pkg/facttable`, `pkg/project`, `pkg/version` | Zero internal importers; consider if these belong in pkg/ |
