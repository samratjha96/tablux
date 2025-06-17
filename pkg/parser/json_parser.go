package parser

import (
	"encoding/json"
	"fmt"

	"tablux/pkg/model"
)

// JSONParser parses JSON data into a tree structure
type JSONParser struct {
	// Configuration options could go here
}

// NewJSONParser creates a new JSON parser
func NewJSONParser() *JSONParser {
	return &JSONParser{}
}

// Parse parses JSON data into a tree structure
func (p *JSONParser) Parse(data []byte) (*model.JSONNode, error) {
	var v interface{}

	// Parse JSON
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Create the root node
	rootNode := model.NewJSONNode("root", v, nil)
	return rootNode, nil
}

// ParseJSONL parses JSONL data (one JSON object per line)
func (p *JSONParser) ParseJSONL(data []byte) ([]*model.JSONNode, error) {
	var nodes []*model.JSONNode

	// Split by lines and parse each line separately
	lines := splitLines(data)
	for i, line := range lines {
		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		var v interface{}
		err := json.Unmarshal(line, &v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line %d: %w", i+1, err)
		}

		// Create a node for this line
		key := fmt.Sprintf("[%d]", i)
		node := model.NewJSONNode(key, v, nil)
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// Helper function to split data into lines
func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			if i > start {
				lines = append(lines, data[start:i])
			}
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
