package tetris

import (
	"math"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func (c Core) String() string {
	var b strings.Builder

	// Game View
	gameView := c.GetGameView()
	gameViewLines := strings.Split(gameView, "\n")
	gameViewLen := runewidth.StringWidth(removeColorFromString(gameViewLines[0]))
	gameViewStartX := int(math.Floor(float64(c.width)/2)) - (gameViewLen / 2)
	gameViewStartY := int(math.Floor(float64(c.height)/2) - 10)

	playtime := lipgloss.NewStyle().Width(gameViewLen).Align(lipgloss.Center).Render(time.Since(c.startTime).Truncate(time.Second).String())
	gameViewLines = append(gameViewLines, playtime)

	// TODO fix score box and next box being out of view
	scoreBox := c.GetScoreView()
	scoreBoxLines := strings.Split(scoreBox, "\n")
	scoreBoxLen := runewidth.StringWidth(removeColorFromString(scoreBoxLines[0]))
	scoreBoxStartX := int(math.Floor(float64(c.width)/2)) - gameViewLen/2 - scoreBoxLen - 3
	scoreBoxStartY := int(math.Floor(float64(c.height)/2) - 10)

	nextBlockView := c.GetNextBlockView()
	nextBlockViewLines := strings.Split(nextBlockView, "\n")
	nextBlockViewLen := runewidth.StringWidth(removeColorFromString(nextBlockViewLines[0]))
	nextBlockViewStartX := int(math.Floor(float64(c.width)/2)) - gameViewLen/2 - nextBlockViewLen - 3
	nextBlockViewStartY := int(math.Floor(float64(c.height)/2)) + len(scoreBoxLines) - 10

	gameOverView := c.GetGameOverView()
	gameOverViewLines := strings.Split(gameOverView, "\n")
	gameOverViewLen := runewidth.StringWidth(removeColorFromString(gameOverViewLines[0]))
	gameOverViewStartX := int(math.Floor(float64(c.width)/2)) - gameViewLen/2 - nextBlockViewLen - 3
	gameOverViewStartY := int(math.Floor(float64(c.height)/2)) + len(scoreBoxLines) + len(nextBlockViewLines) - 10

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

			if x >= nextBlockViewStartX && x < nextBlockViewStartX+nextBlockViewLen && y >= nextBlockViewStartY && y < nextBlockViewStartY+len(nextBlockViewLines) {
				if x == nextBlockViewStartX {
					b.WriteString(nextBlockViewLines[y-nextBlockViewStartY])
				}
				continue
			}

			if x >= scoreBoxStartX && x < scoreBoxStartX+scoreBoxLen && y >= scoreBoxStartY && y < scoreBoxStartY+len(scoreBoxLines) {
				if x == scoreBoxStartX {
					b.WriteString(scoreBoxLines[y-scoreBoxStartY])
				}
				continue
			}

			if x >= gameOverViewStartX && x < gameOverViewStartX+gameOverViewLen && y >= gameOverViewStartY && y < gameOverViewStartY+len(gameOverViewLines) && c.gameOver {
				if x == gameOverViewStartX {
					b.WriteString(gameOverViewLines[y-gameOverViewStartY])
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
