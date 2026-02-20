# Contributing to VibeDiff

Thank you for your interest in contributing!

## Development Setup

```bash
# Clone the repository
git clone https://github.com/vibediff/vibediff.git
cd vibediff

# Install dependencies
go mod download

# Run tests
go test ./...

# Build locally
go build -o bin/vibediff ./cmd/vibediff
# or use make
make build
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Write tests for new features (TDD preferred)
- Keep functions focused and small
- Add godoc comments for exported types

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./internal/git -v -run TestValidateCommitRef
```

We aim for 80%+ test coverage.

## Security

VibeDiff handles user prompts and diffs as sensitive data. When contributing:

- All data files must use 0600 permissions (owner read/write only)
- Directories must use 0700 permissions
- Validate all user inputs
- Use parameterized queries/commands to prevent injection
- Follow the security checklist in the implementation plan

## Submitting Changes

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Ensure all tests pass
5. Submit a pull request

## License

By contributing, you agree that your contributions will be licensed under the Apache 2.0 License.
