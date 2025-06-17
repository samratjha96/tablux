package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"tablux/pkg/parser"
)

var (
	// CSV viewer styles (using theme constants)
	// Using theme-defined styles directly
	headerStyle       = HeaderStyle
	cellStyle         = CellStyle
	selectedRowStyle  = SelectedRowStyle
	selectedColStyle  = SelectedColStyle
	selectedCellStyle = SelectedCellStyle

	// Collapsed/hidden column style
	collapsedColHeaderStyle = CollapsedHeaderStyle
	collapsedColStyle       = CollapsedCellStyle

	// Table styles
	separatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(MutedTextColor))

	tableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(MutedTextColor))

	// Column separators
	columnSeparator = " "

	// Column sorting indicators using theme constants
	sortAscIndicator  = SortAscIndicator
	sortDescIndicator = SortDescIndicator

	// Defaults using theme constants
	defaultCellPadding    = DefaultCellPadding
	defaultColumnMaxWidth = DefaultColumnMaxWidth
	collapsedColumnWidth  = CollapsedColumnWidth

	// Collapsed column indicator
	collapsedIndicator = CollapsedColumn
)

// CSVViewer displays a CSV table
type CSVViewer struct {
	data           *parser.CSVData
	cursorRow      int
	cursorCol      int
	viewportX      int
	viewportY      int
	viewportWidth  int
	viewportHeight int
	columnMaxWidth int   // Max width of a column before truncation
	columnWidths   []int // Pre-calculated widths for columns
}

// NewCSVViewer creates a new CSV viewer
func NewCSVViewer(data *parser.CSVData) *CSVViewer {
	viewer := &CSVViewer{
		data:           data,
		cursorRow:      0,
		cursorCol:      0,
		viewportX:      0,
		viewportY:      0,
		columnMaxWidth: defaultColumnMaxWidth,
	}

	// Pre-calculate column widths
	viewer.calculateColumnWidths()

	return viewer
}

// calculateColumnWidths pre-calculates optimal widths for all columns
func (v *CSVViewer) calculateColumnWidths() {
	colCount := len(v.data.Headers)
	v.columnWidths = make([]int, colCount)

	// Initialize with header widths
	for i, header := range v.data.Headers {
		// Add space for sort indicators
		width := len(header) + 4 // Add padding and space for sort indicators
		v.columnWidths[i] = width
	}

	// Update with data cell widths if needed
	for _, row := range v.data.Rows {
		for i, cell := range row {
			if i < colCount {
				cellWidth := len(cell) + 2
				if cellWidth > v.columnWidths[i] {
					v.columnWidths[i] = cellWidth
				}
			}
		}
	}

	// Cap all widths to maximum and ensure minimum width
	for i, width := range v.columnWidths {
		if width > v.columnMaxWidth {
			v.columnWidths[i] = v.columnMaxWidth
		} else if width < 10 {
			v.columnWidths[i] = 10
		}

		// Ensure even widths for better alignment
		if v.columnWidths[i]%2 == 1 {
			v.columnWidths[i]++
		}
	}
}

// SetViewport sets the viewport dimensions
func (v *CSVViewer) SetViewport(width, height int) {
	v.viewportWidth = width
	v.viewportHeight = height
	v.ensureCursorVisible()
}

// Move cursor methods
func (v *CSVViewer) MoveUp() {
	if v.cursorRow > 0 {
		v.cursorRow--
		v.ensureCursorVisible()
	}
}

func (v *CSVViewer) MoveDown() {
	if v.cursorRow < len(v.data.Rows) {
		v.cursorRow++
		v.ensureCursorVisible()
	}
}

func (v *CSVViewer) MoveLeft() {
	if v.cursorCol > 0 {
		v.cursorCol--
		v.ensureCursorVisible()
	}
}

func (v *CSVViewer) MoveRight() {
	if v.cursorCol < len(v.data.Headers)-1 {
		v.cursorCol++
		v.ensureCursorVisible()
	}
}

// ToggleColumnVisibility toggles visibility of the current column
func (v *CSVViewer) ToggleColumnVisibility() {
	v.data.ToggleColumnVisibility(v.cursorCol)
}

