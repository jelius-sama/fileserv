package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fileserv/internal/models"
	"fileserv/internal/template"
)

// FileServer handles file serving and directory listings
type FileServer struct {
	directories []models.Directory
}

// NewFileServer creates a new file server instance
func NewFileServer(dirs []models.Directory) *FileServer {
	return &FileServer{
		directories: dirs,
	}
}

// HandleRequest handles incoming HTTP requests
func (fs *FileServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Root path - show directory selector
	if path == "/" {
		fs.showRootListing(w, r)
		return
	}

	// Try to match the path to a directory
	for _, dir := range fs.directories {
		// Check if the path starts with the directory name
		prefix := "/" + dir.Name
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			// Remove the prefix to get the relative path
			relPath := strings.TrimPrefix(path, prefix)
			if relPath == "" {
				relPath = "/"
			}

			fsPath := filepath.Join(dir.Path, filepath.Clean(relPath))
			fs.serveFromDirectory(w, r, dir, fsPath, relPath)
			return
		}
	}

	http.NotFound(w, r)
}

// showRootListing shows the root directory selector
func (fs *FileServer) showRootListing(w http.ResponseWriter, r *http.Request) {
	data := models.PageData{
		CurrentPath: "/",
		Files:       nil,
		Directories: fs.directories,
		IsRoot:      true,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := template.RenderListing(w, data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// serveFromDirectory serves files from a specific directory
func (fs *FileServer) serveFromDirectory(w http.ResponseWriter, r *http.Request, dir models.Directory, fsPath, relPath string) {
	info, err := os.Stat(fsPath)
	if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		log.Printf("Error stating file %s: %v", fsPath, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		fs.showDirectoryListing(w, r, dir, fsPath, relPath)
		return
	}

	http.ServeFile(w, r, fsPath)
}

// showDirectoryListing shows the contents of a directory
func (fs *FileServer) showDirectoryListing(w http.ResponseWriter, r *http.Request, dir models.Directory, fsPath, relPath string) {
	f, err := os.Open(fsPath)
	if err != nil {
		log.Printf("Error opening directory %s: %v", fsPath, err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	defer f.Close()

	entries, err := f.Readdir(-1)
	if err != nil {
		log.Printf("Error reading directory %s: %v", fsPath, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var fileInfos []models.FileInfo
	for _, entry := range entries {
		// Build the URL path
		urlPath := "/" + dir.Name + filepath.ToSlash(filepath.Join(relPath, entry.Name()))

		fileInfos = append(fileInfos, models.FileInfo{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Path:  urlPath,
			Size:  entry.Size(),
		})
	}

	// Sort: directories first, then by name
	sort.Slice(fileInfos, func(i, j int) bool {
		if fileInfos[i].IsDir != fileInfos[j].IsDir {
			return fileInfos[i].IsDir
		}
		return strings.ToLower(fileInfos[i].Name) < strings.ToLower(fileInfos[j].Name)
	})

	data := models.PageData{
		CurrentPath: "/" + dir.Name + relPath,
		Files:       fileInfos,
		Directories: fs.directories,
		IsRoot:      false,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := template.RenderListing(w, data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
