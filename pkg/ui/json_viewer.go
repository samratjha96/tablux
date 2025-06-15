package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"tablux/pkg/model"
)

var (
	// JSON node colors
	keyStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#88AAFF"))
	stringStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#7CFC00"))
	numberStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	boolStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9F5F"))
	nullStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F5F"))
	bracketStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F8F8F2"))
	selectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("#404040"))
	separatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	
	// Tree symbols
	treeStyles = map[string]string{
		"pipe":     "│ ",
		"tee":      "├─",
		"last":     "└─",
		"expanded": "▼ ",
		"collapsed": "► ",
		"empty":    "  ",
	}
	
	// Spacing between columns for readability
	valuePadding = 2
)

// JSONViewer displays a JSON tree
type JSONViewer struct {
	root         *model.JSONNode
	cursor       int
	nodes        []*model.JSONNode // Flattened list for navigation
	visibleNodes []*model.JSONNode // Current visible nodes
	viewportY    int
	viewportHeight int
	maxKeyWidth   int              // For alignment
}

// NewJSONViewer creates a new JSON viewer
func NewJSONViewer(root *model.JSONNode) *JSONViewer {
	viewer := &JSONViewer{
		root:         root,
		cursor:       0,
		viewportHeight: 20, // Default height
	}
	viewer.buildNodeList()
	return viewer
}

// buildNodeList creates a flattened list of visible nodes
func (v *JSONViewer) buildNodeList() {
	v.nodes = make([]*model.JSONNode, 0)
	v.visibleNodes = make([]*model.JSONNode, 0)
	v.flattenNode(v.root, 0)
	v.updateVisibleNodes()
}

// flattenNode adds a node and its visible children to the nodes list
func (v *JSONViewer) flattenNode(node *model.JSONNode, depth int) {
	v.nodes = append(v.nodes, node)
	
	if !node.Expanded {
		return
	}
	
	for _, child := range node.Children {
		v.flattenNode(child, depth+1)
	}
}

// updateVisibleNodes updates the list of visible nodes
func (v *JSONViewer) updateVisibleNodes() {
	v.visibleNodes = make([]*model.JSONNode, 0)
	
	for _, node := range v.nodes {
		// Check if node should be visible
		parent := node.Parent
		isVisible := true
		
		for parent != nil {
			if !parent.Expanded {
				isVisible = false
				break
			}
			parent = parent.Parent
		}
		
		if isVisible {
			v.visibleNodes = append(v.visibleNodes, node)
		}
	}
	
	// Make sure cursor is still valid
	if v.cursor >= len(v.visibleNodes) && len(v.visibleNodes) > 0 {
		v.cursor = len(v.visibleNodes) - 1
	}
}

// MoveUp moves the cursor up
func (v *JSONViewer) MoveUp() {
	if v.cursor > 0 {
		v.cursor--
	}
	v.ensureCursorVisible()
}

// MoveDown moves the cursor down
func (v *JSONViewer) MoveDown() {
	if v.cursor < len(v.visibleNodes)-1 {
		v.cursor++
	}
	v.ensureCursorVisible()
}

// ToggleNode expands or collapses the current node
func (v *JSONViewer) ToggleNode() {
	if v.cursor < len(v.visibleNodes) {
		node := v.visibleNodes[v.cursor]
		if node.HasChildren() {
			node.Toggle()
			v.buildNodeList()
		}
	}
}

// ensureCursorVisible adjusts viewport to keep cursor in view
func (v *JSONViewer) ensureCursorVisible() {
	if v.cursor < v.viewportY {
		v.viewportY = v.cursor
	} else if v.cursor >= v.viewportY+v.viewportHeight {
		v.viewportY = v.cursor - v.viewportHeight + 1
	}
}

// CollapseAll collapses all nodes in the tree
func (v *JSONViewer) CollapseAll() {
	v.toggleAllNodes(v.root, false)
	v.buildNodeList()
}

// ExpandAll expands all nodes in the tree
func (v *JSONViewer) ExpandAll() {
	v.toggleAllNodes(v.root, true)
	v.buildNodeList()
}

