package main

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
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
}

func NewModel() model {
	ti := textinput.New()
	ti.Placeholder = "Your API key"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)
	return model{
		state:     pathView,
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
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
					m.dir = getDirectory("")
					m.state = menuView
					m.cursor = 0
				}
			case menuView:
				if m.dir.isDirectory(m.cursor) {
					m.dir = getDirectory(m.dir.Path + m.dir.Files[m.cursor])
					m.cursor = 0
				} else {
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
			}
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		}

	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}
