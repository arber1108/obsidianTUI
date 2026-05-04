package main

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/glamour"
)

const (
	pathView uint = iota
	menuView
	editView
)

type model struct {
	vaultPath string
	state     uint
	textInput textinput.Model
	info      string
	dir       Directory
	cursor    int
	viewport  viewport.Model
}

func NewModel() model {
	ti := textinput.New()
	ti.Placeholder = "Your API key"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)
	m := model{
		state:     pathView,
		textInput: ti,
		viewport:  viewport.New(viewport.WithWidth(80), viewport.WithHeight(40)),
	}

	if key := loadSavedApiKey(); key != "" {
		if code, _ := checkApiKey(key); code == 200 {
			setApiKey(key)
			m.dir = getDirectory("")
			m.state = menuView
		}
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width - 4)
		m.viewport.SetHeight(msg.Height - 4)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			switch m.state {
			case pathView:
				apiKey := m.textInput.Value()
				code, err := checkApiKey(apiKey)
				if code != 200 {
					m.info = fmt.Sprintf("API key is invalid: %v", err)
					m.textInput.SetValue("")
				} else {
					m.info = "API key is Valid"
					setApiKey(apiKey)
					saveApiKey(apiKey)
					m.dir = getDirectory("")
					m.state = menuView
					m.cursor = 0
				}
			case menuView:
				if m.dir.isDirectory(m.cursor) {
					m.dir = getDirectory(m.dir.Path + m.dir.Files[m.cursor])
					m.cursor = 0
				} else {
					file := getFile(m.dir.Path + m.dir.Files[m.cursor])
					renderer, err := glamour.NewTermRenderer(
						glamour.WithStylePath("dark"),
						glamour.WithWordWrap(m.viewport.Width()),
					)
					rendered := file.Content
					if err == nil {
						if out, err := renderer.Render(file.Content); err == nil {
							rendered = out
						}
					}
					m.viewport.SetContent(rendered)
					m.viewport.GotoTop()
					m.state = editView
				}
			}

		case "q":
			return m, tea.Quit

		case "backspace":
			if m.state == menuView && m.dir.Path != "" {
				m.dir = getDirectory(m.dir.parentPath())
				m.cursor = 0
			}
			if m.state == editView {
				m.state = menuView
			}

		case "down":
			if m.state == menuView {
				if m.cursor < len(m.dir.Files)-1 {
					m.cursor++
				}
			} else if m.state == editView {
				m.viewport.ScrollDown(1)
			}

		case "up":
			if m.state == editView {
				m.viewport.ScrollUp(1)
			} else if m.cursor > 0 {
				m.cursor--
			}

		case "o":
			if m.state == editView {
				m.dir.openInObsidian(m.cursor)
			}
		}
	}

	m.viewport, _ = m.viewport.Update(msg)
	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}
