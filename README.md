# VibeDiff

Git-native versioning layer for prompts and AI outputs.

[![Go Report Card](https://goreportcard.com/badge/github.com/vibediff/vibediff)](https://goreportcard.com/report/github.com/vibediff/vibediff)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/vibediff/vibediff?style=social)](https://github.com/vibediff/vibediff)

## Git vs VibeDiff

| Traditional Git | VibeDiff |
|----------------|----------|
| `git commit -m "update file"` | `vibediff commit "Add JWT authentication"` |
| Message: "what changed" | Message: why it changed |
| Lost AI context | Your intent preserved |

## Demo

```
$ vibediff commit "Add OAuth2 Google login with refresh tokens"
Staged 3 files
Committed with prompt attached

$ vibediff show
Commit: abc123f
Prompt: Add OAuth2 Google login with refresh tokens
Model: claude-opus-4
Files: auth.go, auth_test.go, middleware.go
```

## Quick Start

```bash
# Install
curl -sSL https://raw.githubusercontent.com/vibediff/vibediff/main/install.sh | sh

# Or with Go
go install github.com/vibediff/vibediff@latest

# Use with your AI workflow
vibediff commit "Add user authentication"
```

## Features

- Attach prompts to commits
- Review prompts alongside diffs
- Privacy-first: all data local
- Git-native: works with any repository
- Zero configuration

## Commands

```
vibediff commit "message"     Stage and commit with prompt
vibediff show                 Show prompt + diff
vibediff log                  List all vibe commits
vibediff export               Export to HTML
vibediff version              Version info
```

## Installation

### One-line install (Linux/macOS)

```bash
curl -sSL https://raw.githubusercontent.com/vibediff/vibediff/main/install.sh | sh
```

### Go install

```bash
go install github.com/vibediff/vibediff@latest
```

### Homebrew (macOS/Linux)

```bash
brew tap vibediff/vibediff
brew install vibediff
```

### Download binary

Download from [GitHub Releases](https://github.com/vibediff/vibediff/releases):

```bash
# Linux amd64
curl -sSL https://github.com/vibediff/vibediff/releases/latest/download/vibediff-linux-amd64 -o vibediff
chmod +x vibediff
sudo mv vibediff /usr/local/bin/

# macOS (Apple Silicon)
curl -sSL https://github.com/vibediff/vibediff/releases/latest/download/vibediff-darwin-arm64 -o vibediff
chmod +x vibediff
sudo mv vibediff /usr/local/bin/
```

### Build from source

```bash
git clone https://github.com/vibediff/vibediff.git
cd vibediff
make build
sudo make install
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/vibediff/vibediff/main/install.ps1 | iex
```

## Use Cases

- AI-assisted development: Track which prompts produced your best code
- Code reviews: See the intent behind changes, not just the diff
- Documentation: Maintain history of why code was written
- Auditing: Full traceability of AI-generated code changes

## How It Works

VibeDiff stores your prompts in `.vibediff/commits.json` within your git repository.

Files are protected with 0600 permissions - readable only by you.

## License

Apache 2.0 - see [LICENSE](LICENSE) file.

## Security

For security considerations, see [SECURITY.md](SECURITY.md).