// SortByCurrentColumn sorts by the current column
func (v *CSVViewer) SortByCurrentColumn() {
	ascending := true
	if v.data.SortColumn == v.cursorCol {
		// Toggle order if already sorting by this column
		ascending = !v.data.SortAsc
	}
	v.data.SortByColumn(v.cursorCol, ascending)
}

// ensureCursorVisible adjusts viewport to keep cursor in view
func (v *CSVViewer) ensureCursorVisible() {
	// Adjust vertical viewport
	if v.cursorRow < v.viewportY {
		v.viewportY = v.cursorRow
	} else if v.cursorRow >= v.viewportY+v.viewportHeight-1 { // -1 for header
		v.viewportY = v.cursorRow - v.viewportHeight + 2
	}
}

// Render renders the CSV viewer
func (v *CSVViewer) Render() string {
	var table strings.Builder

	// Create header
	headers := v.createHeaderRow()
	table.WriteString(headers)
	table.WriteString("\n")

	// Calculate visible rows
	startRow := v.viewportY
	endRow := min(startRow+v.viewportHeight-2, len(v.data.Rows)) // -2 for header and spacing

	// Create data rows
	for rowIdx := startRow; rowIdx < endRow; rowIdx++ {
		dataRow := v.createDataRow(rowIdx)
		table.WriteString(dataRow)
		table.WriteString("\n")
	}

	// Apply table border
	result := tableStyle.Render(table.String())
	return result
}

// createHeaderRow generates header row with consistent formatting
func (v *CSVViewer) createHeaderRow() string {
	var cells []string

	// Create each header cell
	for i, header := range v.data.Headers {
		// Handle hidden columns
		if !v.data.ColumnVisibility[i] {
			// Create collapsed indicator
			style := collapsedColHeaderStyle
			if i == v.cursorCol {
				style = style.Background(lipgloss.Color(HighlightColor))
			}
			cells = append(cells, style.Render(collapsedIndicator))
			continue
		}

		// Get cell content
		content := header
		width := v.columnWidths[i]

		// Handle sort indicators
		if i == v.data.SortColumn {
			if v.data.SortAsc {
				content += sortAscIndicator
			} else {
				content += sortDescIndicator
			}
		}

		// Truncate if needed
		if len(content) > width-2 {
			content = content[:width-5] + "..."
		}

		// Apply styling with fixed width
		style := headerStyle.Copy().Width(width)
		if i == v.cursorCol {
			style = style.Background(lipgloss.Color(HighlightColor))
		}

		// Render cell with exact width
		cells = append(cells, style.Render(content))
	}

	return strings.Join(cells, "")
}

// createDataRow generates a single data row with consistent formatting
func (v *CSVViewer) createDataRow(rowIdx int) string {
	var cells []string
	row := v.data.Rows[rowIdx]

	// Create each data cell
	for i := range v.data.Headers {
		// Handle hidden columns
		if !v.data.ColumnVisibility[i] {
			// Create collapsed indicator
			style := collapsedColStyle
			if i == v.cursorCol && rowIdx == v.cursorRow {
				style = style.Background(lipgloss.Color(HighlightColor))
			} else if i == v.cursorCol {
				style = style.Background(lipgloss.Color(SecondaryColor))
			} else if rowIdx == v.cursorRow {
				style = style.Background(lipgloss.Color(BackgroundColor))
			}
			cells = append(cells, style.Render(collapsedIndicator))
			continue
		}

		// Get cell content
		var content string
		if i < len(row) {
			content = row[i]
		}
		width := v.columnWidths[i]

		// Truncate if needed
		if len(content) > width-2 {
			content = content[:width-5] + "..."
		}

		// Select styling based on cursor position
		var style lipgloss.Style
		if rowIdx == v.cursorRow && i == v.cursorCol {
			style = selectedCellStyle
		} else if rowIdx == v.cursorRow {
			style = selectedRowStyle
		} else if i == v.cursorCol {
			style = selectedColStyle
		} else {
			style = cellStyle
		}

		// Apply same width as headers for consistent alignment
		style = style.Copy().Width(width)
		cells = append(cells, style.Render(content))
	}

	return strings.Join(cells, "")
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetColumnWidth returns the width of a specific column
func (v *CSVViewer) GetColumnWidth(colIndex int) int {
	if colIndex >= 0 && colIndex < len(v.columnWidths) {
		return v.columnWidths[colIndex]
	}
	return 0
}
