# Research Dimensions

Two dimension sets — one per research mode. Core dimensions are always covered; optional dimensions are included only when the user selects them.

## Single-Tech Deep Dive

Use when researching one technology in depth (learning, understanding, evaluating).

| Dimension | Core/Optional | Description |
|---|---|---|
| Overview & Positioning | Core | What it is, what problem it solves, how it positions itself |
| Architecture & Core Concepts | Core | Internal design, key abstractions, mental model |
| Learning Curve | Core | Onboarding difficulty, prerequisite knowledge, time-to-productivity |
| Ecosystem & Community | Optional | Package count, activity level, enterprise adoption |
| Performance & Benchmarks | Optional | Benchmark results, resource consumption, scalability |
| Security | Optional | Known vulnerabilities, security model, track record |
| Version & Stability | Optional | Release cadence, breaking change policy, LTS support |

## Multi-Candidate Comparison

Use when comparing multiple technologies for a selection decision.

| Dimension | Core/Optional | Description |
|---|---|---|
| Overview & Positioning | Core | Each candidate's identity, target audience, philosophy |
| Ecosystem & Community | Core | Package ecosystem, community size, corporate backing |
| Performance & Benchmarks | Core | Head-to-head benchmarks, resource usage comparison |
| Developer Experience | Core | API design, documentation quality, debugging tools |
| Learning Curve | Core | Time-to-productivity, required knowledge, migration path |
| Security | Optional | Vulnerability history, security architecture |
| Migration Cost | Optional | Effort to adopt or switch, compatibility concerns |
| Version & Stability | Optional | Release stability, backward compatibility track record |

## Usage in AskUserQuestion

When presenting dimension selection (Phase 1, Round 2, Q4):

1. Load the relevant dimension set based on confirmed mode
2. Present core dimensions as pre-selected
3. Present optional dimensions as selectable additions
4. Allow user to deselect core dimensions if they explicitly don't want them
5. Accept custom dimensions the user proposes — these override the predefined list
