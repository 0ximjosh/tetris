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

type Core struct {
	grid         [][]uint8
	width        int
	height       int
	paused       bool
	Fps          int
	hideHelp     bool
	currentBlock *Block
	blocks       [][]uint8
	tick         int
}

type FrameMsg struct{}

func (c *Core) Animate() tea.Cmd {
	if c.paused {
		return nil
	}
	return tea.Tick(time.Second/time.Duration(c.Fps), func(_ time.Time) tea.Msg {
		return FrameMsg{}
	})
}

func (c *Core) WriteState() {
	// First, clear screen
	for i := range len(c.grid) {
		for j := range len(c.grid[0]) {
			c.grid[i][j] = 0
		}
	}

	// If screen is too small, request larger screen
	if c.width < 14 || c.height < 24 {
		return
	}

	// 10 blocks wide
	// 20 blocks tall
	bufferHor := uint8(math.Floor(float64(c.width)/2) - 6)
	bufferVer := uint8(math.Floor(float64(c.height)/2) - 10)
	for i := range uint8(12) {
		for j := range uint8(21) {

			if i == 0 || i == 11 || j == 20 {
				c.grid[i+bufferHor][j+bufferVer] = 1
				continue
			}
			c.grid[i+bufferHor][j+bufferVer] = c.blocks[i-1][j]

			// Write current block
		}
	}

	for i := range uint8(12) {
		for j := range uint8(21) {
			if c.currentBlock != nil {
				if i == c.currentBlock.x && j == c.currentBlock.y {
					for a := range uint8(len(c.currentBlock.shape)) {
						for b := range uint8(len(c.currentBlock.shape[0])) {
							if !c.currentBlock.shape[a][b] {
								continue
							}
							c.grid[i+a+bufferHor+1][j+b+bufferVer] = c.currentBlock.color
						}
					}
				}
			}
		}
	}
}

// Init the game state
// Setup tetris bucket
func (c *Core) Init(w, h int) {
	if w < 0 {
		return
	}
	c.width = w
	c.height = h
	c.grid = make([][]uint8, w)
	for i := range c.grid {
		c.grid[i] = make([]uint8, h)
	}
	c.blocks = make([][]uint8, 10)
	for i := range c.blocks {
		c.blocks[i] = make([]uint8, 20)
	}

	c.currentBlock = &Block{
		x:     2,
		y:     2,
		color: 2,
		shape: [][]bool{{false, true, false}, {true, true, true}},
	}
	c.WriteState()
}

func (c *Core) PlaceCurrentBlock() {
	for i := range uint8(len(c.currentBlock.shape)) {
		for j := range uint8(len(c.currentBlock.shape[0])) {
			if !c.currentBlock.shape[i][j] {
				continue
			}
			c.blocks[i+c.currentBlock.x][j+c.currentBlock.y] = c.currentBlock.color
		}
	}

	c.currentBlock = &Block{
		x:     2,
		y:     2,
		color: 2,
		shape: [][]bool{{false, true, false}, {true, true, true}},
	}
}

// Update the game grid
// Core game logic. Update core.grid accordingly for this frame
func (c *Core) Update() {
	if !c.Ready() {
		return
	}

	c.tick++
	if c.tick%5 == 0 {
		if c.currentBlock.y+uint8(len(c.currentBlock.shape[0])) >= 20 {
			c.PlaceCurrentBlock()
		}
		if c.currentBlock.y+uint8(len(c.currentBlock.shape[0])) < 20 {
			c.currentBlock.y++
		}
	}
	c.WriteState()
}

func (c Core) Ready() bool {
	return len(c.grid) > 0
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

	colors := []string{"#FFFFFF", "#FF8156", "#7CC0FF", "#58E8B4", "#8FA0FE"}

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
