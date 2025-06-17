package parser

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
)

// FileFormat represents supported file formats
type FileFormat int

// File format constants
const (
	// FormatUnknown is the default format when detection fails
	FormatUnknown FileFormat = iota
	// FormatJSON represents JSON format
	FormatJSON
	// FormatJSONL represents JSON Lines format (one JSON object per line)
	FormatJSONL
	// FormatCSV represents CSV format
	FormatCSV
)

// File format string representations for consistent usage
const (
	TypeJSON    = "json"
	TypeJSONL   = "jsonl"
	TypeCSV     = "csv"
	TypeUnknown = "unknown"
)

// String returns the string representation of the file format
func (f FileFormat) String() string {
	switch f {
	case FormatJSON:
		return "JSON"
	case FormatJSONL:
		return "JSONL"
	case FormatCSV:
		return "CSV"
	default:
		return "Unknown"
	}
}

// ToTypeString returns the lowercase type string representation needed by the UI
func (f FileFormat) ToTypeString() string {
	switch f {
	case FormatJSON:
		return TypeJSON
	case FormatJSONL:
		return TypeJSONL
	case FormatCSV:
		return TypeCSV
	default:
		return TypeUnknown
	}
}

// DetectFormat tries to determine the format of the file contents
func DetectFormat(data []byte, extension string) FileFormat {
	// Try to detect by extension first
	switch strings.ToLower(extension) {
	case ".json":
		return FormatJSON
	case ".jsonl":
		return FormatJSONL
	case ".csv":
		return FormatCSV
	}

	// If extension doesn't conclusively determine format, inspect the content
	return detectByContent(data)
}

// DetectFileType returns a string representation of the file type
// This is a helper function used by the UI to determine which viewer to use
func DetectFileType(data []byte) string {
	return detectByContent(data).ToTypeString()
}

// detectJSONFormat attempts to detect JSON format from data
func detectJSONFormat(data []byte) (FileFormat, bool) {
	// Check if it looks like standard JSON (starts with { or [)
	if len(data) > 0 && (data[0] == '{' || data[0] == '[') {
		var js interface{}
		if json.Unmarshal(data, &js) == nil {
			return FormatJSON, true
		}
	}
	return FormatUnknown, false
}

// detectJSONLFormat attempts to detect JSONL format from data
func detectJSONLFormat(lines [][]byte) (FileFormat, bool) {
	if len(lines) <= 1 {
		return FormatUnknown, false
	}

	// Check a sample of lines (up to 10)
	sampleSize := min(10, len(lines))
	jsonlCount := 0

	for _, line := range lines[:sampleSize] {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if line[0] == '{' {
			var js interface{}
			if json.Unmarshal(line, &js) == nil {
				jsonlCount++
			}
		}
	}

	// If most (>50%) of the tested lines are valid JSON objects, assume JSONL
	if jsonlCount > 0 && jsonlCount >= sampleSize/2 {
		return FormatJSONL, true
	}

	return FormatUnknown, false
}

// detectCSVFormat attempts to detect CSV format from data
func detectCSVFormat(data []byte) (FileFormat, bool) {
	r := csv.NewReader(bytes.NewReader(data))
	r.FieldsPerRecord = -1 // Allow variable number of fields
	records, err := r.ReadAll()

	if err != nil || len(records) <= 1 {
		return FormatUnknown, false
	}

	// Count valid rows (rows with at least 2 fields)
	validRows := 0
	for _, record := range records {
		if len(record) >= 2 {
			validRows++
		}
	}

	// If most (>50%) rows are valid, assume it's a CSV
	if validRows > len(records)/2 {
		return FormatCSV, true
	}

	return FormatUnknown, false
}

// detectByContent analyzes the file content to determine its format
func detectByContent(data []byte) FileFormat {
	// Trim whitespace
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return FormatUnknown
	}

	// Try to detect JSON first (fastest check)
	if format, ok := detectJSONFormat(trimmed); ok {
		return format
	}

	// Split into lines for JSONL detection
	lines := bytes.Split(trimmed, []byte("\n"))

	// Try to detect JSONL
	if format, ok := detectJSONLFormat(lines); ok {
		return format
	}

	// Try to detect CSV last (can be expensive for large files)
	if format, ok := detectCSVFormat(trimmed); ok {
		return format
	}

	return FormatUnknown
}

// Helper function for min value
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
