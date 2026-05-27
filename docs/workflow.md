# Workflow

This page covers the full handoff workflow — when to do it, how the two sessions coordinate, and how to structure the content of a knowledge package so the next agent can resume immediately.

---

## The handoff workflow

A handoff involves two sessions: one that stores context, and one that retrieves it.

### Session A — storing context

1. **The trigger.** The user notices the context is getting long and says something like: *"Do a knowledge transfer before we lose context."* The agent can also offer proactively when it notices its context filling up.

2. **The agent composes a summary.** The agent writes a structured markdown document covering the current project state — what was decided, what is done, what is in progress, what comes next, and any important warnings.

3. **The agent runs `handoff store`.**

    ```bash
    echo '<summary>' | handoff store \
      --name "project-state" \
      --summary "One-line description of the current state" \
      --ttl 14d \
      --project "myapp" \
      --tags "relevant,tags"
    ```

4. **The agent reports the ID.** The command prints an 8-character hex ID (e.g. `a3f9c12e`). The agent tells the user: *"Stored as `a3f9c12e`. In the next session, retrieve it with `handoff retrieve a3f9c12e` or `handoff retrieve --name project-state`."*

### Session B — retrieving context

1. **The user starts a fresh session.** The fresh agent has no prior context.

2. **The agent checks for packages.** With [always-on instructions](agents.md) in place, the agent automatically runs `handoff list` at session start to see what is available:

    ```bash
    handoff list
    handoff list --project myapp
    ```

3. **The agent retrieves the package.**

    ```bash
    handoff retrieve a3f9c12e
    # or by name:
    handoff retrieve --name "project-state"
    ```

4. **The agent reads the content and resumes.** It now has full context and can continue work without any explanation from the user.

---

## Knowledge package format

`handoff` does not enforce any structure on the content — it stores whatever markdown text you pipe in. However, a consistent structure makes the stored context much more useful to the agent reading it in the next session.

The recommended format:

```markdown
## Context
What this project is and what we are currently building. Include the tech stack,
the overall goal, and any background the next agent needs to understand the project.

## Key Decisions
- Decision made: rationale and any alternatives that were considered
- Decision made: rationale and any alternatives that were considered

## Current State
What is fully complete. What is currently in progress. What is blocked and why.

## Next Steps
1. First concrete action to take
2. Second concrete action to take
3. Third concrete action to take

## Warnings / Gotchas
- Known issues or fragile areas
- Things that look wrong but are intentional
- Setup steps that are easy to miss

## Files of Note
- `path/to/file.go` — what it does and why it matters
- `path/to/other.ts` — what it does and why it matters
```

### Why each section matters

**Context** — A fresh agent knows nothing. Don't assume it has seen the README or explored the codebase. Describe the project as if briefing someone for the first time.

**Key Decisions** — The most overlooked section. Decisions made during the session (technology choices, API design, trade-offs) exist only in the conversation. Without this section, the next agent may re-examine or reverse them.

**Current State** — Distinguishes what is production-ready from what is half-done. The next agent needs to know where the work boundary is.

**Next Steps** — Ordered, specific actions. Not vague goals but concrete steps: "Add input validation to `PATCH /tasks/:id` using `go-playground/validator`", not "fix tasks endpoint".

**Warnings / Gotchas** — Traps that cost time to rediscover. Include anything that is non-obvious, environment-specific, or that caused confusion during the session.

**Files of Note** — The next agent will likely explore the codebase. Point it directly to the relevant files so it does not have to infer them from the structure.

---

## Full example

This is what a complete, real-world knowledge package looks like:

```bash
handoff store \
  --name "todo-api-state" \
  --summary "Task CRUD in progress, auth complete, tests pending" \
  --ttl 14d \
  --project "todo-api" \
  --tags "api,crud,auth,go" << 'EOF'
## Context
Building a REST API for a multi-user todo application.
Stack: Go + Chi router + PostgreSQL.
Goal: full CRUD for tasks with JWT authentication.
Repo: github.com/example/todo-api

## Key Decisions
- PostgreSQL over SQLite: needed for concurrent multi-user access
- JWT auth: 15-min access tokens + 7-day refresh tokens stored in HttpOnly cookies
- Chi router over Gin: lighter footprint, closer to standard library
- No ORM: using raw database/sql with pgx driver for clarity and control

## Current State
Complete:
- User registration and login
- JWT middleware (applied per-route, not globally)
- GET /tasks and POST /tasks

In progress:
- PATCH /tasks/:id — handler exists, input validation not yet written

Not started:
- DELETE /tasks/:id
- Pagination on GET /tasks
- Integration tests

Blocked: nothing currently blocked.

## Next Steps
1. Finish PATCH /tasks/:id — add validation using go-playground/validator
2. Implement DELETE /tasks/:id
3. Write integration tests using testcontainers-go
4. Add cursor-based pagination to GET /tasks

## Warnings / Gotchas
- DB migrations live in /migrations — always run `make migrate` before testing
- JWT_SECRET env var must be set or the server panics on startup (see main.go:42)
- The test database runs on port 5433, not 5432, to avoid conflicts with local Postgres
- The refresh token rotation logic is in internal/auth/refresh.go, not middleware.go

## Files of Note
- `internal/auth/middleware.go` — JWT validation middleware, applied per-route
- `internal/auth/refresh.go` — refresh token rotation and revocation logic
- `internal/task/handler.go` — all task HTTP handlers
- `migrations/` — SQL migration files, applied in ascending filename order
EOF
```

The agent in Session B retrieves this, reads it, and has everything it needs to continue working on `PATCH /tasks/:id` without asking any questions.
