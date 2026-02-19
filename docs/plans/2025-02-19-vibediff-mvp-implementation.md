# VibeDiff MVP Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a Git CLI plugin that captures AI prompts alongside the code changes they generate.

**Architecture:** Go-based CLI tool with JSON metadata storage. Wrapper around Git commands that captures prompts and diffs together.

**Tech Stack:** Go 1.23+, go-git library, Cobra CLI framework, JSON storage

---

## Task 1: Project Initialization

**Files:**
- Create: `go.mod`
- Create: `.gitignore`
- Create: `LICENSE`

**Step 1: Initialize Go module**

Run: `go mod init github.com/vibediff/vibediff`

**Step 2: Create .gitignore**

```gitignore
# Binaries
vibediff
*.exe
*.dll
*.so
*.dylib

# Metadata (user can choose to commit)
.vibediff/

# Test coverage
*.out

# IDE
.idea/
.vscode/
*.swp
```

**Step 3: Create Apache 2.0 LICENSE**

Run: `curl -o LICENSE https://www.apache.org/licenses/LICENSE-2.0.txt`

**Step 4: Commit**

```bash
git add go.mod .gitignore LICENSE
git commit -m "chore: initialize Go module and project files"
```

---

## Task 2: Core Data Types

**Files:**
- Create: `internal/store/types.go`
- Create: `internal/store/types_test.go`

**Step 1: Write failing test**

Create `internal/store/types_test.go`:

```go
package store

import (
    "testing"
    "time"
)

func TestCommitMetadata_Valid(t *testing.T) {
    meta := CommitMetadata{
        Commit:    "abc123",
        Prompt:    "Add user authentication",
        Timestamp: time.Now().Format(time.RFC3339),
        Files:     []string{"auth.go"},
    }

    if meta.Commit != "abc123" {
        t.Errorf("expected commit abc123, got %s", meta.Commit)
    }
    if meta.Prompt != "Add user authentication" {
        t.Errorf("unexpected prompt: %s", meta.Prompt)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/store/...`
Expected: FAIL with "undefined: CommitMetadata"

**Step 3: Write minimal implementation**

Create `internal/store/types.go`:

```go
package store

import "time"

// CommitMetadata stores the prompt and resulting diff for a commit
type CommitMetadata struct {
    Commit    string   `json:"commit"`
    Prompt    string   `json:"prompt"`
    Timestamp string   `json:"timestamp"`
    Diff      string   `json:"diff,omitempty"`
    Files     []string `json:"files,omitempty"`
    Model     string   `json:"model,omitempty"`
}

// Store manages vibe diff metadata
type Store struct {
    dataPath string
}

// NewStore creates a new store with the given data path
func NewStore(dataPath string) *Store {
    return &Store{dataPath: dataPath}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/store/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/store/types.go internal/store/types_test.go
git commit -m "feat: add core data types for commit metadata"
```

---

## Task 3: JSON Storage Backend

**Files:**
- Modify: `internal/store/types.go`
- Create: `internal/store/json.go`
- Create: `internal/store/json_test.go`

**Step 1: Write failing test**

Create `internal/store/json_test.go`:

```go
package store

import (
    "os"
    "path/filepath"
    "testing"
    "time"
)

func TestJSONStore_SaveAndLoad(t *testing.T) {
    tmpDir := t.TempDir()
    dataFile := filepath.Join(tmpDir, "commits.json")
    js := NewJSONStore(dataFile)

    meta := CommitMetadata{
        Commit:    "abc123",
        Prompt:    "Test prompt",
        Timestamp: time.Now().Format(time.RFC3339),
        Files:     []string{"test.go"},
    }

    // Save
    err := js.Save(meta)
    if err != nil {
        t.Fatalf("failed to save: %v", err)
    }

    // Load
    loaded, err := js.Load("abc123")
    if err != nil {
        t.Fatalf("failed to load: %v", err)
    }

    if loaded.Prompt != "Test prompt" {
        t.Errorf("expected prompt 'Test prompt', got %s", loaded.Prompt)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/store/... -v`
Expected: FAIL with "undefined: NewJSONStore"

**Step 3: Write implementation**

Create `internal/store/json.go`:

