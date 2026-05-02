package main

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/glamour"

	tea "charm.land/bubbletea/v2"
)

var CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

func (m model) View() tea.View {
	var body strings.Builder

	switch m.state {
	case pathView:
		body.WriteString(m.textInput.View())
		if m.info != "" {
			body.WriteString("\n" + m.info)
		}
	case menuView:
		if m.dir.Path != "" {
			body.WriteString(fmt.Sprintf("  %s\n\n", m.dir.Path))
		}
		for i, name := range m.dir.Files {
			cursor := "  "
			if i == m.cursor {
				cursor = CursorStyle.Render("> ")
			}
			body.WriteString(fmt.Sprintf("%s%s\n", cursor, name))
		}
	case editView:
		rendered, err := glamour.Render(getFile(m.dir.Path+m.dir.Files[m.cursor]), "dark")
		if err != nil {
			body.WriteString("Couldn't open the File")
		}

		body.WriteString(rendered)
	}

	return tea.NewView(body.String())
}
