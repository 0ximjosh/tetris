package tetris

import (
	"math"
	"strconv"
	"time"
)

const (
	defaultFg = "#cbd5e1"
	defaultBg = "#020617"
	borderFg  = "#71797E"
)

type Cell struct {
	X     int      `json:"x"`
	Y     int      `json:"y"`
	Ch    string   `json:"ch"`
	Fg    string   `json:"fg,omitempty"`
	Bg    string   `json:"bg,omitempty"`
	Attrs []string `json:"attrs,omitempty"`
}

func (m *Model) HandleGGPInput(key string) {
	switch key {
	case " ":
		m.paused = !m.paused
	case "h", "a", "left":
		m.MoveBlock("left")
	case "j", "s", "down":
		m.Drop()
	case "l", "d", "right":
		m.MoveBlock("right")
	case "r", "w", "up":
		m.Rotate()
	}
}

func (m *Model) Advance() {
	if m.paused {
		return
	}
	m.Tick()
}

func (m Model) FrameInterval() time.Duration {
	if m.tickSpeed < 1 {
		return time.Second
	}
	return time.Second / time.Duration(m.tickSpeed)
}

func (m Model) Score() int64 {
	if m.score > uint64(math.MaxInt64) {
		return math.MaxInt64
	}
	return int64(m.score)
}

func (m *Model) Cells() []Cell {
	if !m.Ready() {
		return nil
	}

	cells := make([]Cell, 0, m.width*m.height)
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			cells = append(cells, Cell{X: x, Y: y, Ch: " ", Fg: defaultFg, Bg: defaultBg})
		}
	}

	if m.width < 32 || m.height < 24 {
		writeText(cells, m.width, m.height, 1, 1, "Game zone is too small", "#f87171", defaultBg, true)
		writeText(cells, m.width, m.height, 1, 2, "Please resize your terminal", defaultFg, defaultBg, false)
		return cells
	}

	m.copyBlocksToGrid()
	boardWidth := len(m.grid) * 2
	boardHeight := len(m.grid[0])
	boardX := maxInt((m.width-boardWidth)/2, 0)
	boardY := maxInt((m.height-boardHeight-1)/2, 0)

	drawBoard(cells, m.width, m.height, boardX, boardY, m.grid)
	drawPlaytime(cells, m.width, m.height, boardX, boardY+boardHeight, boardWidth, m.playtime())

	showSidebars := m.width >= boardWidth+44
	if showSidebars {
		sidebarX := boardX + boardWidth + 3
		drawScore(cells, m.width, m.height, sidebarX, boardY, m.score)
		drawNextBlock(cells, m.width, m.height, sidebarX, boardY+5, m.nextBlock)
		drawHelp(cells, m.width, m.height, sidebarX, boardY+14)
	}

	if m.gameOver {
		drawGameOver(cells, m.width, m.height, boardX, boardY+8, boardWidth)
	}

	return cells
}

func (m *Model) copyBlocksToGrid() {
	for x := range len(m.grid) {
		for y := range len(m.grid[0]) {
			m.grid[x][y] = 0
		}
	}

	for x := range len(m.blocks) {
		for y := range len(m.blocks[0]) {
			m.grid[x][y] = m.blocks[x][y]
		}
	}

	if m.currentBlock == nil || m.gameOver {
		return
	}
	for x := range len(m.currentBlock.shape) {
		for y := range len(m.currentBlock.shape[0]) {
			if !m.currentBlock.shape[x][y] {
				continue
			}
			gx := m.currentBlock.x + x
			gy := m.currentBlock.y + y
			if gx >= 0 && gx < len(m.grid) && gy >= 0 && gy < len(m.grid[0]) {
				m.grid[gx][gy] = m.currentBlock.color
			}
		}
	}
}

func (m Model) playtime() string {
	if m.gameOver {
		return m.endTime.Sub(m.startTime).Truncate(time.Second).String()
	}
	return time.Since(m.startTime).Truncate(time.Second).String()
}

func drawBoard(cells []Cell, width, height, startX, startY int, grid [][]uint8) {
	for y := range len(grid[0]) {
		for x := range len(grid) {
			color := defaultBg
			if grid[x][y] > 0 {
				color = blockColor(grid[x][y])
			}
			setCell(cells, width, height, startX+x*2, startY+y, " ", defaultFg, color, false)
			setCell(cells, width, height, startX+x*2+1, startY+y, " ", defaultFg, color, false)
		}
	}
}

