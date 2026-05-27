# handoff

**When an AI agent's context window fills up, accumulated project knowledge is lost.** `handoff` lets agents store structured knowledge packages to a local SQLite database and retrieve them in the next session — no cloud, no server, no configuration required.

![Session A stores to handoff.db; Session B retrieves from handoff.db](https://raw.githubusercontent.com/Dborasik/handoff/main/assets/flow.png)

---

## The problem

Every AI coding agent — Claude Code, GitHub Copilot, Cursor, OpenAI Codex — works inside a finite context window. When that window fills up, the user must start a fresh session. Any decisions made, architecture discussed, progress tracked, or gotchas discovered during the session is gone.

`handoff` fixes this. Before the window closes, the agent stores a structured markdown summary to a local SQLite database. In the next session, a fresh agent retrieves that summary and picks up exactly where the previous one left off.

---

## Quick start

**Session A** — agent stores its context before the window fills:

```bash
echo "## Context
Building a REST API. Stack: Go + Chi + PostgreSQL.

## Current State
Auth complete. Working on task CRUD.

## Next Steps
1. Finish PATCH /tasks/:id
2. Add pagination to GET /tasks" | handoff store --name "api-state" --ttl 7d
```

Output:

```
a3f9c12e
```

**Session B** — fresh agent retrieves the context instantly:

```bash
handoff retrieve a3f9c12e
# or by name:
handoff retrieve --name "api-state"
```

That's it. The agent reads the stored markdown and resumes work with full context.

---

## Key properties

| Property | Detail |
|----------|--------|
| **Local-only** | All data lives in `~/.handoff/handoff.db` on your machine |
| **No daemon** | No background process, no server, no network calls |
| **Single binary** | Pure Go — one executable, no runtime dependencies |
| **Auto-expiry** | Packages expire by TTL; old context never accumulates forever |
| **Cross-platform** | macOS, Linux, Windows — amd64 and arm64 |

---

## Next steps

- [Install](install.md) — Homebrew, `go install`, or build from source
- [Commands](commands.md) — full CLI reference: `store`, `retrieve`, `list`, `gc`
- [Agent Setup](agents.md) — how to wire `handoff` into your agent's instructions
- [Workflow](workflow.md) — the recommended handoff workflow and package format
- [How It Works](internals.md) — storage, TTL, IDs, database schema
