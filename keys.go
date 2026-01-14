// Package tetris
package tetris

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) HandleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case " ":
		if m.paused {
			m.paused = false
			return m.Animate()
		}
		m.paused = true
		return nil
	case "h":
		m.MoveBlock("h")
		return nil
	case "j":
		m.Drop()
		return nil
	case "l":
		m.MoveBlock("l")
		return nil
	case "q", "ctrl+c":
		return tea.Quit
	case "r":
		m.Rotate()
		return nil
	default:
		return nil
	}
}