func drawPlaytime(cells []Cell, width, height, x, y, boxWidth int, text string) {
	writeText(cells, width, height, x+maxInt((boxWidth-len(text))/2, 0), y, text, mutedFg(), defaultBg, false)
}

func drawScore(cells []Cell, width, height, x, y int, score uint64) {
	drawBox(cells, width, height, x, y, 16, 4, "Score")
	writeText(cells, width, height, x+2, y+2, strconv.FormatUint(score, 10), defaultFg, defaultBg, true)
}

func drawNextBlock(cells []Cell, width, height, x, y int, block *Block) {
	drawBox(cells, width, height, x, y, 16, 8, "Next")
	if block == nil {
		return
	}
	for bx := range len(block.shape) {
		for by := range len(block.shape[0]) {
			if !block.shape[bx][by] {
				continue
			}
			color := blockColor(block.color)
			cellX := x + 4 + bx*2
			cellY := y + 3 + by
			setCell(cells, width, height, cellX, cellY, " ", defaultFg, color, false)
			setCell(cells, width, height, cellX+1, cellY, " ", defaultFg, color, false)
		}
	}
}

func drawHelp(cells []Cell, width, height, x, y int) {
	lines := []string{"Controls", "←/a/h left", "→/d/l right", "↓/s/j drop", "↑/w/r rotate", "space pause"}
	drawBox(cells, width, height, x, y, 20, len(lines)+3, "Help")
	for i, line := range lines {
		writeText(cells, width, height, x+2, y+1+i, line, defaultFg, defaultBg, i == 0)
	}
}

func drawGameOver(cells []Cell, width, height, x, y, boardWidth int) {
	message := "Game Over"
	boxWidth := maxInt(len(message)+4, 12)
	boxX := x + maxInt((boardWidth-boxWidth)/2, 0)
	drawBox(cells, width, height, boxX, y, boxWidth, 5, "")
	writeText(cells, width, height, boxX+maxInt((boxWidth-len(message))/2, 0), y+2, message, "#f87171", defaultBg, true)
}

func drawBox(cells []Cell, width, height, x, y, boxWidth, boxHeight int, title string) {
	if boxWidth < 2 || boxHeight < 2 {
		return
	}
	for dx := 0; dx < boxWidth; dx++ {
		setCell(cells, width, height, x+dx, y, "─", borderFg, defaultBg, false)
		setCell(cells, width, height, x+dx, y+boxHeight-1, "─", borderFg, defaultBg, false)
	}
	for dy := 0; dy < boxHeight; dy++ {
		setCell(cells, width, height, x, y+dy, "│", borderFg, defaultBg, false)
		setCell(cells, width, height, x+boxWidth-1, y+dy, "│", borderFg, defaultBg, false)
	}
	setCell(cells, width, height, x, y, "╭", borderFg, defaultBg, false)
	setCell(cells, width, height, x+boxWidth-1, y, "╮", borderFg, defaultBg, false)
	setCell(cells, width, height, x, y+boxHeight-1, "╰", borderFg, defaultBg, false)
	setCell(cells, width, height, x+boxWidth-1, y+boxHeight-1, "╯", borderFg, defaultBg, false)
	if title != "" {
		writeText(cells, width, height, x+2, y, " "+title+" ", "#7dd3fc", defaultBg, true)
	}
}

func writeText(cells []Cell, width, height, x, y int, text, fg, bg string, bold bool) {
	for offset, r := range []rune(text) {
		setCell(cells, width, height, x+offset, y, string(r), fg, bg, bold)
	}
}

func setCell(cells []Cell, width, height, x, y int, ch, fg, bg string, bold bool) {
	if x < 0 || y < 0 || x >= width || y >= height {
		return
	}
	attrs := []string(nil)
	if bold {
		attrs = []string{"bold"}
	}
	cells[y*width+x] = Cell{X: x, Y: y, Ch: ch, Fg: fg, Bg: bg, Attrs: attrs}
}

func blockColor(color uint8) string {
	idx := int(color) - 1
	if idx < 0 || idx >= len(colors) {
		return defaultBg
	}
	return colors[idx]
}

func mutedFg() string {
	return "#94a3b8"
}

func maxInt(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}
