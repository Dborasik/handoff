# Commands

`handoff` has four subcommands. All of them operate on the same SQLite database, located at `~/.handoff/handoff.db` by default (configurable via [`HANDOFF_DB`](internals.md#configuration)).

| Command | Description |
|---------|-------------|
| [`store`](#handoff-store) | Read content from stdin and save it as a named knowledge package |
| [`retrieve`](#handoff-retrieve) | Fetch a package by ID or name and write its content to stdout |
| [`list`](#handoff-list) | List all non-expired packages in a table |
| [`gc`](#handoff-gc) | Manually delete all expired packages |

---

## handoff store

Reads content from stdin and saves it as a named knowledge package. On success, prints the package ID to stdout.

```bash
echo "<content>" | handoff store --name <name> [flags]
```

### Flags

`--name` *(required)*
:   A short, descriptive name for the package. Used to retrieve it by name in a future session. Names are **not unique** — multiple packages can share the same name. `retrieve --name` always returns the most recently stored match.

`--summary`
:   A single-line human-readable description of the package contents. Optional. Not used for lookup, only for display in `list` output.

`--ttl` *(default: `7d`)*
:   How long the package lives before it is automatically deleted. Format: a positive integer followed by `h` (hours) or `d` (days). Examples: `2h`, `1d`, `7d`, `14d`, `30d`. Cannot be changed after the package is stored.

`--project`
:   An optional project scoping key. Use the same key consistently across a project to group related packages together. Used as a filter in `handoff list --project`.

`--tags`
:   Optional comma-separated labels. Leading and trailing whitespace is trimmed from each tag. Example: `auth,api,decisions`.

### Content

Content is always read from stdin — there is no inline flag for it. The command exits with an error if stdin is empty or contains only whitespace.

### Output

On success, the package ID is printed to stdout:

```text
a3f9c12e
```

The ID is an 8-character lowercase hex string. Always pass this back to the user so they can retrieve the package by ID in the next session.

### Examples

```bash
# Minimal store
echo "notes about current state" | handoff store --name "scratch"

# With all options
cat notes.md | handoff store \
  --name "auth-design" \
  --summary "JWT auth architecture decisions" \
  --ttl 14d \
  --project "myapp" \
  --tags "auth,api,decisions"

# Using a heredoc for multi-line content
handoff store --name "sprint-state" --ttl 7d --project "myapp" << 'EOF'
## Current State
Feature X is complete. Feature Y is in progress.

## Next Steps
1. Finish Y
2. Write tests
EOF
```

### Errors

| Error | Cause |
|-------|-------|
| `Error: --name is required` | The `--name` flag was omitted |
| `Error: no content provided on stdin` | Stdin was empty or whitespace-only |
| `Error: invalid --ttl: unsupported unit 'X' (use h or d)` | TTL unit was not `h` or `d` |
| `Error: invalid --ttl: must be positive` | TTL value was `0` or negative |

---

## handoff retrieve

Retrieves a package and writes its content to stdout. Lookup is by package ID or by name.

```bash
handoff retrieve <id>
handoff retrieve --name <name>
```

### Arguments

`<id>` *(positional)*
:   The exact 8-character hex package ID returned by `store`. ID lookup is precise — there are no partial matches.

`--name`
:   Retrieve by package name. If multiple packages share the same name, the most recently stored one is returned. An expired package with that name is never returned.

!!! note "One or the other"
    You must provide either a positional ID argument or `--name`. Providing both is not an error — the positional ID takes priority.

### Output

The full package content is written to stdout exactly as it was stored. No metadata is added. Use shell redirection to save it to a file:

```bash
handoff retrieve --name "auth-design" > context.md
```

### Examples

```bash
# Retrieve by ID
handoff retrieve a3f9c12e

# Retrieve by name (returns most recent if name is reused)
handoff retrieve --name "auth-design"

# Save to file
handoff retrieve a3f9c12e > context.md

# Pipe to clipboard (macOS)
handoff retrieve --name "sprint-state" | pbcopy
```

### Errors

| Error | Cause |
|-------|-------|
| `Error: provide a package ID as argument or use --name` | Neither a positional ID nor `--name` was provided |
| `Error: package not found` | No non-expired package matched the given ID or name |

---

## handoff list

Lists all non-expired packages in a tab-aligned table, ordered from most recently stored to oldest.

```bash
handoff list
handoff list --project <key>
```

### Flags

`--project`
:   Filter the listing to packages that match the given project key. If omitted, all non-expired packages are shown.

### Output

```text
ID        NAME           PROJECT   TAGS          EXPIRES
a3f9c12e  auth-design    myapp     auth,api      2026-06-09 14:30
b1d2e3f4  db-schema      myapp     db,postgres   2026-06-02 09:15
```

| Column | Description |
|--------|-------------|
| `ID` | 8-character package ID |
| `NAME` | Package name as given at store time |
| `PROJECT` | Project key, if set |
| `TAGS` | Comma-separated tags, if set |
| `EXPIRES` | Expiry timestamp formatted as `YYYY-MM-DD HH:MM` |

!!! info "Empty results"
    If no packages exist (or none match the project filter), `No packages found.` is printed to stderr and the command exits with status `0`.

### Examples

```bash
# List all packages
handoff list

# List packages for a specific project
handoff list --project myapp
```

---

## handoff gc

Manually removes all expired packages from the database and prints the count deleted.

```bash
handoff gc
```

```text
Removed 3 expired package(s).
```

### When to run it

You generally do not need to run `handoff gc` explicitly. Expired packages are automatically deleted at the start of every database operation — whenever any `handoff` command runs. An expired package cannot be retrieved regardless of whether `gc` has been run.

Run `handoff gc` when you want to:

- Confirm how many stale packages have accumulated
- Reclaim disk space immediately rather than waiting for the next natural operation