```go
package store

import (
    "encoding/json"
    "os"
    "sync"
)

// JSONStore persists metadata to a JSON file
type JSONStore struct {
    mu       sync.RWMutex
    dataPath string
}

// NewJSONStore creates a new JSON-based store
func NewJSONStore(dataPath string) *JSONStore {
    return &JSONStore{dataPath: dataPath}
}

// data represents the JSON file structure
type data struct {
    Commits map[string]CommitMetadata `json:"commits"`
}

// Save stores a commit metadata entry
func (js *JSONStore) Save(meta CommitMetadata) error {
    js.mu.Lock()
    defer js.mu.Unlock()

    // Read existing data
    d := data{Commits: make(map[string]CommitMetadata)}
    if content, err := os.ReadFile(js.dataPath); err == nil {
        json.Unmarshal(content, &d)
    }

    // Add new entry
    d.Commits[meta.Commit] = meta

    // Write back
    content, err := json.MarshalIndent(d, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(js.dataPath, content, 0644)
}

// Load retrieves a commit metadata entry
func (js *JSONStore) Load(commit string) (CommitMetadata, error) {
    js.mu.RLock()
    defer js.mu.RUnlock()

    content, err := os.ReadFile(js.dataPath)
    if err != nil {
        return CommitMetadata{}, err
    }

    var d data
    if err := json.Unmarshal(content, &d); err != nil {
        return CommitMetadata{}, err
    }

    meta, ok := d.Commits[commit]
    if !ok {
        return CommitMetadata{}, os.ErrNotExist
    }

    return meta, nil
}

// List returns all commit metadata
func (js *JSONStore) List() ([]CommitMetadata, error) {
    js.mu.RLock()
    defer js.mu.RUnlock()

    content, err := os.ReadFile(js.dataPath)
    if err != nil {
        if os.IsNotExist(err) {
            return []CommitMetadata{}, nil
        }
        return nil, err
    }

    var d data
    if err := json.Unmarshal(content, &d); err != nil {
        return nil, err
    }

    result := make([]CommitMetadata, 0, len(d.Commits))
    for _, meta := range d.Commits {
        result = append(result, meta)
    }
    return result, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/store/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/store/json.go internal/store/json_test.go
git commit -m "feat: add JSON storage backend"
```

---

## Task 4: Git Operations Wrapper

**Files:**
- Create: `internal/git/git.go`
- Create: `internal/git/git_test.go`

**Step 1: Write failing test**

Create `internal/git/git_test.go`:

```go
package git

import (
    "os"
    "path/filepath"
    "testing"
)

func TestGetStatus_GetDiff(t *testing.T) {
    // Create a test repo
    tmpDir := t.TempDir()
    repo, err := Init(tmpDir)
    if err != nil {
        t.Fatalf("failed to init repo: %v", err)
    }

    // Create a test file
    testFile := filepath.Join(tmpDir, "test.txt")
    if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
        t.Fatal(err)
    }

    // Stage the file
    if err := repo.Add("test.txt"); err != nil {
        t.Fatalf("failed to add: %v", err)
    }

    // Get diff
    diff, err := repo.Diff("staged")
    if err != nil {
        t.Fatalf("failed to get diff: %v", err)
    }

    if diff == "" {
        t.Error("expected non-empty diff")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/git/... -v`
Expected: FAIL with "package not found" or "undefined"

**Step 3: Write implementation**

Create `internal/git/git.go`:

```go
package git

import (
    "bytes"
    "fmt"
    "os/exec"
    "strings"
)

// Repo represents a Git repository
type Repo struct {
    path string
}

// Init initializes a new Repo for the given path
func Init(path string) (*Repo, error) {
    // Verify it's a git repo
    _, err := runGit(path, "rev-parse", "--git-dir")
    if err != nil {
        return nil, fmt.Errorf("not a git repository: %w", err)
    }
    return &Repo{path: path}, nil
}

// Add stages files for commit
func (r *Repo) Add(files ...string) error {
    _, err := runGit(r.path, append([]string{"add"}, files...)...)
    return err
}

// Commit creates a new commit
func (r *Repo) Commit(message string) error {
    _, err := runGit(r.path, "commit", "-m", message)
    return err
}

// Diff returns the diff for staged or unstaged changes
func (r *Repo) Diff(which string) (string, error) {
    var args []string
    if which == "staged" {
        args = []string{"diff", "--staged"}
    } else {
        args = []string{"diff"}
    }
    return runGit(r.path, args...)
}

// GetCurrentHead returns the current commit hash
func (r *Repo) GetCurrentHead() (string, error) {
    return runGit(r.path, "rev-parse", "HEAD")
}

// GetChangedFiles returns list of changed files
func (r *Repo) GetChangedFiles(which string) ([]string, error) {
    var args []string
    if which == "staged" {
        args = []string{"diff", "--staged", "--name-only"}
    } else {
        args = []string{"diff", "--name-only"}
    }
    output, err := runGit(r.path, args...)
    if err != nil {
        return nil, err
    }
    if output == "" {
        return []string{}, nil
    }
    return strings.Split(strings.TrimSpace(output), "\n"), nil
}

func runGit(dir string, args ...string) (string, error) {
    cmd := exec.Command("git", args...)
    cmd.Dir = dir
    var out, errOut bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errOut
    err := cmd.Run()
    if err != nil {
        return "", fmt.Errorf("git %v: %s: %w", args, errOut.String(), err)
    }
    return out.String(), nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/git/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/git/git.go internal/git/git_test.go
git commit -m "feat: add Git operations wrapper"
```

