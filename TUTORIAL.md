# VibeDiff Tutorial

A complete guide to using VibeDiff for AI-assisted development.

## What is VibeDiff?

VibeDiff is a Git-native tool that attaches your AI prompts to your commits. It captures the **intent** behind your code changes - perfect for AI-assisted development.

**Key idea:** When AI helps you write code, the commit message should reflect your intent, not just "update file". VibeDiff makes this easy.

---

## Installation

```bash
# Clone and build
git clone https://github.com/vibediff/vibediff.git
cd vibediff
make build

# Or install globally
sudo make install
```

---

## Basic Usage

### Your First Vibe Commit

```bash
# Make some code changes (with AI help or not)
vim auth.go

# Stage and commit with your AI prompt
vibediff commit "Add JWT authentication middleware"
```

What happens:
1. Stages all changes (`git add .`)
2. Creates a commit with message `vibe: Add JWT authentication middleware`
3. Saves your prompt + diff in `.vibediff/commits.json`

### See What You Saved

```bash
# Show the most recent vibe commit
vibediff show
```

Output:
```
Commit: abc1234
When:   2026-02-20T14:30:00+05:45

Prompt:
  Add JWT authentication middleware

Files: auth.go, middleware.go

Diff:
diff --git a/auth.go b/auth.go
@@ -1,1 +1,5 @@
+func validateJWT(token string) bool {
+    // ...
+}
```

---

## Workflow Examples

### 1. Pure AI-Assisted Development

```bash
# 1. Ask Claude/Cursor to add authentication
# AI generates code, you apply it

# 2. Commit with the intent
vibediff commit "Add JWT authentication with refresh token support"

# 3. Later, remember WHY you did this
vibediff show abc123
```

### 2. Hybrid Workflow (vibediff + git commit)

```bash
# AI-assisted feature - use vibediff
vibediff commit "Add user registration with email verification"

# Quick typo fix - use regular git
git commit -am "fix typo in main.go"

# Another AI change - use vibediff
vibediff commit "Refactor database layer for better query performance"

# Regular merge/rebase commits - use git
git commit
```

### 3. See All Your AI-Assisted Commits

```bash
# List all vibe commits with prompts
vibediff log
```

Output:
```
Found 3 vibe commit(s):

1. abc1234
  2026-02-20T14:30:00+05:45
  Add user registration with email verification

2. def5678
  2026-02-20T15:45:00+05:45
  Refactor database layer for better query performance

3. ghi9012
  2026-02-20T16:20:00+05:45
  Add JWT authentication with refresh token support
```

---

## Commands Reference

### `vibediff commit <prompt>`

Stage all changes and create a commit with your prompt.

```bash
vibediff commit "Add OAuth2 Google login"

# Prompt is trimmed automatically
vibediff commit "   Add OAuth2 Google login   "
# Same as: "Add OAuth2 Google login"

# Empty prompts are rejected
vibediff commit "   "
# Error: prompt cannot be empty

# Prompts are limited to 100KB
```

### `vibediff show [commit-ref]`

Display the prompt and diff for a commit.

```bash
# Show most recent vibe commit
vibediff show

# Show specific commit
vibediff show abc1234

# Show by tag
vibediff show v1.0.0
```

### `vibediff log`

Show all vibe commits with their prompts.

```bash
vibediff log
```

### `vibediff export [commit-ref]`

Export a vibe commit to HTML for sharing.

```bash
# Export most recent
vibediff export

# Creates: vibe-abc1234.html
```

### `vibediff version`

Show version information.

```bash
vibediff version
# VibeDiff 1.0.0
# Commit: abc123
# Built: 2026-02-20T00:00:00Z
```

---

## Understanding the Data

### Where Data is Stored

All data lives in your repository:

```
my-project/
├── .vibediff/
│   └── commits.json    # Your prompts + diffs
├── src/
│   └── main.go
└── README.md
```

### What Gets Stored

For each vibe commit, VibeDiff stores:

```json
{
  "commits": {
    "abc1234": {
      "commit": "abc1234def5678",
      "prompt": "Add JWT authentication",
      "timestamp": "2026-02-20T14:30:00Z",
      "diff": "--- a/auth.go\n+++ b/auth.go\n...",
      "files": ["auth.go", "middleware.go"]
    }
  }
}
```

**Privacy:** All data is:
- Stored locally (no cloud)
- Private (0600 permissions)
- Git-native (committed to your repo)

---

## Common Workflows

### Workflow 1: Review AI Decisions Later

```bash
# Weeks ago, you asked AI to build auth
vibediff commit "Add authentication with OAuth2"

# Now you need to remember WHY you chose OAuth2
vibediff show <commit>

# Ah, there's your prompt + the full diff
```

### Workflow 2: Share AI Context with Team

```bash
# Export to HTML and send to teammate
vibediff export abc123

# They can see your intent + actual changes
# Without needing access to your AI conversation
```

### Workflow 3: Document AI-Generated Features

```bash
# See all AI-assisted features
vibediff log

# Export documentation for each
vibediff export v1.0.0 > docs/auth-v1.md
```

### Workflow 4: Understand Code Changes

```bash
# Regular git log shows WHAT changed
git log --oneline

# VibeDiff shows WHY it changed
vibediff log

# Combine both for full picture
git log --oneline | grep "vibe:"
```

---

## Tips & Best Practices

### Writing Good Prompts

```bash
# ✅ Good - Describes intent
vibediff commit "Add JWT authentication with refresh token support"

# ✅ Good - Describes architectural decision
vibediff commit "Switch from REST to GraphQL for user API"

# ❌ Too vague
vibediff commit "update"

# ❌ Just describes files
vibediff commit "changes to auth.go"
```

### When to Use vs Regular Git Commit

| Use `vibediff commit` for... | Use `git commit` for... |
|------------------------------|----------------------|
| AI-assisted features | Typos |
| Complex refactors | Format changes |
| Architectural decisions | Merge commits |
| Bug fixes with reasoning | Changelog updates |
| Experiments | ".gitignore" updates |

### Hybrid Workflow Example

```bash
# Feature development
vibediff commit "Add user profile API endpoint"
vibediff commit "Implement password reset flow"
git commit -am "fix typo in endpoint url"

# Documentation
vibediff commit "Add API documentation for profile endpoints"

# Cleanup
git commit -am "remove debug logging"
```

---

## Privacy & Security

- **Local only:** No data leaves your machine
- **Private files:** `commits.json` uses 0600 permissions
- **Input validation:** All refs and prompts are validated
- **Git-native:** `.vibediff` can be added to `.gitignore` if desired

---

## Troubleshooting

### "no vibe metadata found"

The commit wasn't made with `vibediff commit`. Use `git log` instead.

### ".vibediff in wrong location"

Upgrade to latest version - fixed to always use git root.

### "invalid commit ref"

Commit refs are validated for security. Use valid refs like `HEAD`, `v1.0.0`, or commit hashes.

---

## Examples Repository

See `examples/basic/` for a sample project using VibeDiff.
