// Package tetris
package tetris

import tea "github.com/charmbracelet/bubbletea"

func (c *Core) HandleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case " ":
		if c.paused {
			c.paused = false
			return c.Animate()
		}
		c.paused = true
		return nil
	case "h":
		c.MoveBlock("h")
		return nil
	case "j":
		c.Drop()
		return nil
	case "l":
		c.MoveBlock("l")
		return nil
	case "q":
		return tea.Quit
	case "r":
		c.Rotate()
		return nil
	default:
		return nil
	}
}
