# Forge

> **Spec Driven Development** toolkit — turning Claude Code from "chat" into "engineering"

[![Version](https://img.shields.io/badge/Version-3.0.0-blue.svg)](https://github.com/bigfaner/forge)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

If you've been coding with AI, you know the pain:

- Two hours in, the AI misunderstood the requirement — everything it wrote is useless
- "Add a button" somehow turned into a full module rewrite three hours later
- A bug you fixed yesterday came back today because the AI simply doesn't remember
- Every new session means re-teaching the AI your project conventions from scratch

**Forge** is a **Spec Driven Development (SDD)** toolkit built for Claude Code, powered by **Harness Engineering** — turning free-form AI conversations into controlled engineering pipelines. It replaces ad-hoc prompting with a disciplined pipeline: `brainstorm -> PRD -> design -> tasks -> auto-execute`. No more guessing. Just engineering-grade delivery.

---

## How Forge Compares

| | Forge | Superpowers | Spec Kit | OpenSpec |
|------|:-----:|:-----------:|:--------:|:--------:|
| Structured Workflow | ✓ | ✓ | ✓ | ✓ |
| Quality Gates | ✓ | ✗ | ✗ | ✗ |
| Persistent Context (manifest + worktree) | ✓ | ✗ | ✗ | ✗ |
| Knowledge Capture (`/learn`) | ✓ | ✗ | ✗ | ✗ |
| Cross-Session Continuity | ✓ | ✗ | ✗ | ✗ |
| Agent Orchestration (`/run-tasks`) | ✓ | ✓ | ✗ | ✗ |
| Multi-Agent / Subagent | ✗ | ✓ | ✗ | ✗ |
| Cross-IDE / Cross-Agent Platform | ✗ | ✓ | ✓ | ✓ |

> Data source: as of 2026-06, verified against each project's GitHub README and documentation.

---

## Core Features

### Structured Pipeline

Two modes for different scales of work:

| Mode | Flow | Best For |
|------|------|----------|
| Full Mode | brainstorm → PRD → tech design → task breakdown → auto-execute | Complex features (>10 tasks) |
| Quick Mode | brainstorm → direct task generation → auto-execute | Small features (1–10 tasks) |

### Adversarial Evaluation System

8 specialized evaluators (PRD, tech design, UI design, Proposal, Journey, Contract, etc.) using **expert personas + 1000+ point rubrics** for multi-round iterative revision. Supports cross-document consistency checks (PRD → design → task alignment), ensuring document quality from the source.

### Autonomous Task Execution Engine

`/run-tasks` continuously dispatches tasks to subagents, each strictly following **TDD protocol** (RED → GREEN → REFACTOR). Features a **dynamic fix chain**: failed tasks automatically create fix tasks, which restore the original task upon completion. 7-state task machine (pending → in_progress → completed / blocked / suspended, etc.) with Quality Gate (`compile → fmt → lint → test`) enforcement at every step.

### Journey-Contract Test Model

Extract **Journeys** (user flows + risk ratings) from PRD user stories, generate six-dimension **Contracts** (behavioral specs) from Journeys, and produce executable **Surface-aware test scripts** (web / api / cli / tui / mobile) from Contracts.

### Persistent Context

`manifest.md` tracks the full lifecycle of a feature. `forge worktree` isolates parallel development environments. `index.json` records task dependencies and progress. No more AI amnesia — every session picks up exactly where the last one left off.

### Knowledge Capture

`/learn` distills design decisions, lessons learned, and technical conventions into reusable knowledge artifacts. `/consolidate-specs` automatically detects and fixes spec-code drift. New sessions and new contributors build on accumulated experience instead of starting from zero every time.

---

## Engineering Scale

**21 Skills** · **16 Slash Commands** · **1 Subagent** · **20 Task Types** · **Go CLI** (19 commands) + **Claude Code Plugin** · Hooks system for automatic session-level context injection and cleanup

---

## 5-Minute Quick Start

```bash
# Initialize your project
forge init

# Quick mode — shorter pipeline, skips PRD/design/eval
/quick

# Full pipeline — complete brainstorm -> PRD -> design -> tasks -> execute flow
/brainstorm -> /write-prd -> /tech-design -> /breakdown-tasks -> /run-tasks
```

---

## Installation

### Prerequisites

- [Claude Code](https://docs.anthropic.com/en/docs/claude-code) CLI
- curl (pre-installed on macOS/Linux, included in Windows 10+)

### Install

**macOS / Linux:**

```bash
# 1. Install forge CLI
curl -fsSL https://github.com/bigfaner/forge/releases/latest/download/install.sh | bash

# 2. Install forge Plugin (CLI binary + Plugin in one step)
forge upgrade

# 3. Initialize in your project
cd my-project && forge init
```

**Windows (PowerShell):**

```powershell
# 1. Install forge CLI
irm https://github.com/bigfaner/forge/releases/latest/download/install.ps1 | iex

# 2. Install forge Plugin (CLI binary + Plugin in one step)
forge upgrade

# 3. Initialize in your project
cd my-project; forge init
```

Building from source: `git clone` -> `cd forge-cli && bash scripts/install-local.sh` -> `forge upgrade`

---

## Learn More

- [Architecture Overview](docs/user-guide/architecture-overview.md) — Plugin system, four core components, data flow and state management
- [Usage Guide](docs/user-guide/usage-guide.md) — Full Mode / Quick Mode end-to-end walkthroughs, single-command scenarios, troubleshooting
- [Project Initialization](docs/user-guide/initialization.md) — `forge init` walkthrough, config field reference, Surface detection
- [Environment Setup](docs/user-guide/environment-setup.md) — Setting up your Forge development environment from scratch

---

## Contributing

```bash
git clone git@github.com:bigfaner/forge.git && cd forge
cd forge-cli && go mod download
go test -race -cover ./...
```

Commits follow [Conventional Commits](https://www.conventionalcommits.org/).

## License

[MIT](LICENSE)
