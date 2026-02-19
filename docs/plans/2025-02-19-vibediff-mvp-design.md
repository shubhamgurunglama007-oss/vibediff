# VibeDiff MVP Design

**Date:** 2025-02-19
**Status:** Approved

## Vision

VibeDiff becomes the Git-native versioning layer for prompts and AI outputs.

Free, open-source core for developers. Paid enterprise SaaS for teams and organizations.

## MVP Scope

A Git CLI plugin that captures AI prompts alongside the code changes they generate.

## Architecture

```
vibediff/
├── cmd/
│   └── vibediff/          # Main CLI entrypoint
├── internal/
│   ├── git/               # Git operations wrapper
│   ├── store/             # JSON metadata storage
│   ├── diff/              # Diff capture & formatting
│   └── ui/                # Terminal + HTML output
├── pkg/
│   └── api/               # Shared types for future SaaS
└── .vibediff/             # Metadata (gitignored)
    └── commits.json       # Prompt → commit mappings
```

## Core Commands

| Command | Description |
|---------|-------------|
| `vibediff commit "prompt"` | Stage + commit with prompt attached |
| `vibediff show <ref>` | Display prompt + diff for a commit |
| `vibediff log` | Show all vibe-coded commits with prompts |
| `vibediff export --html` | Generate shareable HTML report |

## Data Model

```json
{
  "commit": "abc123...",
  "prompt": "Add user authentication",
  "timestamp": "2025-02-19T10:00:00Z",
  "diff": "diff --git a/auth.go b/auth.go\n+func Login()...",
  "files": ["auth.go", "middleware.go"],
  "model": "claude-sonnet-4"
}
```

## Tech Stack

- **Language:** Go (easy CLI, great cross-platform binaries)
- **License:** Apache 2.0 (permissive, enterprise-friendly)
- **Storage:** JSON in `.vibediff/` folder

## Key Design Decisions

| Decision | Rationale |
|----------|-----------|
| `vibediff commit` wrapper | Intentional, one-command workflow |
| Prompt + full diff | "This prompt produced exactly this change" — powerful visual |
| `.vibediff/` JSON storage | Simple, flexible, language-agnostic |
| CLI + HTML export | Fast terminal workflow + shareable output |
| Apache 2.0 license | Maximum adoption, enterprise-friendly |

## Success Metrics (Phase 1)

- 1,000+ GitHub stars
- 100+ weekly active developers
- 3 months to public launch

## Future Phases

- **Phase 2:** VS Code extension, GitHub Actions integration
- **Phase 3:** Enterprise SaaS layer (team dashboards, audit logs)
- **Phase 4:** Defensive strategy (trademark, moat building)
- **Phase 5:** Expansions (marketplace, benchmarking, partnerships)