// toggleAllNodes sets the expanded state for the given node and all its children
func (v *JSONViewer) toggleAllNodes(node *model.JSONNode, expanded bool) {
	if node.HasChildren() {
		node.Expanded = expanded
		for _, child := range node.Children {
			v.toggleAllNodes(child, expanded)
		}
	}
}

// SetViewportHeight sets the height of the viewport
func (v *JSONViewer) SetViewportHeight(height int) {
	v.viewportHeight = height
	v.ensureCursorVisible()
}

// Render renders the JSON viewer
func (v *JSONViewer) Render() string {
	if len(v.visibleNodes) == 0 {
		return "Empty JSON"
	}

	var sb strings.Builder

	// Calculate visible range
	endIdx := v.viewportY + v.viewportHeight
	if endIdx > len(v.visibleNodes) {
		endIdx = len(v.visibleNodes)
	}

	// Render visible nodes
	for i := v.viewportY; i < endIdx; i++ {
		node := v.visibleNodes[i]
		line := v.renderNode(node, i == v.cursor)
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderNode renders a single node
func (v *JSONViewer) renderNode(node *model.JSONNode, selected bool) string {
	indent := v.getIndentation(node)
	nodeText := v.formatNode(node)
	
	line := indent + nodeText
	if selected {
		return selectedStyle.Render(line)
	}
	return line
}

// getIndentation returns the tree indentation for a node
func (v *JSONViewer) getIndentation(node *model.JSONNode) string {
	var result strings.Builder
	
	// Calculate ancestry for tree drawing
	var ancestry []*model.JSONNode
	current := node
	for current.Parent != nil {
		ancestry = append([]*model.JSONNode{current.Parent}, ancestry...)
		current = current.Parent
	}
	
	// Draw tree branches
	for i := 1; i < len(ancestry); i++ {
		parent := ancestry[i]
		isLast := false
		
		if i == len(ancestry)-1 {
			// Check if this is the last child of its parent
			children := parent.Children
			for j, child := range children {
				if child == node && j == len(children)-1 {
					isLast = true
				}
			}
		}
		
		if isLast {
			result.WriteString(treeStyles["empty"])
		} else {
			result.WriteString(treeStyles["pipe"])
		}
	}
	
	// Add expand/collapse symbol if needed
	if node.HasChildren() {
		if node.Expanded {
			result.WriteString(treeStyles["expanded"])
		} else {
			result.WriteString(treeStyles["collapsed"])
		}
	} else {
		result.WriteString("  ")
	}
	
	return result.String()
}

// formatNode formats a node for display
func (v *JSONViewer) formatNode(node *model.JSONNode) string {
	key := node.Key
	if key != "" && key != "root" {
		key = fmt.Sprintf("\"%s\"", key)
	} else if key == "root" {
		key = ""
	}
	
	keyFormatted := keyStyle.Render(key)
	
	// Add colon and padding for better readability
	separator := ""
	if key != "" {
		separator = separatorStyle.Render(": " + strings.Repeat(" ", valuePadding))
	}
	
	switch node.Type {
	case model.NodeObject:
		if node.Expanded {
			return keyFormatted + separator + bracketStyle.Render("{")
		} else {
			childCount := len(node.Children)
			return keyFormatted + separator + bracketStyle.Render(fmt.Sprintf("{ %d %s }", childCount, pluralize("item", childCount)))
		}
	case model.NodeArray:
		if node.Expanded {
			return keyFormatted + separator + bracketStyle.Render("[")
		} else {
			childCount := len(node.Children)
			return keyFormatted + separator + bracketStyle.Render(fmt.Sprintf("[ %d %s ]", childCount, pluralize("item", childCount)))
		}
	case model.NodeString:
		return keyFormatted + separator + stringStyle.Render(fmt.Sprintf("\"%s\"", node.Value.(string)))
	case model.NodeNumber:
		return keyFormatted + separator + numberStyle.Render(model.String(node.Value))
	case model.NodeBoolean:
		return keyFormatted + separator + boolStyle.Render(model.String(node.Value))
	case model.NodeNull:
		return keyFormatted + separator + nullStyle.Render("null")
	default:
		return keyFormatted + separator + model.String(node.Value)
	}
}

// pluralize returns singular or plural word form
func pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}