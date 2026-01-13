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
	// TODO
	//if !c.hideHelp {
	//	help = HelpBox()
	//}
	c.score = 1000000000

	scoreBox := c.GetGameStateString()
	scoreLines := strings.Split(scoreBox, "\n")
	var b strings.Builder
	boxLen := runewidth.StringWidth(removeColorFromString(scoreLines[0]))
	// startBoxX := c.width - boxLen
	// endBoxY := len(scoreLines)

	bufferHor := int(math.Floor(float64(c.width)/2)) - (boxLen / 2)
	bufferVer := int(math.Floor(float64(c.height)/2) - 10)

	// If screen is too small, request larger screen
	if c.width < boxLen || c.height < 24 {
		return "Game zone is too small\nPlease zoom out"
	}

	for y := range c.height {
		for x := range c.width {
			if x < bufferHor {
				b.WriteString(" ")
				continue
			}
			if x == bufferHor && y > bufferVer && y < bufferVer+len(scoreLines) {
				b.WriteString(scoreLines[y-bufferVer])
				continue
			}
			//if y < endBoxY && x >= startBoxX && x < startBoxX+1 {
			//	continue
			//}
		}
		if y == c.height-1 {
			continue
		}
		b.WriteRune('\n')
	}
	return b.String()
}
