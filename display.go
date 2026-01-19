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

// Animate sends a frame msg
// TODO rework so we only call this when we need, not at a crazy fps
func (m Model) Animate() tea.Cmd {
	if m.paused {
		return nil
	}
	return tea.Tick(time.Second/time.Duration(m.tickSpeed), func(_ time.Time) tea.Msg {
		return FrameMsg{}
	})
}

func removeColorFromString(text string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(text, "")
}

func (m Model) String() string {
	var b strings.Builder

	// Game View
	gameView := m.GetGameView()
	gameViewLines := strings.Split(gameView, "\n")
	gameViewLen := runewidth.StringWidth(removeColorFromString(gameViewLines[0]))
	gameViewStartX := int(math.Floor(float64(m.width)/2)) - (gameViewLen / 2)
	gameViewStartY := int(math.Floor(float64(m.height)/2) - 10)
	playtime := lipgloss.NewStyle().Width(gameViewLen).Align(lipgloss.Center).Render(time.Since(m.startTime).Truncate(time.Second).String())
	if m.gameOver {
		playtime = lipgloss.NewStyle().Width(gameViewLen).Align(lipgloss.Center).Render(m.endTime.Sub(m.startTime).Truncate(time.Second).String())
	}
	gameViewLines = append(gameViewLines, playtime)
	gameViewBox := ViewBox{Lines: gameViewLines, Len: gameViewLen, X: gameViewStartX, Y: gameViewStartY, Visable: true}

	scoreBox := m.GetScoreView()
	scoreBoxLines := strings.Split(scoreBox, "\n")
	scoreBoxLen := runewidth.StringWidth(removeColorFromString(scoreBoxLines[0]))
	scoreBoxStartX := int(math.Floor(float64(m.width)/2)) - gameViewLen/2 - scoreBoxLen - 3
	scoreBoxStartY := int(math.Floor(float64(m.height)/2) - 10)
	scoreViewBox := ViewBox{Lines: scoreBoxLines, Len: scoreBoxLen, X: scoreBoxStartX, Y: scoreBoxStartY, Visable: !m.gameOver}

	nextBlockView := m.GetNextBlockView()
	nextBlockViewLines := strings.Split(nextBlockView, "\n")
	nextBlockViewLen := runewidth.StringWidth(removeColorFromString(nextBlockViewLines[0]))
	nextBlockViewStartX := int(math.Floor(float64(m.width)/2)) - gameViewLen/2 - nextBlockViewLen - 3
	nextBlockViewStartY := int(math.Floor(float64(m.height)/2)) + len(scoreBoxLines) - 10
	nextBlockViewBox := ViewBox{Lines: nextBlockViewLines, Len: nextBlockViewLen, X: nextBlockViewStartX, Y: nextBlockViewStartY, Visable: !m.gameOver}

	gameOverView := m.GetGameOverView()
	gameOverViewLines := strings.Split(gameOverView, "\n")
	gameOverViewLen := runewidth.StringWidth(removeColorFromString(gameOverViewLines[0]))
	gameOverViewStartX := int(math.Floor(float64(m.width)/2)) - gameViewLen/2 - nextBlockViewLen - 3
	gameOverViewStartY := int(math.Floor(float64(m.height)/2)) + len(scoreBoxLines) + len(nextBlockViewLines) - 10
	gameOverViewBox := ViewBox{Lines: gameOverViewLines, Len: gameOverViewLen, X: gameOverViewStartX, Y: gameOverViewStartY, Visable: m.gameOver}

	helpView := m.GetHelpBox()
	helpViewLines := strings.Split(helpView, "\n")
	helpViewLen := runewidth.StringWidth(removeColorFromString(helpViewLines[0]))
	helpViewStartX := int(math.Floor(float64(m.width)/2)) + gameViewLen/2 + 3
	helpViewStartY := int(math.Floor(float64(m.height)/2)) - 10
	helpViewBox := ViewBox{Lines: helpViewLines, Len: helpViewLen, X: helpViewStartX, Y: helpViewStartY, Visable: true}

	// If screen is too small, request larger screen
	if m.width < gameViewLen || m.height < 24 {
		return "Game zone is too small\nPlease zoom out"
	}
	shouldPrintBoxes := m.width > gameViewLen+nextBlockViewLen+24
	nextBlockViewBox.Visable = shouldPrintBoxes && nextBlockViewBox.Visable
	helpViewBox.Visable = shouldPrintBoxes
	scoreViewBox.Visable = shouldPrintBoxes
	gameOverViewBox.Visable = shouldPrintBoxes && gameOverViewBox.Visable

	boxes := []ViewBox{gameViewBox, nextBlockViewBox, scoreViewBox, gameOverViewBox, helpViewBox}

	for y := range m.height {
	layer: // label for continueing outer loop
		for x := range m.width {
			for _, box := range boxes {
				if box.MaybeWriteLine(x, y, &b) {
					continue layer
				}
			}
			b.WriteString(" ")
		}
		if y == m.height-1 {
			continue
		}
		b.WriteRune('\n')
	}
	return b.String()
}

type ViewBox struct {
	Lines   []string
	Len     int
	X       int
	Y       int
	Visable bool
}

func (v *ViewBox) MaybeWriteLine(x, y int, b *strings.Builder) bool {
	if !v.Visable {
		return false
	}
	if x >= v.X && x < v.X+v.Len && y >= v.Y && y < v.Y+len(v.Lines) {
		if x == v.X {
			b.WriteString(v.Lines[y-v.Y])
		}
		return true
	}
	return false
}
