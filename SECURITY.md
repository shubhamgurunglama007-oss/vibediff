# Security Policy

## Supported Versions

Current version: 1.0.0

Security updates are provided for the latest version.

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly.

**Email:** security@vibediff.com

Please include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if known)

We will acknowledge receipt within 48 hours and provide regular updates on our progress.

## Security Design

VibeDiff is designed with security in mind:

- **Local-only storage:** All data is stored locally in `.vibediff/`
- **File permissions:** Data files use 0600 permissions (owner read/write only)
- **Directory permissions:** Directories use 0700 permissions (owner access only)
- **No network calls:** VibeDiff does not make any network requests
- **No telemetry:** No data is sent to external servers

## Data Storage

VibeDiff stores prompts in `.vibediff/commits.json` within your git repository.

- File permissions: 0600
- Contains: prompts, commit hashes, timestamps, model names, tags
- Format: JSON

This file should NOT be committed to git if it contains sensitive prompts.

## Recommendations

- Add `.vibediff/` to your `.gitignore` to avoid committing sensitive prompts
- Review `.vibediff/commits.json` before sharing repositories
- Use environment-specific prompts for sensitive operations
