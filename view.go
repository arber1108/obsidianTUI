package main

import (
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"

	tea "charm.land/bubbletea/v2"
)

var BorderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63"))

func renderTree(m model) string {
	root := "vault"
	if m.dir.Path != "" {
		trimmed := strings.TrimSuffix(m.dir.Path, "/")
		if idx := strings.LastIndex(trimmed, "/"); idx >= 0 {
			root = trimmed[idx+1:]
		} else {
			root = trimmed
		}
	}

	children := make([]any, len(m.dir.Files))
	for i, name := range m.dir.Files {
		children[i] = name
	}

	t := tree.Root(root+"/").
		Child(children...).
		EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginRight(1)).
		RootStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("33")).Bold(true)).
		ItemStyleFunc(func(_ tree.Children, i int) lipgloss.Style {
			if i == m.cursor {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
			}
			if m.dir.isDirectory(i) {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
			}
			return lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
		})

	return t.String()
}

func (m model) View() tea.View {
	var body strings.Builder
	switch m.state {
	case pathView:
		body.WriteString(m.textInput.View())
		if m.info != "" {
			body.WriteString("\n" + m.info)
		}
	case menuView:
		body.WriteString(BorderStyle.Render(renderTree(m)))
	case editView:
		body.WriteString(BorderStyle.Render(m.viewport.View()))
	}

	return tea.NewView(body.String())
}
