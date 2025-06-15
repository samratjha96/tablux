# Tablux

A TUI file/text visualizer for JSON, CSV, and other formats, built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea).

## Features

- Interactive visualization of JSON, JSONL, and CSV files
- Collapsible JSON tree view for easy navigation
- File format auto-detection
- Syntax highlighting
- Keyboard-driven navigation

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/tablux.git
cd tablux

# Build the application
go build -o tablux cmd/tablux/main.go

# Or install it
go install github.com/yourusername/tablux/cmd/tablux@latest
```

## Usage

```bash
tablux path/to/file.json
tablux path/to/file.jsonl
tablux path/to/file.csv
```

## Keyboard Controls

- `↑`/`k`: Navigate up
- `↓`/`j`: Navigate down
- `Enter`/`Space`: Expand/collapse current node
- `c`: Collapse all nodes (great for large JSONs)
- `e`: Expand all nodes
- `q`: Quit

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

Planned features:
- CSV viewer with column manipulation
- Enhanced navigation and search
- More file formats

## License

MIT