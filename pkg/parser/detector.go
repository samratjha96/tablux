package parser

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
)

// FileFormat represents supported file formats
type FileFormat int

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

// detectByContent analyzes the file content to determine its format
func detectByContent(data []byte) FileFormat {
	// Trim whitespace
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return FormatUnknown
	}

	// Check if it looks like JSON (starts with { or [)
	if trimmed[0] == '{' || trimmed[0] == '[' {
		var js interface{}
		if json.Unmarshal(trimmed, &js) == nil {
			return FormatJSON
		}
	}

	// Check if it looks like JSONL (multiple lines, each a valid JSON object)
	lines := bytes.Split(trimmed, []byte("\n"))
	if len(lines) > 1 {
		jsonlCount := 0
		for _, line := range lines[:min(10, len(lines))] {
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
		if jsonlCount > len(lines)/2 {
			return FormatJSONL
		}
	}

	// Check if it looks like CSV
	r := csv.NewReader(bytes.NewReader(trimmed))
	r.FieldsPerRecord = -1 // Allow variable number of fields
	records, err := r.ReadAll()
	if err == nil && len(records) > 0 {
		// CSV detection is a bit tricky, but let's assume it's CSV if:
		// 1. We can parse it as CSV
		// 2. There are multiple rows
		// 3. All rows have at least 2 fields
		if len(records) > 1 {
			validRows := 0
			for _, record := range records {
				if len(record) >= 2 {
					validRows++
				}
			}
			if validRows > len(records)/2 {
				return FormatCSV
			}
		}
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