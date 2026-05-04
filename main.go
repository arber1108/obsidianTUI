package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

func main() {

	m := NewModel()

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Unable to run TUI: %v", err)
	}
}
