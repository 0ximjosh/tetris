package main

// A simple example demonstrating how to draw and animate on a cellular grid.
// Note that the cellbuffer implementation in this example does not support
// double-width runes.

import (
	"flag"
	"fmt"
	"os"

	"tetris"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Core tetris.Core
}

func (m model) Init() tea.Cmd {
	return m.Core.Animate()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, m.Core.HandleKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.Core.UpdateDims(msg.Width, msg.Height)
		return m, nil
	case tetris.FrameMsg:
		m.Core.Update()
		return m, m.Core.Animate()
	default:
		return m, nil
	}
}

func (m model) View() string {
	return m.Core.String()
}

func main() {
	fps := flag.Int("fps", 30, "fps of the simulation")
	flag.Parse()
	m := model{}
	m.Core.Fps = *fps
	m.Core.Init(0, 0)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
}
