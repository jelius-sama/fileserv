package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fileserv/internal/handler"
	"fileserv/internal/server"
)

const version = "1.0.0"

// expandTilde expands ~ to the user's home directory
func expandTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	if path == "~" {
		home, err := os.UserHomeDir()
		if err == nil {
			return home
		}
	}
	return path
}

func main() {
	var directories []string
	var port string
	var showVersion bool
	var showHelp bool

	// Manual flag parsing to handle mixed flags and arguments
	args := os.Args[1:]

	i := 0
	for i < len(args) {
		arg := args[i]

		switch arg {
		case "-port", "--port":
			if i+1 < len(args) {
				port = args[i+1]
				i += 2
			} else {
				i++
			}
		case "-dir", "--dir":
			i++
			// Collect all following arguments until we hit another flag
			for i < len(args) && !strings.HasPrefix(args[i], "-") {
				// Check if this argument contains commas (comma-separated list)
				if strings.Contains(args[i], ",") {
					// Split by comma and add each part
					parts := strings.Split(args[i], ",")
					for _, part := range parts {
						trimmed := strings.TrimSpace(part)
						if trimmed != "" {
							directories = append(directories, expandTilde(trimmed))
						}
					}
				} else {
					// Single directory (tilde already expanded by shell for space-separated args)
					directories = append(directories, args[i])
				}
				i++
			}
		case "-version", "--version":
			showVersion = true
			i++
		case "-help", "--help", "-h":
			showHelp = true
			i++
		default:
			// If it doesn't start with -, treat it as a directory
			if !strings.HasPrefix(arg, "-") {
				directories = append(directories, args[i])
			}
			i++
		}
	}

	// Set default port if not specified
	if port == "" {
		port = "8000"
	}

	// Handle version flag
	if showVersion {
		fmt.Printf("fileserv version %s\n", version)
		os.Exit(0)
	}

	// Handle help flag
	if showHelp {
		printHelp()
		os.Exit(0)
	}

	// If no directories specified, use current directory
	if len(directories) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		directories = append(directories, dir)
	}

	// Validate directories
	validDirs, err := server.ValidateDirectories(directories)
	if err != nil {
		log.Fatal(err)
	}

	// Create file server
	fs := handler.NewFileServer(validDirs)

	// Setup routes
	http.HandleFunc("/", fs.HandleRequest)

	log.Printf("Serving directories %v on HTTP port: %s\n", validDirs, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func printHelp() {
	fmt.Println(`
  ╔══════════════════════════════════════════════════════════════════╗
  ║                                                                  ║
  ║   ███████╗██╗██╗     ███████╗███████╗███████╗██████╗ ██╗   ██╗   ║
  ║   ██╔════╝██║██║     ██╔════╝██╔════╝██╔════╝██╔══██╗██║   ██║   ║
  ║   █████╗  ██║██║     █████╗  ███████╗█████╗  ██████╔╝██║   ██║   ║
  ║   ██╔══╝  ██║██║     ██╔══╝  ╚════██║██╔══╝  ██╔══██╗╚██╗ ██╔╝   ║
  ║   ██║     ██║███████╗███████╗███████║███████╗██║  ██║ ╚████╔╝    ║
  ║   ╚═╝     ╚═╝╚══════╝╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝  ╚═══╝     ║
  ║                                                                  ║
  ║              A Modern Multi-Directory File Server                ║
  ╚══════════════════════════════════════════════════════════════════╝
`)
	fmt.Printf("  Version: %s\n\n", version)

	fmt.Println("  USAGE")
	fmt.Println("  ─────────────────────────────────────────────────────────────")
	fmt.Println("    fileserv [options] [directories...]")
	fmt.Println()

	fmt.Println("  OPTIONS")
	fmt.Println("  ─────────────────────────────────────────────────────────────")
	fmt.Println("    -port <number>")
	fmt.Println("        Port to serve HTTP on (default: 8000)")
	fmt.Println()
	fmt.Println("    -dir <paths>")
	fmt.Println("        Directories to serve (space or comma-separated)")
	fmt.Println("        Can also pass directories as arguments after flags")
	fmt.Println("        If not specified, serves the current directory")
	fmt.Println()
	fmt.Println("    -version")
	fmt.Println("        Show version information")
	fmt.Println()
	fmt.Println("    -help")
	fmt.Println("        Show this help message")
	fmt.Println()

	fmt.Println("  EXAMPLES")
	fmt.Println("  ─────────────────────────────────────────────────────────────")
	fmt.Println("    # Serve current directory on default port")
	fmt.Println("    $ fileserv")
	fmt.Println()
	fmt.Println("    # Serve on custom port")
	fmt.Println("    $ fileserv -port 3000")
	fmt.Println()
	fmt.Println("    # Serve multiple directories (space-separated)")
	fmt.Println("    $ fileserv -dir ~/Documents ~/Downloads ~/Pictures")
	fmt.Println()
	fmt.Println("    # Serve multiple directories (comma-separated)")
	fmt.Println("    $ fileserv -dir ~/Documents,~/Downloads,~/Pictures")
	fmt.Println()
	fmt.Println("    # Serve with custom port (flags in any order)")
	fmt.Println("    $ fileserv -port 3000 -dir ~/Documents ~/Downloads")
	fmt.Println("    $ fileserv -dir ~/Documents ~/Downloads -port 3000")
	fmt.Println()
	fmt.Println("    # Serve directories as standalone arguments")
	fmt.Println("    $ fileserv ~/Documents ~/Downloads -port 3000")
	fmt.Println()
	fmt.Println("  ─────────────────────────────────────────────────────────────")
	fmt.Println()
}
