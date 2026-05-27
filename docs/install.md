# Install

`handoff` ships as a single static binary. There are three ways to get it.

---

## Homebrew (recommended)

The fastest path on macOS or Linux. No Go toolchain required.

```bash
brew tap Dborasik/tap
brew install handoff
```

To upgrade later:

```bash
brew upgrade handoff
```

---

## Go install

Requires Go 1.26+. No CGO, no system dependencies — the binary is fully self-contained.

```bash
go install github.com/Dborasik/handoff@latest
```

The binary is placed in your `$GOPATH/bin` (typically `~/go/bin`). Make sure that directory is on your `$PATH`.

---

## Build from source

```bash
git clone https://github.com/Dborasik/handoff.git
cd handoff
go build -o handoff .
```

Move the resulting binary somewhere on your `$PATH`:

```bash
mv handoff /usr/local/bin/
```

---

## Verify the installation

```bash
handoff --help
```

You should see the top-level help output listing the four available subcommands.

---

## Pre-built binaries

Every [GitHub release](https://github.com/Dborasik/handoff/releases) includes pre-built archives for all supported platforms:

| Platform | Archive |
|----------|---------|
| macOS (Apple Silicon) | `handoff_darwin_arm64.tar.gz` |
| macOS (Intel) | `handoff_darwin_amd64.tar.gz` |
| Linux (x86-64) | `handoff_linux_amd64.tar.gz` |
| Linux (ARM64) | `handoff_linux_arm64.tar.gz` |
| Windows (x86-64) | `handoff_windows_amd64.zip` |
| Windows (ARM64) | `handoff_windows_arm64.zip` |

A `checksums.txt` file is included with each release for verification.

---

## Uninstall

See [Uninstall](internals.md#uninstall) for instructions, including how to remove the data directory.
