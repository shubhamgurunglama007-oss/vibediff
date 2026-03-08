// Package ui generates user-facing output formats for VibeDiff.
// It currently supports HTML export for sharing commits.
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
        body { font-family: system-ui, -apple-system, sans-serif; max-width: 900px; margin: 40px auto; padding: 0 20px; line-height: 1.6; }
        .header { border-bottom: 1px solid #eee; padding-bottom: 20px; margin-bottom: 20px; }
        .prompt { background: #f0f7ff; padding: 15px; border-radius: 8px; border-left: 4px solid #0066cc; }
        .details { background: #fff8e6; padding: 15px; border-radius: 8px; border-left: 4px solid #ff9800; margin-top: 15px; }
        .diff { background: #f6f8fa; padding: 15px; border-radius: 8px; font-family: monospace; white-space: pre-wrap; font-size: 13px; overflow-x: auto; }
        .add { color: #22863a; }
        .remove { color: #b31d28; }
        h1 { color: #333; }
        h2 { margin-top: 30px; }
        .meta { color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>VibeDiff</h1>
        <div class="meta">Commit: <code>{{.Commit}}</code> | {{.Timestamp}}</div>
    </div>
    <div class="prompt">
        <strong>Prompt:</strong><br>
        {{.Prompt}}
    </div>
    {{if .Details}}
    <div class="details">
        <strong>Details:</strong><br>
        {{.Details}}
    </div>
    {{end}}
    {{if .Files}}
    <p><strong>Files:</strong> {{range .Files}}{{.}} {{end}}</p>
    {{end}}
    <h2>Diff</h2>
    <div class="diff">{{.Diff}}</div>
</body>
</html>
`

// ExportHTML writes metadata to an HTML file.
// The template automatically escapes HTML entities to prevent XSS.
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
