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
		c.currentBlock.Move("h")
		return nil
	case "j":
		c.currentBlock.Move("j")
		return nil
	case "l":
		c.currentBlock.Move("l")
		return nil
	case "q":
		return tea.Quit
	case "r":
		c.currentBlock.Rotate()
		return nil
	default:
		return nil
	}
}