---

## Task 5: CLI Foundation with Cobra

**Files:**
- Create: `cmd/vibediff/main.go`
- Create: `cmd/vibediff/root.go`
- Create: `cmd/vibediff/commit.go`

**Step 1: Create main entry point**

Create `cmd/vibediff/main.go`:

```go
package main

import (
    "github.com/vibediff/vibediff/cmd/vibediff"
)

func main() {
    vibediff.Execute()
}
```

**Step 2: Create root command**

Create `cmd/vibediff/root.go`:

```go
package vibediff

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "vibediff",
    Short: "Git-native versioning layer for prompts and AI outputs",
    Long: `VibeDiff captures AI prompts alongside the code changes they generate.

Use "vibediff commit" to attach a prompt to your commit.`,
}

// Execute runs the root command
func Execute() {
    cobra.CheckErr(rootCmd.Execute())
}
```

**Step 3: Install Cobra dependency**

Run: `go get github.com/spf13/cobra@latest`

**Step 4: Commit**

```bash
git add cmd/vibediff/main.go cmd/vibediff/root.go go.sum go.mod
git commit -m "feat: add CLI foundation with Cobra"
```

---

## Task 6: Commit Command

**Files:**
- Modify: `cmd/vibediff/root.go`
- Create: `cmd/vibediff/commit.go`

**Step 1: Create commit command**

Create `cmd/vibediff/commit.go`:

```go
package vibediff

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "github.com/vibediff/vibediff/internal/git"
    "github.com/vibediff/vibediff/internal/store"
)

var commitCmd = &cobra.Command{
    Use:   "commit <prompt>",
    Short: "Stage files, commit, and attach prompt",
    Long: `Stage all changes, create a commit, and attach the AI prompt
that generated these changes.

Example:
  vibediff commit "Add user authentication with JWT tokens"`,
    Args: cobra.ExactArgs(1),
    RunE: runCommit,
}

func runCommit(cmd *cobra.Command, args []string) error {
    prompt := args[0]

    // Get current directory
    wd, err := os.Getwd()
    if err != nil {
        return fmt.Errorf("failed to get working directory: %w", err)
    }

    // Initialize git repo
    repo, err := git.Init(wd)
    if err != nil {
        return fmt.Errorf("not in a git repository: %w", err)
    }

    // Stage all changes
    if err := repo.Add("."); err != nil {
        return fmt.Errorf("failed to stage files: %w", err)
    }

    // Get the diff
    diff, err := repo.Diff("staged")
    if err != nil {
        return fmt.Errorf("failed to get diff: %w", err)
    }

    // Get changed files
    files, err := repo.GetChangedFiles("staged")
    if err != nil {
        return fmt.Errorf("failed to get changed files: %w", err)
    }

    // Create commit
    commitMsg := fmt.Sprintf("vibe: %s", prompt)
    if err := repo.Commit(commitMsg); err != nil {
        return fmt.Errorf("failed to commit: %w", err)
    }

    // Get the commit hash
    hash, err := repo.GetCurrentHead()
    if err != nil {
        return fmt.Errorf("failed to get commit hash: %w", err)
    }

    // Store metadata
    vibeDir := filepath.Join(wd, ".vibediff")
    if err := os.MkdirAll(vibeDir, 0755); err != nil {
        return fmt.Errorf("failed to create .vibediff directory: %w", err)
    }

    dataFile := filepath.Join(vibeDir, "commits.json")
    js := store.NewJSONStore(dataFile)

    meta := store.CommitMetadata{
        Commit:    strings.TrimSpace(hash),
        Prompt:    prompt,
        Timestamp: time.Now().Format(time.RFC3339),
        Diff:      diff,
        Files:     files,
    }

    if err := js.Save(meta); err != nil {
        return fmt.Errorf("failed to save metadata: %w", err)
    }

    fmt.Printf("✅ Committed %s\n", hash[:7])
    fmt.Printf("📝 Prompt: %s\n", prompt)

    return nil
}

func init() {
    rootCmd.AddCommand(commitCmd)
}
```

