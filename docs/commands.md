# Commands

`handoff` has four subcommands. All of them share the same database, located at `~/.handoff/handoff.db` by default (see [Configuration](internals.md#configuration)).

---

## `handoff store`

Reads content from stdin and stores it as a named knowledge package. On success, prints the package ID — an 8-character lowercase hex string — to stdout.

```bash
echo "<content>" | handoff store --name <name> [options]
```

### Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--name` | **Yes** | — | Name for the package. Used to retrieve it by name later. |
| `--summary` | No | — | A short one-line description of the package contents. |
| `--ttl` | No | `7d` | How long the package lives before it expires. Format: a positive integer followed by `h` (hours) or `d` (days). |
| `--project` | No | — | A project scoping key. Use consistently across a project to group related packages. |
| `--tags` | No | — | Comma-separated tags, e.g. `auth,api,decisions`. Leading and trailing whitespace is trimmed from each tag. |

### Content

Content is always read from stdin. The command exits with an error if stdin is empty or contains only whitespace. There is no flag for passing content inline.

### Output

On success, the package ID is printed to stdout followed by a newline:

```
a3f9c12e
```

Report this ID to the user so they can retrieve the package by ID in the next session.

### Examples

```bash
# Minimal — just name and content
echo "context notes..." | handoff store --name "scratch"

# Full options
cat notes.md | handoff store \
  --name "auth-design" \
  --summary "JWT auth architecture decisions" \
  --ttl 14d \
  --project "myapp" \
  --tags "auth,api,decisions"

# Multi-line content using a heredoc
handoff store --name "sprint-state" --ttl 7d --project "myapp" << 'EOF'
## Current State
Feature X is complete. Feature Y is in progress.

## Next Steps
1. Finish Y
2. Write tests
EOF
```

### TTL reference

| Value | Duration |
|-------|----------|
| `2h` | 2 hours |
| `1d` | 1 day |
| `7d` | 7 days (default) |
| `14d` | 14 days |
| `30d` | 30 days |

TTL is specified at store time and cannot be changed after the fact. If a package expires before you retrieve it, it is deleted and no longer accessible.

---

## `handoff retrieve`

Retrieves a package and writes its content to stdout. Lookup is by ID or by name.

```bash
handoff retrieve <id>
handoff retrieve --name <name>
```

### Arguments

| Form | Description |
|------|-------------|
| `handoff retrieve <id>` | Retrieve by exact 8-character hex package ID. |
| `handoff retrieve --name <name>` | Retrieve by name. If multiple packages share the same name, the most recently stored one is returned. |

At least one of a positional ID argument or `--name` must be provided. If neither is given, the command exits with an error.

### Output

The package content is written to stdout exactly as it was stored. No trailing newline is added beyond what the original content contained. If no package matches, the command exits with a non-zero status and prints `package not found` to stderr.

### Examples

```bash
# Retrieve by ID
handoff retrieve a3f9c12e

# Retrieve by name
handoff retrieve --name "auth-design"

# Pipe content into a file for review
handoff retrieve --name "sprint-state" > context.md

# Pipe content to another tool
handoff retrieve a3f9c12e | pbcopy
```

---

## `handoff list`

Lists all non-expired packages in a tab-aligned table, ordered from most recently stored to oldest.

```bash
handoff list
handoff list --project <key>
```

### Flags

| Flag | Description |
|------|-------------|
| `--project` | Filter the results to packages matching the given project key. If omitted, all non-expired packages are shown. |

### Output

```
ID        NAME           PROJECT   TAGS          EXPIRES
a3f9c12e  auth-design    myapp     auth,api      2026-06-09 14:30
b1d2e3f4  db-schema      myapp     db,postgres   2026-06-02 09:15
```

Columns:

| Column | Description |
|--------|-------------|
| `ID` | 8-character package ID |
| `NAME` | Package name as given at store time |
| `PROJECT` | Project key, if set |
| `TAGS` | Comma-separated tags, if set |
| `EXPIRES` | Expiry timestamp in local time, formatted as `YYYY-MM-DD HH:MM` |

If no packages exist (or none match the project filter), the message `No packages found.` is printed to stderr and the command exits successfully.

### Examples

```bash
# List all packages
handoff list

# List packages for a specific project
handoff list --project myapp
```

---

## `handoff gc`

Manually removes all expired packages from the database and prints how many were deleted.

```bash
handoff gc
# Removed 3 expired package(s).
```

### When to use it

You generally do not need to run this command. Expired packages are automatically deleted on every database operation — whenever you run `store`, `retrieve`, `list`, or `gc`. Expired packages cannot be retrieved; they are simply taking up space on disk.

Run `handoff gc` explicitly if you want to reclaim disk space immediately, or to confirm how many stale packages have been cleaned up.
