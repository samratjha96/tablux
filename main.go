package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"tablux/pkg/parser"
	"tablux/pkg/ui"
)

// UI constants
const (
	AppName  = "TABLUX"
	AppTitle = " TABLUX "

	// File type constants
	TypeJSON  = "json"
	TypeJSONL = "jsonl"
	TypeCSV   = "csv"

	// Viewport padding
	HeaderFooterSpace = 4 // Space needed for header and footer
	CSVBorderSpace    = 6 // Extra space needed for CSV borders and padding

	// Default sizes for non-interactive mode
	DefaultHeight = 30
	DefaultWidth  = 100

	// Input methods
	InputStdin = "<stdin>"
)

// Colors
var (
	PrimaryColor = lipgloss.Color("#7D56F4")
	TextColor    = lipgloss.Color("#FAFAFA")
	ErrorColor   = lipgloss.Color("#FF5555")
)

// Define styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextColor).
			Background(PrimaryColor).
			PaddingLeft(2).
			PaddingRight(2)

	infoStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Italic(true).
			PaddingLeft(1)
)

// Model represents the application state
type Model struct {
	title      string
	filePath   string
	width      int
	height     int
	jsonViewer *ui.JSONViewer
	csvViewer  *ui.CSVViewer
	viewerType string
	isLoading  bool
	errorMsg   string
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		loadSourceCmd(m.filePath),
	)
}

// FileLoadedMsg is sent when a file is loaded
type FileLoadedMsg struct {
	viewerType string
	jsonViewer *ui.JSONViewer
	csvViewer  *ui.CSVViewer
	error      error
}

// parseFile parses data and returns appropriate viewer based on file type
// If a specific format is provided, it will use that instead of auto-detection
func parseFile(data []byte, forcedFormat string) (string, *ui.JSONViewer, *ui.CSVViewer, error) {
	fileType := forcedFormat

	// Auto-detect format if not forced
	if fileType == "" {
		fileType = parser.DetectFileType(data)
	}

	switch fileType {
	case TypeJSON, TypeJSONL:
		// Parse JSON data
		jsonParser := parser.NewJSONParser()
		root, err := jsonParser.Parse(data)
		if err != nil {
			return "", nil, nil, err
		}

		// Create JSON viewer
		viewer := ui.NewJSONViewer(root)
		return fileType, viewer, nil, nil

	case TypeCSV:
		// Parse CSV data
		csvParser := parser.NewCSVParser()
		csvData, err := csvParser.Parse(data)
		if err != nil {
			return "", nil, nil, err
		}

		// Create CSV viewer
		viewer := ui.NewCSVViewer(csvData)
		return fileType, nil, viewer, nil

	default:
		return "", nil, nil, fmt.Errorf("unsupported file type")
	}
}

// loadSourceCmd loads data from a file or stdin and returns the appropriate viewer
func loadSourceCmd(source string) tea.Cmd {
	return func() tea.Msg {
		// Get the format flag value
		formatFlag := ""
		flag.Visit(func(f *flag.Flag) {
			if f.Name == "format" {
				formatFlag = f.Value.String()
			}
		})

		// Load data from source
		data, err := readDataFromSource(source)
		if err != nil {
			return FileLoadedMsg{error: err}
		}

		// Parse data with optional format
		fileType, jsonViewer, csvViewer, err := parseFile(data, formatFlag)
		if err != nil {
			return FileLoadedMsg{error: err}
		}

		return FileLoadedMsg{
			viewerType: fileType,
			jsonViewer: jsonViewer,
			csvViewer:  csvViewer,
		}
	}
}

// handleJSONKeyMsg processes key presses for JSON viewer
func (m *Model) handleJSONKeyMsg(key string) {
	if m.jsonViewer == nil {
		return
	}

	switch key {
	case "up":
		m.jsonViewer.MoveUp()
	case "down":
		m.jsonViewer.MoveDown()
	case "enter", " ":
		m.jsonViewer.ToggleNode()
	}
}

