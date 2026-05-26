# Agent Handoff — Architecture

## Problem

When an AI agent's context window fills up, the user loses accumulated project knowledge.
Existing solutions (`.copilot-instructions.md`, skills files) are too rigid for evolving,
undocumented, or exploratory projects.

## Solution

A CLI tool (`handoff`) that lets agents **store** and **retrieve** knowledge packages — structured
blobs of project context that survive across sessions.

```
Session A (full context)          Session B (fresh context)
        │                                  │
        ▼                                  ▼
  handoff store ──► SQLite DB ◄── handoff retrieve
```

## Core Concepts

| Concept | Description |
|---------|-------------|
| **Knowledge Package** | A named, timestamped blob of markdown content with metadata |
| **TTL** | Time-to-live. Packages auto-expire. Default: 7 days |
| **Lazy GC** | Expired packages are purged on every DB operation |
| **Project Scope** | Optional grouping key to isolate packages per project |

## CLI Interface

```bash
# Store a knowledge package (content via stdin)
echo "..." | handoff store --name "auth-design" --ttl 7d --project "myapp" --tags "auth,api"
# Returns: package ID

# Retrieve by ID or name
handoff retrieve abc123
handoff retrieve --name "auth-design"
# Outputs content to stdout

# List available packages
handoff list
handoff list --project "myapp"

# Manual garbage collection
handoff gc
```

## Agent Workflow

1. User tells Agent A: "perform a knowledge transfer"
2. Agent A composes a markdown summary of relevant context
3. Agent A runs: `echo '<content>' | handoff store --name "project-x-state" --ttl 14d`
4. Agent A reports back the package ID to the user
5. User starts Session B and says: "retrieve knowledge package project-x-state"
6. Agent B runs: `handoff retrieve --name "project-x-state"`
7. Agent B reads the content and has full context

## Data Model

```sql
CREATE TABLE packages (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    summary    TEXT,
    content    TEXT NOT NULL,
    tags       TEXT,          -- JSON array: ["auth", "api"]
    project    TEXT,          -- optional project grouping
    created_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL
);
```

## Storage

- Location: `~/.handoff/handoff.db` (override via `HANDOFF_DB` env var)
- Engine: SQLite via `modernc.org/sqlite` (pure Go, no CGO)
- Single file, no server, no config

## TTL Strategy

- TTL is specified as a duration at creation time (e.g., `7d`, `30d`, `2h`)
- Converted to absolute `expires_at` timestamp
- Every DB operation runs `DELETE FROM packages WHERE expires_at < NOW()` first
- Explicit `handoff gc` command for manual cleanup

## Package Content Guidelines (for agents)

The tool doesn't enforce content structure, but agents should aim for:

```markdown
## Context
What the project is, what we're building, current goals.

## Key Decisions
- Decision: rationale
- Decision: rationale

## Current State
Where things stand right now. What's done, what's in progress.

## Technical Details
Stack, architecture patterns, file structure, key APIs.

## Open Issues / Next Steps
What remains to be done or resolved.
```

## Project Structure

```
├── main.go
├── cmd/
│   ├── root.go          # CLI root + global flags
│   ├── store.go         # `handoff store` command
│   ├── retrieve.go      # `handoff retrieve` command
│   ├── list.go          # `handoff list` command
│   └── gc.go            # `handoff gc` command
├── internal/
│   ├── db/
│   │   ├── db.go        # DB connection + lazy GC
│   │   └── schema.go    # Table creation / migrations
│   ├── model/
│   │   └── package.go   # Package struct + helpers
│   └── config/
│       └── config.go    # DB path resolution
└── go.mod
```

## Design Principles

1. **Zero config** — works out of the box with sensible defaults
2. **Single binary** — `go build` produces one executable, no runtime deps
3. **Cross-platform** — pure Go, no CGO, builds for linux/mac/windows
4. **Stdin/Stdout** — content flows through pipes, natural for agent terminal use
5. **Ephemeral by default** — TTL ensures nothing accumulates forever
6. **Simple over clever** — no plugin system, no network, no auth
