package template

import (
	"fmt"
	"html/template"
	"io"

	"fileserv/internal/models"
)

var tmpl = template.Must(template.New("listing").Funcs(template.FuncMap{
	"formatSize": formatSize,
}).Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{if .IsRoot}}File Server{{else}}{{.CurrentPath}}{{end}}</title>
    <style>
        :root {
            --bg-primary: #ffffff;
            --bg-secondary: #f8f9fa;
            --bg-hover: #e9ecef;
            --text-primary: #212529;
            --text-secondary: #6c757d;
            --border-color: #dee2e6;
            --accent-color: #0d6efd;
            --accent-hover: #0a58ca;
            --shadow: rgba(0, 0, 0, 0.1);
        }

        @media (prefers-color-scheme: dark) {
            :root {
                --bg-primary: #1a1a1a;
                --bg-secondary: #2d2d2d;
                --bg-hover: #3a3a3a;
                --text-primary: #e9ecef;
                --text-secondary: #adb5bd;
                --border-color: #495057;
                --accent-color: #4dabf7;
                --accent-hover: #339af0;
                --shadow: rgba(0, 0, 0, 0.3);
            }
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: var(--bg-primary);
            color: var(--text-primary);
            line-height: 1.6;
            min-height: 100vh;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 2rem 1rem;
        }

        header {
            margin-bottom: 2rem;
            padding-bottom: 1rem;
            border-bottom: 2px solid var(--border-color);
        }

        h1 {
            font-size: 1.75rem;
            font-weight: 600;
            margin-bottom: 0.5rem;
            color: var(--text-primary);
        }

        .breadcrumb {
            color: var(--text-secondary);
            font-size: 0.9rem;
        }

        .directory-selector {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
            gap: 1rem;
            margin-top: 2rem;
        }

        .directory-card {
            background: var(--bg-secondary);
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 1.5rem;
            text-decoration: none;
            color: var(--text-primary);
            transition: all 0.2s ease;
            box-shadow: 0 2px 4px var(--shadow);
        }

        .directory-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 8px var(--shadow);
            border-color: var(--accent-color);
        }

        .directory-card-icon {
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
        }

        .directory-card-name {
            font-weight: 600;
            font-size: 1.1rem;
            margin-bottom: 0.25rem;
        }

        .directory-card-path {
            font-size: 0.85rem;
            color: var(--text-secondary);
            word-break: break-all;
        }

        .file-list {
            background: var(--bg-secondary);
            border: 1px solid var(--border-color);
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 4px var(--shadow);
        }

        .file-item {
            display: flex;
            align-items: center;
            padding: 1rem 1.5rem;
            text-decoration: none;
            color: var(--text-primary);
            border-bottom: 1px solid var(--border-color);
            transition: background 0.15s ease;
        }

        .file-item:last-child {
            border-bottom: none;
        }

        .file-item:hover {
            background: var(--bg-hover);
        }

        .file-icon {
            font-size: 1.5rem;
            margin-right: 1rem;
            min-width: 2rem;
            text-align: center;
        }

        .file-info {
            flex: 1;
            min-width: 0;
        }

        .file-name {
            font-weight: 500;
            margin-bottom: 0.2rem;
            word-break: break-all;
        }

        .file-meta {
            font-size: 0.85rem;
            color: var(--text-secondary);
        }

        .file-size {
            margin-left: auto;
            padding-left: 1rem;
            color: var(--text-secondary);
            font-size: 0.9rem;
            white-space: nowrap;
        }

        .directory-nav {
            background: var(--bg-secondary);
            padding: 0.75rem 1.5rem;
            border-radius: 8px;
            margin-bottom: 1rem;
            border: 1px solid var(--border-color);
        }

        .directory-nav select {
            background: var(--bg-primary);
            color: var(--text-primary);
            border: 1px solid var(--border-color);
            padding: 0.5rem 1rem;
            border-radius: 6px;
            font-size: 0.95rem;
            cursor: pointer;
            transition: all 0.2s ease;
        }

        .directory-nav select:hover {
            border-color: var(--accent-color);
        }

        .directory-nav select:focus {
            outline: none;
            border-color: var(--accent-color);
            box-shadow: 0 0 0 3px rgba(13, 110, 253, 0.1);
        }

        .back-link {
            display: inline-flex;
            align-items: center;
            color: var(--accent-color);
            text-decoration: none;
            margin-bottom: 1rem;
            font-weight: 500;
            transition: color 0.2s ease;
        }

        .back-link:hover {
            color: var(--accent-hover);
        }

        .empty-state {
            text-align: center;
            padding: 3rem 1rem;
            color: var(--text-secondary);
        }

        @media (max-width: 768px) {
            .container {
                padding: 1rem;
            }

            h1 {
                font-size: 1.5rem;
            }

            .file-item {
                padding: 0.75rem 1rem;
            }

            .file-size {
                display: none;
            }

            .directory-selector {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            {{if .IsRoot}}
                <h1>📁 File Server</h1>
                <div class="breadcrumb">Select a directory to browse</div>
            {{else}}
                <a href="/" class="back-link">← Back to all directories</a>
                <h1>{{.CurrentPath}}</h1>
            {{end}}
        </header>

        {{if .IsRoot}}
            <div class="directory-selector">
                {{range .Directories}}
                <a href="/{{.Name}}" class="directory-card">
                    <div class="directory-card-icon">📂</div>
                    <div class="directory-card-name">{{.Name}}</div>
                    <div class="directory-card-path">{{.Path}}</div>
                </a>
                {{end}}
            </div>
        {{else}}
            {{if gt (len .Directories) 1}}
            <div class="directory-nav">
                <label for="dir-select">Switch directory: </label>
                <select id="dir-select" onchange="window.location.href='/' + this.value">
                    {{range .Directories}}
                    <option value="{{.Name}}">{{.Name}}</option>
                    {{end}}
                </select>
            </div>
            {{end}}

            {{if .Files}}
            <div class="file-list">
                {{range .Files}}
                <a href="{{.Path}}" class="file-item">
                    <div class="file-icon">{{if .IsDir}}📁{{else}}📄{{end}}</div>
                    <div class="file-info">
                        <div class="file-name">{{.Name}}</div>
                        <div class="file-meta">{{if .IsDir}}Directory{{else}}File{{end}}</div>
                    </div>
                    {{if not .IsDir}}
                    <div class="file-size">{{formatSize .Size}}</div>
                    {{end}}
                </a>
                {{end}}
            </div>
            {{else}}
            <div class="empty-state">
                <p>📭 This directory is empty</p>
            </div>
            {{end}}
        {{end}}
    </div>
</body>
</html>
`))

// RenderListing renders the directory listing template
func RenderListing(w io.Writer, data models.PageData) error {
	return tmpl.Execute(w, data)
}

// formatSize formats file size in human-readable format
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