// handleCSVKeyMsg processes key presses for CSV viewer
func (m *Model) handleCSVKeyMsg(key string) {
	if m.csvViewer == nil {
		return
	}

	switch key {
	case "up":
		m.csvViewer.MoveUp()
	case "down":
		m.csvViewer.MoveDown()
	case "left":
		m.csvViewer.MoveLeft()
	case "right":
		m.csvViewer.MoveRight()
	case "enter", " ":
		m.csvViewer.ToggleColumnVisibility()
	case "s":
		m.csvViewer.SortByCurrentColumn()
	}
}

// Update handles messages and user input
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if key == "q" || key == "ctrl+c" {
			return m, tea.Quit
		}

		// Handle viewer-specific keys
		switch m.viewerType {
		case TypeJSON, TypeJSONL:
			m.handleJSONKeyMsg(key)
		case TypeCSV:
			m.handleCSVKeyMsg(key)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update viewers with new size
		if m.jsonViewer != nil {
			m.jsonViewer.SetViewportHeight(m.height - HeaderFooterSpace) // Subtract space for header/footer
		}
		if m.csvViewer != nil {
			m.csvViewer.SetViewport(m.width-HeaderFooterSpace, m.height-CSVBorderSpace) // Subtract space for borders and header/footer
		}

	case FileLoadedMsg:
		m.isLoading = false
		if msg.error != nil {
			m.errorMsg = fmt.Sprintf("Error: %v", msg.error)
			return m, nil
		}

		m.viewerType = msg.viewerType
		if msg.viewerType == TypeJSON || msg.viewerType == TypeJSONL {
			m.jsonViewer = msg.jsonViewer
			m.jsonViewer.SetViewportHeight(m.height - HeaderFooterSpace)
		} else if msg.viewerType == TypeCSV {
			m.csvViewer = msg.csvViewer
			m.csvViewer.SetViewport(m.width-HeaderFooterSpace, m.height-CSVBorderSpace)
		}
	}

	return m, nil
}

// renderError renders an error message
func renderError(msg string) string {
	return fmt.Sprintf("%s\n\n%s",
		titleStyle.Render(AppTitle),
		lipgloss.NewStyle().Foreground(ErrorColor).Render(msg))
}

// renderLoading renders a loading message
func renderLoading(path string) string {
	return fmt.Sprintf("%s\n\nLoading %s...",
		titleStyle.Render(AppTitle),
		path)
}

// getControlsForViewer returns help text based on viewer type
func getControlsForViewer(viewerType string) string {
	switch viewerType {
	case TypeJSON, TypeJSONL:
		return infoStyle.Render("↑/↓: Navigate | Space/Enter: Toggle | q: Quit")
	case TypeCSV:
		return infoStyle.Render("↑/↓/←/→: Navigate | Space/Enter: Toggle visibility | s: Sort | q: Quit")
	default:
		return infoStyle.Render("q: Quit")
	}
}

func (m Model) View() string {
	if m.errorMsg != "" {
		return renderError(m.errorMsg)
	}

	if m.isLoading {
		return renderLoading(m.filePath)
	}

	// Create header with title and file info
	header := lipgloss.JoinHorizontal(lipgloss.Top,
		titleStyle.Render(AppTitle),
		lipgloss.NewStyle().PaddingLeft(2).Render(fmt.Sprintf("File: %s | Type: %s", m.filePath, m.viewerType)))

	// Create content based on viewer type
	var content string
	switch m.viewerType {
	case TypeJSON, TypeJSONL:
		if m.jsonViewer != nil {
			content = m.jsonViewer.Render()
		}
	case TypeCSV:
		if m.csvViewer != nil {
			content = m.csvViewer.Render()
		}
	default:
		content = "No content to display"
	}

	// Get controls for current viewer
	controls := getControlsForViewer(m.viewerType)

	// Combine all elements
	return fmt.Sprintf("%s\n\n%s\n\n%s", header, content, controls)
}

