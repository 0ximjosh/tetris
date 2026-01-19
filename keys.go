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
	case "h", "a", "left":
		m.MoveBlock("left")
		return nil
	case "j", "s", "down":
		m.Drop()
		return nil
	case "l", "d", "right":
		m.MoveBlock("right")
		return nil
	case "q", "ctrl+c":
		return tea.Quit
	case "r", "w", "up":
		m.Rotate()
		return nil
	default:
		return nil
	}
}
