# How It Works

This page describes the internals of `handoff`: how data is stored, how package IDs are generated, how TTL and garbage collection work, and the design principles behind the tool.

---

## Storage

All data is stored in a single SQLite file at `~/.handoff/handoff.db`. The directory is created automatically (with permissions `0755`) the first time any `handoff` command is run — no setup step required.

`handoff` uses [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite), a pure-Go SQLite driver. There is no CGO dependency and no requirement to have SQLite installed on the system. The binary carries everything it needs.

---

## Database schema

```sql
CREATE TABLE IF NOT EXISTS packages (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    summary    TEXT DEFAULT '',
    content    TEXT NOT NULL,
    tags       TEXT DEFAULT '[]',
    project    TEXT DEFAULT '',
    created_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_packages_name       ON packages(name);
CREATE INDEX IF NOT EXISTS idx_packages_project    ON packages(project);
CREATE INDEX IF NOT EXISTS idx_packages_expires_at ON packages(expires_at);
```

Three indexes are maintained:

| Index | Used by |
|-------|---------|
| `idx_packages_name` | `retrieve --name` lookups |
| `idx_packages_project` | `list --project` filtering |
| `idx_packages_expires_at` | Garbage collection (`DELETE WHERE expires_at < NOW()`) |

**Tags** are stored as a JSON array string (e.g. `["auth","api"]`), not as a separate table. This keeps the schema simple and avoids joins.

**Timestamps** (`created_at`, `expires_at`) are stored as UTC datetimes. When reading them back, `handoff` tries multiple datetime formats to handle the various string representations SQLite may return.

---

## Package IDs

Package IDs are generated using [`crypto/rand`](https://pkg.go.dev/crypto/rand): 4 random bytes are read from the OS cryptographic random source and hex-encoded, producing an 8-character lowercase hex string such as `a3f9c12e`.

This gives 2³² (about 4.3 billion) possible IDs. IDs are not guaranteed to be globally unique across all time, but the probability of a collision within a single user's local database — which typically holds only a handful of packages — is negligible.

---

## TTL and garbage collection

TTL is specified at store time as a duration string (e.g. `7d`, `2h`) and converted immediately to an absolute `expires_at` UTC timestamp. It cannot be changed after the package is stored.

Supported TTL units:

| Unit | Meaning |
|------|---------|
| `h` | Hours |
| `d` | Days |

Both must be a positive integer (e.g. `14d`, not `0d` or `-1h`).

### Lazy garbage collection

`handoff` does not run a background daemon. Instead, it garbage-collects expired packages *lazily* — at the start of every database operation:

```sql
DELETE FROM packages WHERE expires_at < <current UTC time>
```

This means expired packages are silently removed whenever you run any `handoff` command. You will never see an expired package in `list` output, and `retrieve` will return `package not found` for a package whose TTL has elapsed.

### Explicit GC

Running `handoff gc` triggers the same deletion and reports the count:

```
Removed 3 expired package(s).
```

Use this if you want to confirm cleanup has happened or want to reclaim disk space immediately.

---

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `HANDOFF_DB` | `~/.handoff/handoff.db` | Path to the SQLite database file |

To use a custom location:

```bash
export HANDOFF_DB=/tmp/scratch.db
handoff store --name "temp" --ttl 2h
```

To use a project-scoped database that lives alongside the project:

```bash
export HANDOFF_DB="$(pwd)/.handoff.db"
```

There is no config file. This environment variable is the only configuration surface.

---

## Design principles

`handoff` was deliberately kept small. The design choices reflect these principles:

**Zero config.** Works out of the box with no setup. The first command you run creates the database automatically.

**Single binary.** `go build` produces one executable. No runtime dependencies, no external libraries to install.

**Cross-platform.** Pure Go, no CGO. Builds identically on macOS, Linux, and Windows for amd64 and arm64.

**Stdin/stdout.** Content flows through pipes. This is natural for agent terminal use and makes `handoff` composable with other tools.

**Ephemeral by default.** TTL ensures stored context eventually expires. Old knowledge packages do not accumulate forever.

**Simple over clever.** No plugin system, no network, no authentication, no daemon. A single SQLite file is the entire data layer.

---

## Uninstall

### Remove the binary (Homebrew)

```bash
brew uninstall handoff
```

### Remove the data directory

Homebrew only removes the binary. The data directory at `~/.handoff/` is left intact. To remove it:

```bash
rm -rf ~/.handoff
```

!!! warning "This deletes all stored packages"
    `~/.handoff/handoff.db` contains all your knowledge packages. Only delete it if you are sure you no longer need them. If you used a custom path via `HANDOFF_DB`, remove that file instead.

### If installed via `go install`

```bash
rm $(which handoff)
rm -rf ~/.handoff
```