**Step 2: Add missing imports**

Run: `go mod tidy`

**Step 3: Build and test**

Run: `go build -o vibediff ./cmd/vibediff`

**Step 4: Test manually**

Run: `./vibediff commit "test commit"` (in a git repo with changes)

**Step 5: Commit**

```bash
git add cmd/vibediff/commit.go
git commit -m "feat: add commit command"
```

---

## Task 7: Show Command

**Files:**
- Create: `cmd/vibediff/show.go`

**Step 1: Create show command**

Create `cmd/vibediff/show.go`:

```go
package vibediff

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"
    "github.com/vibediff/vibediff/internal/git"
    "github.com/vibediff/vibediff/internal/store"
)

var showCmd = &cobra.Command{
    Use:   "show [commit-ref]",
    Short: "Show prompt and diff for a commit",
    Long: `Display the AI prompt and resulting diff for a commit.

If no commit is specified, shows the most recent vibe commit.`,
    Args: cobra.MaximumNArgs(1),
    RunE: runShow,
}

func runShow(cmd *cobra.Command, args []string) error {
    wd, err := os.Getwd()
    if err != nil {
        return err
    }

    // Get commit ref
    ref := "HEAD"
    if len(args) > 0 {
        ref = args[0]
    }

    // Resolve to hash
    repo, err := git.Init(wd)
    if err != nil {
        return err
    }

    hash, err := runGit(repo, "rev-parse", ref)
    if err != nil {
        return fmt.Errorf("invalid commit ref: %w", err)
    }
    hash = strings.TrimSpace(hash)

    // Load metadata
    dataFile := filepath.Join(wd, ".vibediff", "commits.json")
    js := store.NewJSONStore(dataFile)

    meta, err := js.Load(hash)
    if err != nil {
        return fmt.Errorf("no vibe metadata found for commit %s: %w", hash[:7], err)
    }

    // Display
    fmt.Printf("Commit: %s\n", hash[:7])
    fmt.Printf("When:   %s\n", meta.Timestamp)
    fmt.Printf("\n📝 Prompt:\n")
    fmt.Printf("  %s\n\n", meta.Prompt)

    if len(meta.Files) > 0 {
        fmt.Printf("📁 Files: %s\n\n", strings.Join(meta.Files, ", "))
    }

    fmt.Printf("📊 Diff:\n%s\n", meta.Diff)

    return nil
}

func init() {
    rootCmd.AddCommand(showCmd)
}
```

**Step 2: Build**

Run: `go build -o vibediff ./cmd/vibediff`

**Step 3: Test**

Run: `./vibediff show`

**Step 4: Commit**

```bash
git add cmd/vibediff/show.go
git commit -m "feat: add show command"
```

---

## Task 8: Log Command

**Files:**
- Create: `cmd/vibediff/log.go`

**Step 1: Create log command**

Create `cmd/vibediff/log.go`:

```go
package vibediff

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"
    "github.com/vibediff/vibediff/internal/store"
)

var logCmd = &cobra.Command{
    Use:   "log",
    Short: "Show all vibe commits with their prompts",
    Long: `Display a log of all commits with attached vibe prompts.`,
    Args: cobra.NoArgs,
    RunE: runLog,
}

func runLog(cmd *cobra.Command, args []string) error {
    wd, err := os.Getwd()
    if err != nil {
        return err
    }

    dataFile := filepath.Join(wd, ".vibediff", "commits.json")
    js := store.NewJSONStore(dataFile)

    commits, err := js.List()
    if err != nil {
        return fmt.Errorf("failed to load commits: %w", err)
    }

    if len(commits) == 0 {
        fmt.Println("No vibe commits found.")
        return nil
    }

    fmt.Printf("Found %d vibe commit(s):\n\n", len(commits))

    for i, meta := range commits {
        fmt.Printf("%d. %s", len(commits)-i, meta.Commit[:7])
        fmt.Printf("  %s\n", meta.Timestamp)
        fmt.Printf("  📝 %s\n\n", meta.Prompt)
    }

    return nil
}

func init() {
    rootCmd.AddCommand(logCmd)
}
```

