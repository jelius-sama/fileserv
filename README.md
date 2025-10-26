# Modern File Server

A multi-directory file server written in Go with a modern, responsive UI that supports both light and dark themes.

## Features

- ✅ Serve multiple directories simultaneously
- ✅ Modern, responsive UI with dark/light theme support
- ✅ Clean project structure with separated concerns
- ✅ Directory navigation and file browsing
- ✅ Human-readable file sizes
- ✅ Production-ready error handling and logging
- ✅ Proper MIME type detection

## Project Structure

```
fileserv/
├── main.go                          # Entry point
├── go.mod                           # Go module definition
├── internal/
│   ├── models/
│   │   └── types.go                # Data models
│   ├── server/
│   │   └── validator.go            # Directory validation
│   ├── handler/
│   │   └── handler.go              # HTTP request handling
│   └── template/
│       └── template.go             # HTML templates
└── README.md
```

## Installation

```bash
# Clone the repository
git clone https://github.com/jelius-sama/fileserv
cd fileserv

# Build the binary
./build.sh
```

## Usage

### Basic Usage (Current Directory)

```bash
./bin/fileserv
```

This will serve the current directory on port 8000.

### Specify Port

```bash
./bin/fileserv -port 3000
```

### Serve Multiple Directories

```bash
./bin/fileserv -port 8080 -dir /path/to/dir1 /path/to/dir2 /path/to/dir3
```

When serving multiple directories:
- The root page (/) displays all available directories
- Click on any directory to browse its contents
- Use the dropdown to switch between directories
- Each directory is accessible at `/<directory-name>/`

## Command Line Options

- `-port`: Port to serve HTTP on (default: 8000)

## Examples

```bash
# Serve current directory on port 8000
./bin/fileserv

# Serve on custom port
./bin/fileserv -port 9000

# Serve multiple directories
./bin/fileserv /home/user/documents /home/user/downloads /var/www

# Serve with specific port and directories
./bin/fileserv -port 3000 ~/projects ~/music ~/photos
```

## UI Features

- **Responsive Design**: Works on desktop, tablet, and mobile devices
- **Theme Support**: Automatically adapts to system light/dark mode preference
- **Directory Cards**: Visual cards for selecting root directories
- **File Icons**: Visual distinction between files and directories
- **File Sizes**: Human-readable file sizes (B, KB, MB, GB, etc.)
- **Breadcrumb Navigation**: Easy navigation with back links
- **Directory Switcher**: Quick dropdown to switch between served directories

## Security Considerations

This is a simple file server intended for local or trusted network use. For production use over the internet:

1. Add authentication
2. Implement rate limiting
3. Add HTTPS support
4. Implement access control lists
5. Add request logging and monitoring
6. Consider using reverse proxy (nginx, caddy)

## Development

To run in development mode:

```bash
go run ./ -port 8000 /path/to/test/dir
```

## License

[MIT License](./LICENSE).
