package tetris

import (
	"math/rand/v2"
	"time"
)

// Block is the data for a single tetris block
// x and y are the top left coord
// Each block is in a grid
type Block struct {
	x     int
	y     int
	shape [][]bool
	color uint8
}

// Rotate rotates the shape 90 degrees
func (m *Model) Rotate() {
	if m.gameOver {
		return
	}
	b := m.currentBlock
	tmp := make([][]bool, len(b.shape))
	for i := range len(b.shape) {
		tmp[i] = make([]bool, len(b.shape[0]))
		for j := range len(b.shape[0]) {
			tmp[i][len(b.shape[0])-j-1] = b.shape[i][j]
		}
	}
	tmp2 := make([][]bool, len(tmp[0]))
	for i := range len(tmp[0]) {
		tmp2[i] = make([]bool, len(tmp))
		for j := range len(tmp) {
			tmp2[i][j] = tmp[j][i]
		}
	}
	tmpBlock := Block{
		x:     b.x,
		y:     b.y,
		shape: tmp2,
	}
	if m.CanPlace(tmpBlock) {
		b.shape = tmp2
		return
	}
	tmpBlock.x = b.x + 1
	if m.CanPlace(tmpBlock) {
		b.shape = tmp2
		b.x += 1
		return
	}

	tmpBlock.x = b.x + 2
	if m.CanPlace(tmpBlock) {
		b.shape = tmp2
		b.x += 2
		return
	}

	tmpBlock.x = b.x - 1
	if m.CanPlace(tmpBlock) {
		b.shape = tmp2
		b.x -= 1
		return
	}
}

func (m *Model) MoveBlock(d string) {
	if m.gameOver {
		return
	}
	switch d {
	case "left":
		tmp := *m.currentBlock
		tmp.x -= 1
		if m.CanPlace(tmp) {
			m.currentBlock.x--
		}
	case "right":
		tmp := *m.currentBlock
		tmp.x += 1
		if m.CanPlace(tmp) {
			m.currentBlock.x++
		}
	}
}

func (m *Model) CanPlace(b Block) bool {
	// Collision with placed blocks and walls
	for i := range len(b.shape) {
		for j := range len(b.shape[0]) {
			if !b.shape[i][j] {
				continue
			}
			if i+b.x < 0 {
				return false
			}
			if m.blocks[i+b.x][j+b.y] != 0 {
				return false
			}
		}
	}
	return true
}

func (m *Model) NextBlock() {
	b := Blocks[rand.IntN(len(Blocks))]
	m.currentBlock = m.nextBlock
	m.nextBlock = &b
	if !m.CanPlace(*m.nextBlock) {
		m.gameOver = true
		m.endTime = time.Now()
	}
}

func (m *Model) Drop() {
	if m.gameOver {
		return
	}
	b := *m.currentBlock
	b.y++
	canDrop := m.CanPlace(b)

	switch true {
	case canDrop:
		m.pendingPlacement = false
		m.currentBlock.y++
	case m.pendingPlacement && !canDrop:
		m.PlaceCurrentBlock()
		m.NextBlock()
		m.pendingPlacement = false
	case !m.pendingPlacement && !canDrop:
		m.pendingPlacement = true
	}
}

func (m *Model) ProcessRows() {
	shift := 0
	for yr := range len(m.blocks[0]) - 1 {
		y := len(m.blocks[0]) - yr - 2
		score := true
		for x := range len(m.blocks) - 2 {
			if shift != 0 {
				m.blocks[x+1][y+shift] = m.blocks[x+1][y]
			}
			if m.blocks[x+1][y] == 0 {
				score = false
			}
		}
		if score {
			shift++
		}
	}
	m.score += uint64(shift * shift)
}

func (m *Model) PlaceCurrentBlock() {
	for i := range len(m.currentBlock.shape) {
		for j := range len(m.currentBlock.shape[0]) {
			if !m.currentBlock.shape[i][j] {
				continue
			}
			m.blocks[i+m.currentBlock.x][j+m.currentBlock.y] = m.currentBlock.color
		}
	}
}
