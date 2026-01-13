package tetris

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/lucasb-eyer/go-colorful"
)

func (c *Core) GetScoreView() string {
	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#71797E")).
		Padding(1, 3, 1, 2).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Width(12)
	return dialogBoxStyle.Render(fmt.Sprintf("Score\n%d", c.score))
}

func (c *Core) GetGameOverView() string {
	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#71797E")).
		Padding(1).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Width(12)
	return dialogBoxStyle.Render("Game Over!")
}

func (c *Core) GetNextBlockView() string {
	if c.nextBlock == nil || c.currentBlock == nil {
		return ""
	}
	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#71797E")).
		Padding(1).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Width(12).
		Height(8)
	var b strings.Builder
	b.WriteString("Next Block\n\n")
	for y := range len(c.nextBlock.shape[0]) {
		b.WriteString("  ")
		for x := range len(c.nextBlock.shape) {
			if !c.nextBlock.shape[x][y] {
				b.WriteString("  ")
				continue
			}
			color, _ := colorful.Hex(colors[c.nextBlock.color-1])
			s := lipgloss.NewStyle().SetString("  ").Background(lipgloss.Color(color.Hex()))
			b.WriteString(s.String())
		}
		if y == len(c.nextBlock.shape[0])-1 {
			break
		}
		b.WriteRune('\n')
	}
	return dialogBoxStyle.Render(b.String())
}

// GetHelpBox TODO
func (c *Core) GetHelpBox() string {
	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#71797E")).
		Padding(1, 3, 1, 2).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)
	options := list.New(
		"h Show help",
		"q Quit",
		"b Blank canvas",
		"r Random canvas",
		"j Prev algo",
		"k Next algo",
		"c Cycle colors",
		"[space] Pause",
		"[click] toggle cell",
	).Enumerator(list.Dash)
	return dialogBoxStyle.Render(options.String())
}

func (c *Core) GetGameView() string {
	var b strings.Builder

	// If screen is too small, request larger screen
	if c.width == 0 || c.height == 0 {
		return ""
	}

	// First, clear screen
	for i := range len(c.grid) {
		for j := range len(c.grid[0]) {
			c.grid[i][j] = 0
		}
	}

	// 10 blocks wide
	// 20 blocks tall
	for i := range 12 {
		for j := range 21 {
			c.grid[i][j] = c.blocks[i][j]
		}
	}

	if c.currentBlock != nil && !c.gameOver {
		for a := range len(c.currentBlock.shape) {
			for b := range len(c.currentBlock.shape[0]) {
				if !c.currentBlock.shape[a][b] {
					continue
				}
				c.grid[c.currentBlock.x+a][c.currentBlock.y+b] = c.currentBlock.color
			}
		}
	}

	for y := range len(c.grid[0]) {
		for x := range len(c.grid) {
			if c.grid[x][y] == 0 {
				b.WriteString("  ")
				continue
			}
			color, _ := colorful.Hex(colors[c.grid[x][y]-1])
			s := lipgloss.NewStyle().SetString("  ").Background(lipgloss.Color(color.Hex()))
			b.WriteString(s.String())
		}
		b.WriteRune('\n')
	}

	return b.String()
}
