package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingLeft(2).
			PaddingRight(2)

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Italic(true).
		PaddingLeft(1)
)

// Model represents the application state
type Model struct {
	title  string
	width  int
	height int
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and user input
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renders the UI
func (m Model) View() string {
	title := titleStyle.Render(" TABLUX ")
	info := infoStyle.Render("Press 'q' to quit")

	return fmt.Sprintf("%s\n\n%s\n\nWindow size: %d x %d\n", 
		title, info, m.width, m.height)
}

func main() {
	m := Model{
		title: "Tablux",
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}