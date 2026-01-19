package tetris

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/lucasb-eyer/go-colorful"
)

func (m *Model) GetScoreView() string {
	dialogBoxStyle := m.Renderer.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#71797E")).
		Padding(1, 3, 1, 2).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Width(12)
	return dialogBoxStyle.Render(fmt.Sprintf("Score\n%d", m.score))
}

func (m *Model) GetGameOverView() string {
	dialogBoxStyle := m.Renderer.NewStyle().
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

func (m *Model) GetNextBlockView() string {
	if m.nextBlock == nil || m.currentBlock == nil {
		return ""
	}
	dialogBoxStyle := m.Renderer.NewStyle().
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
	for y := range len(m.nextBlock.shape[0]) {
		b.WriteString("  ")
		for x := range len(m.nextBlock.shape) {
			if !m.nextBlock.shape[x][y] {
				b.WriteString("  ")
				continue
			}
			color, _ := colorful.Hex(colors[m.nextBlock.color-1])
			s := m.Renderer.NewStyle().SetString("  ").Background(lipgloss.Color(color.Hex()))
			b.WriteString(s.String())
		}
		if y == len(m.nextBlock.shape[0])-1 {
			break
		}
		b.WriteRune('\n')
	}
	return dialogBoxStyle.Render(b.String())
}

func (m *Model) GetHelpBox() string {
	dialogBoxStyle := m.Renderer.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#71797E")).
		Padding(1, 3, 1, 2).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)
	options := list.New(
		"h Move Left",
		"j Move Down",
		"l Move Right",
		"r Rotate",
		"q Quit",
	).Enumerator(list.Dash)
	return dialogBoxStyle.Render(options.String())
}

func (m *Model) GetGameView() string {
	var b strings.Builder

	// If screen is too small, request larger screen
	if m.width == 0 || m.height == 0 {
		return ""
	}

	// First, clear screen
	for i := range len(m.grid) {
		for j := range len(m.grid[0]) {
			m.grid[i][j] = 0
		}
	}

	// 10 blocks wide
	// 20 blocks tall
	for i := range 12 {
		for j := range 21 {
			m.grid[i][j] = m.blocks[i][j]
		}
	}

	if m.currentBlock != nil && !m.gameOver {
		for a := range len(m.currentBlock.shape) {
			for b := range len(m.currentBlock.shape[0]) {
				if !m.currentBlock.shape[a][b] {
					continue
				}
				m.grid[m.currentBlock.x+a][m.currentBlock.y+b] = m.currentBlock.color
			}
		}
	}

	for y := range len(m.grid[0]) {
		for x := range len(m.grid) {
			if m.grid[x][y] == 0 {
				b.WriteString("  ")
				continue
			}
			color, _ := colorful.Hex(colors[m.grid[x][y]-1])
			s := m.Renderer.NewStyle().SetString("  ").Background(lipgloss.Color(color.Hex()))
			b.WriteString(s.String())
		}
		b.WriteRune('\n')
	}

	return b.String()
}