**Step 2: Build**

Run: `go build -o vibediff ./cmd/vibediff`

**Step 3: Test**

Run: `./vibediff log`

**Step 4: Commit**

```bash
git add cmd/vibediff/log.go
git commit -m "feat: add log command"
```

---

## Task 9: Export HTML Command

**Files:**
- Create: `internal/ui/html.go`
- Create: `cmd/vibediff/export.go`

**Step 1: Create HTML generator**

Create `internal/ui/html.go`:

```go
package ui

import (
    "html/template"
    "os"

    "github.com/vibediff/vibediff/internal/store"
)

var htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>VibeDiff - {{.Commit}}</title>
    <style>
        body { font-family: system-ui; max-width: 900px; margin: 40px auto; padding: 0 20px; }
        .header { border-bottom: 1px solid #eee; padding-bottom: 20px; margin-bottom: 20px; }
        .prompt { background: #f0f7ff; padding: 15px; border-radius: 8px; border-left: 4px solid #0066cc; }
        .diff { background: #f6f8fa; padding: 15px; border-radius: 8px; font-family: monospace; white-space: pre-wrap; font-size: 13px; }
        .add { color: #22863a; }
        .remove { color: #b31d28; }
        h1 { color: #333; }
        .meta { color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🔮 VibeDiff</h1>
        <div class="meta">Commit: <code>{{.Commit}}</code> | {{.Timestamp}}</div>
    </div>
    <div class="prompt">
        <strong>📝 Prompt:</strong><br>
        {{.Prompt}}
    </div>
    {{if .Files}}
    <p><strong>📁 Files:</strong> {{range .Files}}{{.}} {{end}}</p>
    {{end}}
    <h2>📊 Diff</h2>
    <div class="diff">{{.Diff}}</div>
</body>
</html>
`

// ExportHTML writes metadata to an HTML file
func ExportHTML(meta store.CommitMetadata, path string) error {
    tmpl, err := template.New("vibe").Parse(htmlTemplate)
    if err != nil {
        return err
    }

    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()

    return tmpl.Execute(f, meta)
}
```

**Step 2: Create export command**

Create `cmd/vibediff/export.go`:

```go
package vibediff

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"
    "github.com/vibediff/vibediff/internal/git"
    "github.com/vibediff/vibediff/internal/store"
    "github.com/vibediff/vibediff/internal/ui"
)

var exportHTML bool

var exportCmd = &cobra.Command{
    Use:   "export [commit-ref]",
    Short: "Export vibe commit to HTML",
    Long: `Generate a shareable HTML file showing the prompt and diff.

If no commit is specified, exports the most recent vibe commit.`,
    Args: cobra.MaximumNArgs(1),
    RunE: runExport,
}

func runExport(cmd *cobra.Command, args []string) error {
    wd, err := os.Getwd()
    if err != nil {
        return err
    }

    ref := "HEAD"
    if len(args) > 0 {
        ref = args[0]
    }

    repo, err := git.Init(wd)
    if err != nil {
        return err
    }

    hash, err := runGitString(repo, "rev-parse", ref)
    if err != nil {
        return fmt.Errorf("invalid commit ref: %w", err)
    }
    hash = strings.TrimSpace(hash)

    dataFile := filepath.Join(wd, ".vibediff", "commits.json")
    js := store.NewJSONStore(dataFile)

    meta, err := js.Load(hash)
    if err != nil {
        return fmt.Errorf("no vibe metadata found: %w", err)
    }

    outputPath := fmt.Sprintf("vibe-%s.html", hash[:7])
    if err := ui.ExportHTML(meta, outputPath); err != nil {
        return fmt.Errorf("failed to export HTML: %w", err)
    }

    fmt.Printf("✅ Exported to %s\n", outputPath)

    return nil
}

func runGitString(repo *git.Repo, args ...string) (string, error) {
    // This needs to be added to git.go or use runGit directly
    // For now, using internal helper
}

