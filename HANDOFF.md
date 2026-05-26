# Knowledge Transfer: Publishing `handoff` CLI via Homebrew

## What This Is

You are picking up a Go CLI project called `handoff`. It stores and retrieves "knowledge packages" — markdown blobs that AI agents pass between sessions via a local SQLite database. The code is complete and compiles. Your job is to publish it so users can install with `brew install`.

---

## Current State

- **Code**: Complete, compiles, tested manually
- **Module path**: `github.com/agent-handoff/agent-handoff`
- **Binary name**: `handoff`
- **Location**: This repo's root directory
- **Dependencies**: `modernc.org/sqlite` (pure Go, no CGO), `github.com/spf13/cobra`
- **GoReleaser config**: `.goreleaser.yml` is ready
- **License**: MIT

---

## If You Want to Rename from `agent-handoff`

The current GitHub org + repo name is `agent-handoff/agent-handoff`. To change it:

1. **Pick your new `<org>/<repo>`** (e.g., `myorg/handoff`)
2. **go.mod**: Change `module github.com/agent-handoff/agent-handoff` → `module github.com/<org>/<repo>`
3. **All imports**: Find/replace `github.com/agent-handoff/agent-handoff/` → `github.com/<org>/<repo>/` in:
   - `main.go`
   - `cmd/store.go`, `cmd/retrieve.go`, `cmd/list.go`, `cmd/gc.go`
   - `internal/db/db.go`
4. **.goreleaser.yml**: Update `owner` and `name` fields under `brews[0].repository`, and under `release.github`
5. **Homebrew tap repo**: Name it `homebrew-tap` under your new org (or `homebrew-<repo>`)
6. Run `go build .` to verify it still compiles

---

## Steps to Publish

### 1. Create GitHub Repos

You need TWO repos:

| Repo | Purpose |
|------|---------|
| `github.com/agent-handoff/agent-handoff` | The source code (this project) |
| `github.com/agent-handoff/homebrew-tap` | Homebrew formula (GoReleaser auto-updates this) |

Create both on GitHub. The `homebrew-tap` repo can be empty (just needs to exist with a README).

### 2. Push the Source Code

```bash
cd <this-project-root>
git init
git add .
git commit -m "initial commit"
git remote add origin git@github.com:agent-handoff/agent-handoff.git
git branch -M main
git push -u origin main
```

### 3. Create a GitHub Personal Access Token

GoReleaser needs a token to:
- Create releases on the source repo
- Push the formula to the homebrew-tap repo

Create a **fine-grained personal access token** at https://github.com/settings/tokens with:
- **Repository access**: Both repos above
- **Permissions**: Contents (read/write), Metadata (read)

Export it:
```bash
export GITHUB_TOKEN="ghp_xxxxxxxxxxxxx"
```

### 4. Install GoReleaser

```bash
brew install goreleaser
```

### 5. Tag and Release

```bash
git tag v0.1.0
git push origin v0.1.0
goreleaser release --clean
```

This will:
- Build binaries for macOS (amd64 + arm64), Linux (amd64 + arm64), Windows
- Create a GitHub Release with the binaries attached
- Auto-generate and push a Homebrew formula to `agent-handoff/homebrew-tap`

### 6. Verify Installation

```bash
brew tap agent-handoff/tap
brew install handoff
handoff --help
```

---

## Ongoing Releases

For future versions:
```bash
git tag v0.2.0
git push origin v0.2.0
goreleaser release --clean
```

Users update with:
```bash
brew upgrade handoff
```

---

## How Users Install (Final UX)

```bash
brew tap agent-handoff/tap
brew install handoff
```

Or as a one-liner:
```bash
brew install agent-handoff/tap/handoff
```

---

## Project Structure Reference

```
├── .goreleaser.yml          # Build + release config
├── LICENSE                  # MIT
├── ARCHITECTURE.md          # Design doc
├── main.go                  # Entry point
├── cmd/
│   ├── root.go              # CLI root (cobra)
│   ├── store.go             # `handoff store` — stdin → DB
│   ├── retrieve.go          # `handoff retrieve` — DB → stdout
│   ├── list.go              # `handoff list`
│   └── gc.go                # `handoff gc`
├── internal/
│   ├── config/config.go     # DB path (~/.handoff/handoff.db)
│   ├── db/db.go             # All SQLite operations
│   ├── db/schema.go         # Table DDL
│   └── model/package.go     # Package struct
├── go.mod
└── go.sum
```

---

## Notes

- No CGO. Pure Go. Cross-compiles cleanly.
- GoReleaser handles all multi-platform builds. Don't build manually for releases.
- The homebrew formula is auto-generated. Don't edit it by hand.
- If you change the binary name from `handoff`, update `.goreleaser.yml` → `builds[0].binary` and `brews[0].install`.
