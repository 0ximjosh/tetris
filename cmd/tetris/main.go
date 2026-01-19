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

func main() {
	flag.Parse()
	m := tetris.Model{}
	m.Reset()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
}
