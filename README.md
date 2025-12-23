# wtree

A docker-ps style git worktree manager.

- Windows: Opens worktrees in Windows Terminal
- macOS/Linux: Opens worktrees in tmux

## Install

```bash
go install github.com/satoruhiga/wtree@latest
```

Or build from source:

```bash
go build -o wtree.exe .
```

## Usage

```bash
# Initialize (creates .wtree/config.toml)
wtree init

# Create a new worktree and open in Windows Terminal
wtree new
wtree new --pane    # Open in split pane

# List all worktrees
wtree ls

# Open existing worktree (partial ID match supported)
wtree open a3f8

# Remove a worktree
wtree rm a3f8
wtree rm a3f8 --force

# Merge and remove
wtree merge a3f8

# Clean up stale/merged worktrees
wtree prune
```

## Configuration

`.wtree/config.toml`:

```toml
[worktree]
worktree_base_dir = "../worktree"
branch_prefix = "wt/"
base_branch = "main"

[setup]
copy = [".env", ".claude/"]
commands = ["npm install"]

[terminal]
mode = "tab"  # "tab" | "pane" | "window"
```

## License

MIT
