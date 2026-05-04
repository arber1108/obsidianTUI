package main

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"

	tea "charm.land/bubbletea/v2"
)

var (
	colorPrimary = lipgloss.Color("#7C3AED")
	colorAccent  = lipgloss.Color("#A78BFA")
	colorDir     = lipgloss.Color("#60A5FA")
	colorFile    = lipgloss.Color("#CBD5E1")
	colorDim     = lipgloss.Color("#475569")
	colorSuccess = lipgloss.Color("#34D399")
	colorError   = lipgloss.Color("#F87171")
	colorBorder  = lipgloss.Color("#312E81")
	colorEnum    = lipgloss.Color("#6D28D9")
)

var (
	headerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1E1B4B")).
			Foreground(colorAccent).
			Bold(true).
			PaddingLeft(2).PaddingRight(2)

	panelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	successStyle = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	labelStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)
)

func hint(key, desc string) string {
	k := lipgloss.NewStyle().Foreground(colorAccent).Bold(true).Render(key)
	d := dimStyle.Render(" " + desc)
	return k + d
}

func renderFooterHints(pairs [][2]string) string {
	sep := dimStyle.Render("  ·  ")
	parts := make([]string, len(pairs))
	for i, p := range pairs {
		parts[i] = hint(p[0], p[1])
	}
	return strings.Join(parts, sep)
}

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
		EnumeratorStyle(lipgloss.NewStyle().Foreground(colorEnum).MarginRight(1)).
		RootStyle(lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)).
		ItemStyleFunc(func(_ tree.Children, i int) lipgloss.Style {
			if i == m.cursor {
				return lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
			}
			if m.dir.isDirectory(i) {
				return lipgloss.NewStyle().Foreground(colorDir)
			}
			return lipgloss.NewStyle().Foreground(colorFile)
		})

	return t.String()
}

func (m model) View() tea.View {
	var body strings.Builder

	switch m.state {
	case pathView:
		title := lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			Render("◆ Obsidian TUI")

		subtitle := dimStyle.Render("connect via Obsidian Local REST API")
		label := labelStyle.Render("API Key")
		input := panelStyle.Render(m.textInput.View())

		var statusLine string
		if m.info != "" {
			if strings.Contains(strings.ToLower(m.info), "invalid") {
				statusLine = errorStyle.Render("✗  " + m.info)
			} else {
				statusLine = successStyle.Render("✓  " + m.info)
			}
		}

		footer := renderFooterHints([][2]string{{"enter", "confirm"}, {"q", "quit"}})

		body.WriteString("\n\n  " + title + "\n")
		body.WriteString("  " + subtitle + "\n\n")
		body.WriteString("  " + label + "\n")
		body.WriteString("  " + input + "\n")
		if statusLine != "" {
			body.WriteString("\n  " + statusLine + "\n")
		}
		body.WriteString("\n  " + footer)

	case menuView:
		pathDisplay := m.dir.Path
		if pathDisplay == "" {
			pathDisplay = "/"
		}
		header := headerStyle.Render("◆  vault  /  " + pathDisplay)
		content := panelStyle.Render(renderTree(m))
		footer := renderFooterHints([][2]string{
			{"↑↓", "navigate"}, {"enter", "open"}, {"⌫", "back"}, {"q", "quit"},
		})

		body.WriteString(header + "\n")
		body.WriteString(content + "\n")
		body.WriteString(footer)

	case editView:
		fileName := ""
		if m.cursor >= 0 && m.cursor < len(m.dir.Files) {
			fileName = m.dir.Files[m.cursor]
		}

		header := headerStyle.Render("◆  " + fileName)
		content := panelStyle.Render(m.viewport.View())

		scrollPct := int(m.viewport.ScrollPercent() * 100)
		scroll := dimStyle.Render(fmt.Sprintf("%d%%", scrollPct))
		footer := renderFooterHints([][2]string{
			{"↑↓", "scroll"}, {"o", "open in obsidian"}, {"⌫", "back"}, {"q", "quit"},
		}) + "  " + scroll

		body.WriteString(header + "\n")
		body.WriteString(content + "\n")
		body.WriteString(footer)
	}

	return tea.NewView(body.String())
}
