package tetris

import (
	"math"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mattn/go-runewidth"
)

type FrameMsg struct{}

var colors = []string{"#FFFFFF", "#ff595e", "#ff924c", "#ffca3a", "#c5ca30", "#8ac926", "#52a675", "#1982c4", "#4267ac", "#6a4c93"}

func (c *Core) Animate() tea.Cmd {
	if c.paused {
		return nil
	}
	return tea.Tick(time.Second/time.Duration(c.Fps), func(_ time.Time) tea.Msg {
		return FrameMsg{}
	})
}

func removeColorFromString(text string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(text, "")
}

func (c *Core) GetGameStateString() string {
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

	if c.currentBlock != nil {
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

func (c Core) String() string {
	var b strings.Builder

	// Game View
	gameView := c.GetGameStateString()
	gameViewLines := strings.Split(gameView, "\n")
	gameViewLen := runewidth.StringWidth(removeColorFromString(gameViewLines[0]))
	gameViewStartX := int(math.Floor(float64(c.width)/2)) - (gameViewLen / 2)
	gameViewStartY := int(math.Floor(float64(c.height)/2) - 10)

	scoreBox := c.GetScoreBox()
	scoreBoxLines := strings.Split(scoreBox, "\n")
	scoreBoxLen := runewidth.StringWidth(removeColorFromString(scoreBoxLines[0]))
	scoreBoxStartX := int(math.Floor(float64(c.width)/2)) + gameViewLen/2 + 3
	scoreBoxStartY := int(math.Floor(float64(c.height)/2) - 10)

	// If screen is too small, request larger screen
	if c.width < gameViewLen || c.height < 24 {
		return "Game zone is too small\nPlease zoom out"
	}

	for y := range c.height {
		for x := range c.width {

			// gameView area
			if x >= gameViewStartX && x < gameViewStartX+gameViewLen && y >= gameViewStartY && y < gameViewStartY+len(gameViewLines) {
				if x == gameViewStartX {
					b.WriteString(gameViewLines[y-gameViewStartY])
				}
				continue
			}

			if x >= scoreBoxStartX && x < scoreBoxStartX+scoreBoxLen && y >= scoreBoxStartY && y < scoreBoxStartY+len(scoreBoxLines) {
				if x == scoreBoxStartX {
					b.WriteString(scoreBoxLines[y-scoreBoxStartY])
				}
				continue
			}
			b.WriteString(" ")
		}
		if y == c.height-1 {
			continue
		}
		b.WriteRune('\n')
	}
	return b.String()
}
