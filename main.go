package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

type FileInfo struct {
	Name  string
	IsDir bool
	Path  string
}

type Dir struct {
	Name string
	Path string
}

var directories []string

func dirListing(w http.ResponseWriter, r *http.Request, fsPath string) {
	f, err := os.Open(fsPath)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var fileInfos []FileInfo
	for _, file := range files {
		fileInfos = append(fileInfos, FileInfo{
			Name:  file.Name(),
			IsDir: file.IsDir(),
			Path:  filepath.Join(r.URL.Path, file.Name()),
		})
	}

	// Sort directories first
	sort.Slice(fileInfos, func(i, j int) bool {
		if fileInfos[i].IsDir && !fileInfos[j].IsDir {
			return true
		} else if !fileInfos[i].IsDir && fileInfos[j].IsDir {
			return false
		}
		return fileInfos[i].Name < fileInfos[j].Name
	})

	tmpl := `
<!DOCTYPE html>
<html>
<head>
<title>Index of {{.Path}}</title>
<style>
	body { font-family: sans-serif; }
	.dir { font-weight: bold; }
</style>
</head>
<body>
<h1>Index of {{.Path}}</h1>
<ul>
{{range .Files}}
	<li class="{{if .IsDir}}dir{{end}}">
		<a href="{{.Path}}{{if .IsDir}}/{{end}}">{{.Name}}{{if .IsDir}}/{{end}}</a>
	</li>
{{end}}
</ul>
</body>
</html>
`
	t, err := template.New("dir").Parse(tmpl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Path  string
		Files []FileInfo
	}{
		Path:  r.URL.Path,
		Files: fileInfos,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, data)
}

func fileHandler(rootDirs []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Match the request path to the correct root directory
		for _, dir := range rootDirs {
			fsPath := filepath.Join(dir, r.URL.Path)
			info, err := os.Stat(fsPath)
			if os.IsNotExist(err) {
				continue
			} else if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if info.IsDir() {
				dirListing(w, r, fsPath)
				return
			}

			// Serve the file with correct MIME type
			http.ServeFile(w, r, fsPath)
			return
		}

		http.NotFound(w, r)
	}
}

func main() {
	port := flag.String("port", "8000", "Port to serve HTTP on")
	flag.Parse()

	directories = flag.Args()
	if len(directories) == 0 {
		// Default to current directory if none provided
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		directories = append(directories, dir)
	} else {
		// Validate directories
		for _, d := range directories {
			info, err := os.Stat(d)
			if err != nil || !info.IsDir() {
				log.Fatalf("Invalid directory: %s", d)
			}
		}
	}

	http.HandleFunc("/", fileHandler(directories))

	log.Printf("Serving directories %v on HTTP port: %s\n", directories, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
