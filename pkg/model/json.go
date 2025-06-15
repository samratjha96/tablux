package model

import (
	"fmt"
)

// NodeType represents the type of a JSON node
type NodeType int

const (
	// NodeObject represents a JSON object
	NodeObject NodeType = iota
	// NodeArray represents a JSON array
	NodeArray
	// NodeString represents a string value
	NodeString
	// NodeNumber represents a numeric value
	NodeNumber
	// NodeBoolean represents a boolean value
	NodeBoolean
	// NodeNull represents a null value
	NodeNull
)

// JSONNode represents a node in the JSON tree
type JSONNode struct {
	Key      string
	Type     NodeType
	Value    interface{}
	Children []*JSONNode
	Parent   *JSONNode
	Expanded bool
	Path     string
}

// NewJSONNode creates a new JSON node
func NewJSONNode(key string, value interface{}, parent *JSONNode) *JSONNode {
	node := &JSONNode{
		Key:      key,
		Value:    value,
		Parent:   parent,
		Children: []*JSONNode{},
		Expanded: true,
	}

	// Set path
	if parent != nil {
		if parent.Path == "" {
			node.Path = key
		} else if key == "" { // Array element
			node.Path = parent.Path
		} else {
			node.Path = parent.Path + "." + key
		}
	} else {
		node.Path = ""
	}

	// Determine node type and create children for complex types
	switch v := value.(type) {
	case map[string]interface{}:
		node.Type = NodeObject
		for k, val := range v {
			child := NewJSONNode(k, val, node)
			node.Children = append(node.Children, child)
		}
	case []interface{}:
		node.Type = NodeArray
		for _, val := range v {
			child := NewJSONNode("", val, node)
			node.Children = append(node.Children, child)
		}
	case string:
		node.Type = NodeString
	case float64, int, int64:
		node.Type = NodeNumber
	case bool:
		node.Type = NodeBoolean
	case nil:
		node.Type = NodeNull
	}

	return node
}

// Toggle expands or collapses a node
func (n *JSONNode) Toggle() {
	if n.Type == NodeObject || n.Type == NodeArray {
		n.Expanded = !n.Expanded
	}
}

// IsLeaf returns true if the node is a leaf node (has no children)
func (n *JSONNode) IsLeaf() bool {
	return len(n.Children) == 0
}

// HasChildren returns true if the node has children
func (n *JSONNode) HasChildren() bool {
	return len(n.Children) > 0
}

// TypeString returns a string representation of the node type
func (n *JSONNode) TypeString() string {
	switch n.Type {
	case NodeObject:
		return "object"
	case NodeArray:
		return "array"
	case NodeString:
		return "string"
	case NodeNumber:
		return "number"
	case NodeBoolean:
		return "boolean"
	case NodeNull:
		return "null"
	default:
		return "unknown"
	}
}

// GetDisplayValue returns a string representation of the node's value
func (n *JSONNode) GetDisplayValue() string {
	switch n.Type {
	case NodeObject:
		if n.HasChildren() {
			return "{...}"
		}
		return "{}"
	case NodeArray:
		if n.HasChildren() {
			return "[...]"
		}
		return "[]"
	case NodeString:
		return "\"" + n.Value.(string) + "\""
	case NodeNumber:
		return InterfaceToString(n.Value)
	case NodeBoolean:
		if n.Value.(bool) {
			return "true"
		}
		return "false"
	case NodeNull:
		return "null"
	default:
		return ""
	}
}

// InterfaceToString converts an interface{} to a string
func InterfaceToString(v interface{}) string {
	if v == nil {
		return "null"
	}
	
	switch val := v.(type) {
	case string:
		return val
	case float64:
		// Format float without trailing zeros
		s := FloatToString(val)
		return s
	case int:
		return IntToString(val)
	case int64:
		return Int64ToString(val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return "?"
	}
}

// Helper functions for string conversion
func FloatToString(f float64) string {
	// Check if it's an integer value
	if f == float64(int(f)) {
		return IntToString(int(f))
	}
	return ToString(f)
}

func IntToString(i int) string {
	return ToString(i)
}

func Int64ToString(i int64) string {
	return ToString(i)
}

func ToString(v interface{}) string {
	return stringifyValue(v)
}

func stringifyValue(v interface{}) string {
	// Use sprintf for general case
	return toString(v)
}

func toString(v interface{}) string {
	if v == nil {
		return "null"
	}
	return String(v)
}

// String converts a value to string
func String(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", val)
	}
}