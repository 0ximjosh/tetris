package tetris

import (
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mattn/go-runewidth"
)

type FrameMsg struct{}

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

func (c Core) String() string {
	var help string
	// TODO
	//if !c.hideHelp {
	//	help = HelpBox()
	//}
	hlines := strings.Split(help, "\n")
	var b strings.Builder
	startBoxX := c.width - runewidth.StringWidth(removeColorFromString(hlines[0]))
	endBoxY := len(hlines)

	colors := []string{"#FFFFFF", "#ff595e", "#ff924c", "#ffca3a", "#c5ca30", "#8ac926", "#52a675", "#1982c4", "#4267ac", "#6a4c93"}

	for y := range c.height {
		for x := range c.width {
			if y < endBoxY && x >= startBoxX {
				b.WriteString(hlines[y])
				continue
			}
			if c.grid[x][y] == 0 {
				b.WriteRune(' ')
				continue
			}
			color, _ := colorful.Hex(colors[c.grid[x][y]-1])
			s := lipgloss.NewStyle().SetString(" ").Background(lipgloss.Color(color.Hex()))
			b.WriteString(s.String())
		}
		if y == c.height-1 {
			continue
		}
		b.WriteRune('\n')
	}
	return b.String()
}
