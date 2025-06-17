# Tablux

A TUI file/text visualizer for JSON, CSV, and other formats, built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea).

## Features

- Interactive visualization of JSON, JSONL, and CSV files
- Collapsible JSON tree view for easy navigation
- CSV table view with column sorting and visibility control
- File format auto-detection
- Syntax highlighting
- Keyboard-driven navigation

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/tablux.git
cd tablux

# Build using Make
make build

# Or build manually
go build -o tablux cmd/tablux/main.go

# Install to GOPATH/bin
make install
# Or
go install github.com/yourusername/tablux/cmd/tablux@latest
```

## Usage

```bash
# Interactive mode (default)
tablux path/to/file.json
tablux path/to/file.jsonl
tablux path/to/file.csv

# Non-interactive mode (output rendered content to stdout)
tablux path/to/file.json --no-interactive
```

The `--no-interactive` flag is useful for:
- Testing the rendering without a TTY
- Piping output to other commands
- Debugging formatting issues

## Keyboard Controls

### Common Controls
- `q`: Quit
- `Ctrl+C`: Quit

### JSON/JSONL Viewer Controls
- `↑`/`k`: Navigate up
- `↓`/`j`: Navigate down
- `Enter`/`Space`: Expand/collapse current node
- `c`: Collapse all nodes (great for large JSONs)
- `e`: Expand all nodes

### CSV Viewer Controls
- `↑`/`k`: Navigate up
- `↓`/`j`: Navigate down
- `←`/`h`: Navigate left
- `→`/`l`: Navigate right
- `v`: Toggle column visibility (columns stay visible as collapsed indicators)
- `s`: Sort by current column (toggle ascending/descending)

## Project Structure

- `cmd/tablux`: Main application
- `pkg/loader`: File loading functionality
- `pkg/parser`: Format detection and parsing
- `pkg/model`: Data models
- `pkg/ui`: UI components
- `pkg/utils`: Utility functions

## Development

Currently implemented:
- JSON/JSONL file format support with interactive tree view
- CSV file support with sortable columns and column visibility toggle
- Non-interactive mode for debugging and piping

Planned features:
- Enhanced search functionality across different formats
- Support for more file formats (YAML, XML, etc.)
- Advanced filtering options

## License

MIT