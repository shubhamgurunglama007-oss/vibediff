# VibeDiff

> Git-native versioning layer for prompts and AI outputs.

Attach your AI prompts to git commits. Perfect for AI-assisted development.

## Quick Start

```bash
# Install
go install github.com/vibediff/vibediff@latest

# Make changes (with AI help or not)
vim auth.go

# Commit with your AI prompt
vibediff commit "Add JWT authentication middleware"

# Later, see what you intended
vibediff show
```

## What It Does

| Traditional Git | VibeDiff |
|----------------|----------|
| `git commit -m "update file"` | `vibediff commit "Add user authentication"` |
| Message: "what changed" | Message: "why it changed" |
| Lost context | Your intent preserved |

## Features

- ✅ **Attach prompts to commits** - Capture your AI intent
- ✅ **Review later** - See prompts + diffs with `vibediff show`
- ✅ **Privacy-first** - All data local, 0600 file permissions
- ✅ **Git-native** - Works with any git repository
- ✅ **Zero config** - Just run `vibediff commit`
- ✅ **Hybrid workflow** - Mix `vibediff` and `git commit` freely

## Commands

```bash
vibediff commit "Add user authentication"   # Stage + commit with prompt
vibediff show                               # Show prompt + diff
vibediff log                                 # List all vibe commits
vibediff export                              # Export to HTML
vibediff version                             # Version info
```

## Example

```bash
# AI helps you add authentication
vibediff commit "Add OAuth2 Google login with refresh tokens"

# Later, you can see:
# - Your prompt: "Add OAuth2 Google login with refresh tokens"
# - The diff: Full code changes
# - Context: Files changed, timestamp

vibediff show
```

## Learn More

- **[Tutorial](TUTORIAL.md)** - Complete guide with examples
- **[Troubleshooting](docs/troubleshooting.md)** - Common issues and fixes
- **[Contributing](CONTRIBUTING.md)** - Development guide

## Installation

```bash
go install github.com/vibediff/vibediff@latest

# Or build from source
git clone https://github.com/vibediff/vibediff.git
cd vibediff
make build
sudo make install
```

## License

Apache 2.0
