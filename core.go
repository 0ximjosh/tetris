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
	grid     [][]uint8
	width    int
	height   int
	paused   bool
	Fps      int
	hideHelp bool
	block    *Block
}

// Block is the data for a single tetris block
// x and y are the top left coord
// Each block is in a grid
type Block struct {
	x     uint8
	y     uint8
	shape [][]bool
}

func (b *Block) Rotate() {
	tmp := make([][]bool, len(b.shape[0]))
	for i := range len(b.shape[0]) {
		tmp[i] = make([]bool, len(b.shape))
		for j := range len(b.shape) {
			tmp[i][j] = b.shape[j][i]
		}
	}
	b.shape = tmp
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

// Init the game state
// Setup tetris bucket
func (c *Core) Init(wi, hi int) {
	if wi == 0 {
		return
	}
	c.width = wi
	c.height = hi
	c.grid = make([][]uint8, wi)
	for i := range c.grid {
		c.grid[i] = make([]uint8, hi)
	}

	c.block = &Block{
		x:     2,
		y:     2,
		shape: [][]bool{{false, true, false}, {true, true, true}},
	}

	// 10 blocks wide
	// 20 blocks tall
	bufferHor := int(math.Floor(float64(wi-11) / 2))
	bufferVer := int(math.Floor(float64(hi-20) / 2))
	for i := range 12 {
		for j := range 20 {
			if i != 0 && i != 11 && j != 19 {
				continue
			}
			c.grid[i+bufferHor][j+bufferVer] = 1
		}
	}
}

// Update the game grid
// Core game logic. Update core.grid accordingly for this frame
func (c *Core) Update() {
	if !c.Ready() {
		return
	}

	c.grid = make([][]uint8, c.width)
	for i := range c.grid {
		c.grid[i] = make([]uint8, c.height)
	}

	// 10 blocks wide
	// 20 blocks tall
	bufferHor := uint8(math.Floor(float64(c.width-11) / 2))
	bufferVer := uint8(math.Floor(float64(c.height-20) / 2))
	for i := range uint8(12) {
		for j := range uint8(20) {
			if c.block != nil {
				if i == c.block.x && j == c.block.y {
					for a := range uint8(len(c.block.shape)) {
						for b := range uint8(len(c.block.shape[0])) {
							if !c.block.shape[a][b] {
								continue
							}
							c.grid[i+a+bufferHor][j+b+bufferVer] = 2
						}
					}
				}
			}
			if i != 0 && i != 11 && j != 19 {
				continue
			}
			c.grid[i+bufferHor][j+bufferVer] = 1
		}
	}
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
