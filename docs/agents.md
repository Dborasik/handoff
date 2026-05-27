# Agent Setup

To get your agent to use `handoff` automatically, you need to add instructions to your project. There are two approaches: always-on instruction files and on-demand skill files. Both approaches are supported across all major AI coding agents.

---

## Option A — Always-on instruction files

The agent reads these files on every session start, so it will always be aware of `handoff` and check for existing packages automatically.

| Agent | File to add to your project | Source in this repo |
|-------|------------------------------|---------------------|
| Claude Code | `CLAUDE.md` | [`instructions/CLAUDE.md`](https://github.com/Dborasik/handoff/blob/main/instructions/CLAUDE.md) |
| GitHub Copilot | `.github/copilot-instructions.md` | [`instructions/copilot-instructions.md`](https://github.com/Dborasik/handoff/blob/main/instructions/copilot-instructions.md) |
| OpenAI Codex | `AGENTS.md` | [`instructions/AGENTS.md`](https://github.com/Dborasik/handoff/blob/main/instructions/AGENTS.md) |
| Cursor | `.cursor/rules/handoff.mdc` | [`instructions/cursor.mdc`](https://github.com/Dborasik/handoff/blob/main/instructions/cursor.mdc) |

Use `curl` to pull the file directly into your project:

=== "Claude Code"

    ```bash
    curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/instructions/CLAUDE.md \
      -o CLAUDE.md
    ```

=== "GitHub Copilot"

    ```bash
    mkdir -p .github
    curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/instructions/copilot-instructions.md \
      -o .github/copilot-instructions.md
    ```

=== "OpenAI Codex"

    ```bash
    curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/instructions/AGENTS.md \
      -o AGENTS.md
    ```

=== "Cursor"

    ```bash
    mkdir -p .cursor/rules
    curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/instructions/cursor.mdc \
      -o .cursor/rules/handoff.mdc
    ```

!!! note "Existing instruction files"
    If your project already has a `CLAUDE.md` or `AGENTS.md`, append the `handoff` instructions to your existing file rather than replacing it. The instructions are plain markdown — they compose cleanly.

---

## Option B — On-demand skill files

A skill file is loaded by the agent only when it is relevant to the current task — it does not appear in every session's context. This is a lighter-weight option that works well if you already have large instruction files.

| Location in your project | Supported by |
|--------------------------|--------------|
| `.github/skills/handoff/SKILL.md` | GitHub Copilot |
| `.agents/skills/handoff/SKILL.md` | OpenAI Codex and other agents |
| `.claude/skills/handoff/SKILL.md` | Claude Code |

=== "Claude Code"

    ```bash
    mkdir -p .claude/skills/handoff
    curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/.claude/skills/handoff/SKILL.md \
      -o .claude/skills/handoff/SKILL.md
    ```

=== "GitHub Copilot"

    ```bash
    mkdir -p .github/skills/handoff
    curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/.github/skills/handoff/SKILL.md \
      -o .github/skills/handoff/SKILL.md
    ```

=== "OpenAI Codex / other agents"

    ```bash
    mkdir -p .agents/skills/handoff
    curl -fsSL https://raw.githubusercontent.com/Dborasik/handoff/main/.agents/skills/handoff/SKILL.md \
      -o .agents/skills/handoff/SKILL.md
    ```

---

## What the instructions tell the agent

Regardless of which option you choose, the instruction content teaches the agent to:

1. **Check for existing packages at session start** — run `handoff list` (or `handoff list --project <name>`) and retrieve any relevant packages before doing anything else.
2. **Offer a knowledge transfer proactively** — when the context is getting large, suggest storing a package before the window is exhausted.
3. **Respond to explicit requests** — phrases like "do a handoff", "save context", "knowledge transfer", or "store session state" trigger an immediate store.
4. **Verify installation** — check that `handoff` is on the PATH before trying to use it, and guide the user to install it if not.

---

## Choosing between Option A and Option B

| | Option A (always-on) | Option B (skill file) |
|--|---|---|
| Agent reads instructions | Every session | Only when relevant |
| Context cost | Small constant overhead | Near zero when not in use |
| Best for | Projects where you regularly reach context limits | Projects where handoff is only occasionally needed |

Both options result in the same agent behaviour once triggered. If you are unsure, start with Option B (skill file) and switch to Option A if you find yourself manually prompting the agent to use `handoff` frequently.
