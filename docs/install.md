# Install

`handoff` ships as a single static binary with no runtime dependencies. Choose the method that suits your setup.

---

## Methods

=== "Homebrew"

    The recommended method on macOS and Linux. No Go toolchain required.

    ```bash
    brew tap Dborasik/tap
    brew install handoff
    ```

    To upgrade to the latest release:

    ```bash
    brew upgrade handoff
    ```

=== "Go install"

    !!! info "Requirement"
        Requires Go 1.26 or later. Run `go version` to check. The binary is fully self-contained — no CGO, no system libraries.

    ```bash
    go install github.com/Dborasik/handoff@latest
    ```

    The binary is placed in `$GOPATH/bin` (typically `~/go/bin`). Make sure that directory is on your `$PATH`:

    ```bash
    export PATH="$PATH:$(go env GOPATH)/bin"
    ```

=== "Build from source"

    ```bash
    git clone https://github.com/Dborasik/handoff.git
    cd handoff
    go build -o handoff .
    ```

    Then move the binary somewhere on your `$PATH`:

    ```bash
    mv handoff /usr/local/bin/
    ```

---

## Verify

After installing, confirm the binary is on your PATH:

```bash
handoff --help
```

You should see:

```text
A CLI tool for storing and retrieving knowledge packages across AI agent context windows.

Usage:
  handoff [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  gc          Remove all expired packages
  help        Help about any command
  list        List available knowledge packages
  retrieve    Retrieve a knowledge package by ID or name
  store       Store a knowledge package (reads content from stdin)

Flags:
  -h, --help   help for handoff

Use "handoff [command] --help" for more information about a command.
```

---

## Pre-built binaries

Every [GitHub release](https://github.com/Dborasik/handoff/releases) includes pre-built archives. A `checksums.txt` file is included for verification.

| Platform | Architecture | Archive |
|----------|-------------|---------|
| macOS | Apple Silicon (arm64) | `handoff_darwin_arm64.tar.gz` |
| macOS | Intel (amd64) | `handoff_darwin_amd64.tar.gz` |
| Linux | x86-64 (amd64) | `handoff_linux_amd64.tar.gz` |
| Linux | ARM64 | `handoff_linux_arm64.tar.gz` |
| Windows | x86-64 (amd64) | `handoff_windows_amd64.zip` |
| Windows | ARM64 | `handoff_windows_arm64.zip` |

**To use a pre-built binary:**

=== "macOS / Linux"

    ```bash
    # Download and extract (example: macOS Apple Silicon)
    curl -L https://github.com/Dborasik/handoff/releases/latest/download/handoff_darwin_arm64.tar.gz \
      | tar -xz

    # Move to a directory on your PATH
    mv handoff /usr/local/bin/
    ```

=== "Windows"

    Download the `.zip` from the [releases page](https://github.com/Dborasik/handoff/releases), extract it, and place `handoff.exe` in a directory that is on your `%PATH%`.

---

## Uninstall

=== "Homebrew"

    ```bash
    brew uninstall handoff
    ```

=== "Go install / manual"

    ```bash
    rm $(which handoff)
    ```

After removing the binary, the data directory at `~/.handoff/` is left intact.

!!! warning "Data directory"
    `~/.handoff/handoff.db` contains all your stored knowledge packages. Delete it only if you are sure you no longer need them.

    ```bash
    rm -rf ~/.handoff
    ```

    If you used a custom database path via `HANDOFF_DB`, remove that file instead.