func init() {
    exportCmd.Flags().BoolVar(&exportHTML, "html", true, "Export as HTML")
    rootCmd.AddCommand(exportCmd)
}
```

**Step 3: Add helper to git.go**

Modify `internal/git/git.go`, add:

```go
// RunCommand executes a git command and returns output
func (r *Repo) RunCommand(args ...string) (string, error) {
    return runGit(r.path, args...)
}
```

**Step 4: Update export.go to use new method**

Replace `runGitString` with:

```go
hash, err := repo.RunCommand("rev-parse", ref)
```

**Step 5: Build and test**

Run: `go build -o vibediff ./cmd/vibediff && ./vibediff export`

**Step 6: Commit**

```bash
git add internal/ui/html.go cmd/vibediff/export.go internal/git/git.go
git commit -m "feat: add export command with HTML output"
```

---

## Task 10: Documentation and Examples

**Files:**
- Modify: `README.md`
- Create: `examples/basic/README.md`
- Create: `examples/basic/main.go`

**Step 1: Update README**

Create comprehensive `README.md`:

```markdown
# VibeDiff

Git-native versioning layer for prompts and AI outputs.

## What is VibeDiff?

VibeDiff captures AI prompts alongside the code changes they generate.
Never lose track of "how did I get this code?" again.

## Installation

\`\`\`bash
go install github.com/vibediff/vibediff/cmd/vibediff@latest
\`\`\`

## Quick Start

\`\`\`bash
# Make some code changes with AI help
vim auth.go

# Commit with the prompt that generated it
vibediff commit "Add JWT authentication with refresh tokens"

# View the prompt + diff anytime
vibediff show

# See all vibe commits
vibediff log

# Share with teammates
vibediff export
\`\`\`

## Commands

| Command | Description |
|---------|-------------|
| `vibediff commit <prompt>` | Stage + commit with prompt attached |
| `vibediff show [ref]` | Display prompt + diff for a commit |
| `vibediff log` | Show all vibe commits |
| `vibediff export [ref]` | Generate shareable HTML |

## License

Apache 2.0
```

**Step 2: Create example project**

Create `examples/basic/README.md`:

```markdown
# VibeDiff Example

A simple example demonstrating VibeDiff usage.
```

Create `examples/basic/main.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, VibeDiff!")
}
```

**Step 3: Commit**

```bash
git add README.md examples/
git commit -m "docs: add README and examples"
```

---

## Task 11: Build and Release

**Files:**
- Create: `Makefile`
- Create: `.github/workflows/release.yml`

**Step 1: Create Makefile**

Create `Makefile`:

```makefile
.PHONY: build test install clean

build:
	go build -o bin/vibediff ./cmd/vibediff

test:
	go test -v ./...

install:
	go install ./cmd/vibediff

clean:
	rm -rf bin/

# Build for all platforms
release:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/vibediff-linux-amd64 ./cmd/vibediff
	GOOS=darwin GOARCH=amd64 go build -o bin/vibediff-darwin-amd64 ./cmd/vibediff
	GOOS=darwin GOARCH=arm64 go build -o bin/vibediff-darwin-arm64 ./cmd/vibediff
	GOOS=windows GOARCH=amd64 go build -o bin/vibediff-windows-amd64.exe ./cmd/vibediff
```

**Step 2: Create GitHub Actions workflow**

Create `.github/workflows/release.yml`:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: make test
      - run: make release
      - uses: softprops/action-gh-release@v1
        with:
          files: bin/*
```

**Step 3: Commit**

```bash
git add Makefile .github/
git commit -m "ci: add build and release automation"
```

---

## Task 12: Final Testing

**Step 1: Run all tests**

Run: `make test`

**Step 2: Build binary**

Run: `make build`

**Step 3: Manual integration test**

```bash
cd /tmp
mkdir test-vibe && cd test-vibe
git init
echo "initial" > file.txt
git add . && git commit -m "init"

echo "changed by AI" > file.txt
../../vibediff/bin/vibediff commit "Update file with new content"

../../vibediff/bin/vibediff log
../../vibediff/bin/vibediff show
../../vibediff/bin/vibediff export
```

**Step 4: Verify HTML output**

Open the generated HTML file in browser

**Step 5: Final commit**

```bash
git add .
git commit -m "chore: final polish and testing"
```

---

## Post-Implementation

- Tag first release: `git tag v0.1.0`
- Push to GitHub
- Create GitHub release
- Announce on HackerNews / Reddit
