# Troubleshooting

## "not a git repository" Error

**Problem:** `vibediff commit` fails with "not in a git repository"

**Solution:** You must be inside a git repository. Run `git status` to verify.

## "no vibe metadata found" Error

**Problem:** `vibediff show` can't find metadata for a commit

**Solutions:**
1. The commit wasn't made with `vibediff commit`
2. You're in a different repository than where the vibe commit was made
3. The `.vibediff/commits.json` file was deleted

## .vibediff Directory in Wrong Location

**Problem:** `.vibediff` appears in subdirectories instead of repo root

**Solution:** This was a bug in versions < 0.1.0. Upgrade to latest version and move your `.vibediff` directory to the repository root.

## Corrupt commits.json

**Problem:** `vibediff log` fails with JSON error

**Solution:**
```bash
# The file may be corrupted. Check it:
cat .vibediff/commits.json | jq .

# If corrupted, you may need to restore from backup or manually fix the JSON
```

## Permissions Error on commits.json

**Problem:** Can't read/write `.vibediff/commits.json`

**Solution:**
```bash
# Fix permissions
chmod 600 .vibediff/commits.json
chmod 700 .vibediff
```

## "invalid commit ref" Error

**Problem:** Getting validation errors when using `show` or `export` with a commit ref

**Solution:** Commit refs are validated for security. Only alphanumeric characters, slashes, dots, dashes, and colons are allowed. Parent directory traversal (`../`) is blocked.

## Empty Prompt Error

**Problem:** `vibediff commit` rejects empty or whitespace-only prompts

**Solution:** This is intentional. Provide a meaningful prompt describing your changes.

## Prompt Too Large Error

**Problem:** Getting "prompt too large" error

**Solution:** Prompts are limited to 100KB to prevent denial-of-service. If you have a legitimate use case for larger prompts, please file an issue.