// testCSVViewer tests the CSV viewer alignment
func testCSVViewer() {
	// Load sample CSV file
	data, err := os.ReadFile("test/sample.csv")
	if err != nil {
		fmt.Printf("Error reading sample CSV: %v\n", err)
		os.Exit(1)
	}

	// Parse CSV data
	parser := parser.NewCSVParser()
	csvData, err := parser.Parse(data)
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		os.Exit(1)
	}

	// Create CSV viewer
	viewer := ui.NewCSVViewer(csvData)
	viewer.SetViewport(DefaultWidth-HeaderFooterSpace, DefaultHeight-CSVBorderSpace)

	// Render and output result
	result := viewer.Render()
	fmt.Println("CSV Viewer Output:")
	fmt.Println(result)

	// Output raw header widths for debugging
	fmt.Println("\nHeader widths:")
	for i, header := range csvData.Headers {
		fmt.Printf("%d. '%s' - Width: %d\n", i+1, header, viewer.GetColumnWidth(i))
	}
}

// readDataFromSource reads data from either a file or stdin
func readDataFromSource(source string) ([]byte, error) {
	// Read from stdin if specified
	if source == InputStdin {
		return io.ReadAll(os.Stdin)
	}

	// Otherwise read from file
	return os.ReadFile(source)
}

// runNonInteractiveMode shows content without TUI
func runNonInteractiveMode(source string) {
	// Get the format flag value
	formatFlag := ""
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "format" {
			formatFlag = f.Value.String()
		}
	})

	data, err := readDataFromSource(source)
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}

	fileType, jsonViewer, csvViewer, err := parseFile(data, formatFlag)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	switch fileType {
	case TypeJSON, TypeJSONL:
		jsonViewer.SetViewportHeight(DefaultHeight - HeaderFooterSpace)
		fmt.Println(jsonViewer.Render())

	case TypeCSV:
		csvViewer.SetViewport(DefaultWidth-HeaderFooterSpace, DefaultHeight-CSVBorderSpace)
		fmt.Println(csvViewer.Render())
	}
}

// printHelp prints usage information
func printHelp() {
	fmt.Printf("Usage: %s [OPTIONS]\n\n", os.Args[0])
	fmt.Println("A TUI file/text visualizer for JSON, CSV, and other formats.")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  # View a file interactively")
	fmt.Println("  tablux --file data.json")
	fmt.Println("\n  # Process stdin input")
	fmt.Println("  cat data.csv | tablux")
	fmt.Println("\n  # Force specific format")
	fmt.Println("  cat data.txt | tablux --format json")
	fmt.Println("\n  # Output to stdout (non-interactive)")
	fmt.Println("  tablux --file data.json --no-interactive")
	fmt.Println("\nKeyboard controls:")
	fmt.Println("  q, Ctrl+C: Quit")
	fmt.Println("  ↑/↓: Navigate")
	fmt.Println("  Space/Enter: Toggle expand/collapse (JSON) or column visibility (CSV)")
	fmt.Println("  s: Sort column (CSV only)")
}

func main() {
	// Parse command-line flags
	filePath := flag.String("file", "", "Path to the file to open (omit to use stdin)")
	noInteractive := flag.Bool("no-interactive", false, "Run in non-interactive mode")
	testCSV := flag.Bool("test-csv", false, "Run CSV viewer test")
	format := flag.String("format", "", "Force a specific format: json, jsonl, or csv")
	help := flag.Bool("help", false, "Show usage information")
	flag.Parse()

	// Show help if requested
	if *help {
		printHelp()
		return
	}

	// Handle CSV test mode
	if *testCSV {
		testCSVViewer()
		return
	}

	// Validate format if provided
	if *format != "" && *format != TypeJSON && *format != TypeJSONL && *format != TypeCSV {
		fmt.Printf("Invalid format: %s. Use json, jsonl, or csv.\n", *format)
		os.Exit(1)
	}

	// Determine input source
	source := *filePath

	// Check if we should use stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Data is being piped to stdin
		source = InputStdin
	} else if source == "" {
		// No input source provided
		fmt.Println("No input provided. Use --file flag or pipe data to stdin.")
		fmt.Println("Run with --help for usage information.")
		os.Exit(1)
	}

	// Run in non-interactive mode if requested
	if *noInteractive {
		runNonInteractiveMode(source)
		return
	}

	// Create initial model
	m := Model{
		title:     AppName,
		filePath:  source,
		isLoading: true,
	}

	// Run interactive mode
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
