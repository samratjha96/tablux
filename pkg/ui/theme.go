package ui

import "github.com/charmbracelet/lipgloss"

// Common theme colors
const (
	// Base colors
	PrimaryColor    = "#4B6BEF"
	SecondaryColor  = "#5A5AA0"
	HighlightColor  = "#5555CC"
	TextColor       = "#FFFFFF"
	MutedTextColor  = "#AAAAAA"
	BackgroundColor = "#333333"

	// JSON node colors
	KeyColor     = "#88AAFF"
	StringColor  = "#7CFC00"
	NumberColor  = "#FFD700"
	BoolColor    = "#FF9F5F"
	NullColor    = "#FF5F5F"
	BracketColor = "#F8F8F2"

	// UI symbols
	ExpandedIndicator  = "▼ "
	CollapsedIndicator = "► "
	CollapsedColumn    = "│"
	SortAscIndicator   = " ▲"
	SortDescIndicator  = " ▼"
)

// Default dimensions and spacing
const (
	DefaultCellPadding    = 1
	DefaultColumnMaxWidth = 30
	CollapsedColumnWidth  = 2
)

// CreateStyle creates a standard cell style with common settings
func CreateStyle(fg, bg string, bold bool) lipgloss.Style {
	style := lipgloss.NewStyle().
		PaddingLeft(DefaultCellPadding).
		PaddingRight(DefaultCellPadding).
		AlignHorizontal(lipgloss.Left)

	if fg != "" {
		style = style.Foreground(lipgloss.Color(fg))
	}

	if bg != "" {
		style = style.Background(lipgloss.Color(bg))
	}

	if bold {
		style = style.Bold(true)
	}

	return style
}

// Common styles that can be shared across viewers
var (
	// Basic styles
	HeaderStyle = CreateStyle(TextColor, PrimaryColor, true)
	CellStyle   = CreateStyle("", "", false)

	// Selection styles
	SelectedRowStyle  = CreateStyle("", BackgroundColor, false)
	SelectedColStyle  = CreateStyle(TextColor, SecondaryColor, false)
	SelectedCellStyle = CreateStyle(TextColor, HighlightColor, true)

	// Collapsed styles
	CollapsedHeaderStyle = CreateStyle(TextColor, "#777777", true).Width(CollapsedColumnWidth)
	CollapsedCellStyle   = CreateStyle(MutedTextColor, BackgroundColor, false).Width(CollapsedColumnWidth)

	// JSON specific styles
	KeyStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color(KeyColor))
	StringStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color(StringColor))
	NumberStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color(NumberColor))
	BoolStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color(BoolColor))
	NullStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color(NullColor))
	BracketStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(BracketColor))
	SelectedNodeStyle = lipgloss.NewStyle().Background(lipgloss.Color(BackgroundColor))
	SeparatorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(MutedTextColor))
)

// Tree symbols for JSON viewer
var TreeSymbols = map[string]string{
	"pipe":      "│ ",
	"tee":       "├─",
	"last":      "└─",
	"expanded":  ExpandedIndicator,
	"collapsed": CollapsedIndicator,
	"empty":     "  ",
}
