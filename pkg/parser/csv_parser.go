package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// CSVData represents parsed CSV data
type CSVData struct {
	Headers []string
	Rows    [][]string
	// Track column widths for display formatting
	ColumnWidths []int
	// Track column visibility
	ColumnVisibility []bool
	// Track sorting order
	SortColumn int
	SortAsc    bool
}

// NewCSVData creates a new empty CSVData structure
func NewCSVData() *CSVData {
	return &CSVData{
		Headers:          []string{},
		Rows:             [][]string{},
		ColumnWidths:     []int{},
		ColumnVisibility: []bool{},
		SortColumn:       -1, // No sorting by default
	}
}

// CSVParser parses CSV data
type CSVParser struct {
	// Configuration options (delimiter etc.)
	Comma                rune
	Comment              rune
	UseFirstLineAsHeader bool
}

// NewCSVParser creates a new CSV parser with default settings
func NewCSVParser() *CSVParser {
	return &CSVParser{
		Comma:                ',',
		Comment:              '#',
		UseFirstLineAsHeader: true,
	}
}

// Parse parses CSV data from a byte array
func (p *CSVParser) Parse(data []byte) (*CSVData, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))
	reader.Comma = p.Comma
	reader.Comment = p.Comment

	csvData := NewCSVData()

	// Read all records at once
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) == 0 {
		return csvData, nil
	}

	// Extract headers if configured
	startRow := 0
	if p.UseFirstLineAsHeader {
		csvData.Headers = records[0]
		startRow = 1
	} else {
		// Generate default headers (Column 1, Column 2, etc.)
		csvData.Headers = make([]string, len(records[0]))
		for i := range csvData.Headers {
			csvData.Headers[i] = fmt.Sprintf("Column %d", i+1)
		}
	}

	// Initialize column visibility (all visible by default)
	csvData.ColumnVisibility = make([]bool, len(csvData.Headers))
	for i := range csvData.ColumnVisibility {
		csvData.ColumnVisibility[i] = true
	}

	// Add data rows
	for i := startRow; i < len(records); i++ {
		csvData.Rows = append(csvData.Rows, records[i])
	}

	// Calculate column widths for display formatting
	csvData.calculateColumnWidths()

	return csvData, nil
}

// ParseStream parses CSV data from an io.Reader
func (p *CSVParser) ParseStream(reader io.Reader) (*CSVData, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = p.Comma
	csvReader.Comment = p.Comment

	csvData := NewCSVData()

	// Read headers if configured
	if p.UseFirstLineAsHeader {
		headers, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				return csvData, nil
			}
			return nil, fmt.Errorf("failed to read CSV headers: %w", err)
		}
		csvData.Headers = headers

		// Initialize column visibility
		csvData.ColumnVisibility = make([]bool, len(headers))
		for i := range csvData.ColumnVisibility {
			csvData.ColumnVisibility[i] = true
		}
	}

	// Read data rows
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}

		// If we haven't read headers yet, generate them based on the first row
		if len(csvData.Headers) == 0 {
			csvData.Headers = make([]string, len(record))
			for i := range csvData.Headers {
				csvData.Headers[i] = fmt.Sprintf("Column %d", i+1)
			}

			// Initialize column visibility
			csvData.ColumnVisibility = make([]bool, len(csvData.Headers))
			for i := range csvData.ColumnVisibility {
				csvData.ColumnVisibility[i] = true
			}
		}

		csvData.Rows = append(csvData.Rows, record)
	}

	// Calculate column widths
	csvData.calculateColumnWidths()

	return csvData, nil
}

// calculateColumnWidths updates the ColumnWidths field based on the current data
func (c *CSVData) calculateColumnWidths() {
	// Initialize column widths based on headers
	c.ColumnWidths = make([]int, len(c.Headers))
	for i, header := range c.Headers {
		c.ColumnWidths[i] = len(header)
	}

	// Update column widths based on data
	for _, row := range c.Rows {
		for i, cell := range row {
			if i < len(c.ColumnWidths) && len(cell) > c.ColumnWidths[i] {
				c.ColumnWidths[i] = len(cell)
			}
		}
	}
}

// ToggleColumnVisibility toggles the visibility of a column
func (c *CSVData) ToggleColumnVisibility(colIndex int) {
	if colIndex >= 0 && colIndex < len(c.ColumnVisibility) {
		c.ColumnVisibility[colIndex] = !c.ColumnVisibility[colIndex]
	}
}

// GetVisibleColumns returns indices of visible columns
func (c *CSVData) GetVisibleColumns() []int {
	var visible []int
	for i, isVisible := range c.ColumnVisibility {
		if isVisible {
			visible = append(visible, i)
		}
	}
	return visible
}

// IsColumnVisible returns whether a column is visible
func (c *CSVData) IsColumnVisible(colIndex int) bool {
	if colIndex >= 0 && colIndex < len(c.ColumnVisibility) {
		return c.ColumnVisibility[colIndex]
	}
	return false
}

// SortByColumn sorts the data by the specified column
func (c *CSVData) SortByColumn(colIndex int, ascending bool) {
	if colIndex < 0 || colIndex >= len(c.Headers) {
		return
	}

	c.SortColumn = colIndex
	c.SortAsc = ascending

	// Simple string comparison sort
	if ascending {
		for i := 0; i < len(c.Rows)-1; i++ {
			for j := i + 1; j < len(c.Rows); j++ {
				if colIndex < len(c.Rows[i]) && colIndex < len(c.Rows[j]) &&
					c.Rows[i][colIndex] > c.Rows[j][colIndex] {
					c.Rows[i], c.Rows[j] = c.Rows[j], c.Rows[i]
				}
			}
		}
	} else {
		for i := 0; i < len(c.Rows)-1; i++ {
			for j := i + 1; j < len(c.Rows); j++ {
				if colIndex < len(c.Rows[i]) && colIndex < len(c.Rows[j]) &&
					c.Rows[i][colIndex] < c.Rows[j][colIndex] {
					c.Rows[i], c.Rows[j] = c.Rows[j], c.Rows[i]
				}
			}
		}
	}
}
