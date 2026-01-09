package tetris

import (
	"math"
	"math/rand/v2"
)

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
	bufferHor := int(math.Floor(float64(c.width)/2) - 6)
	bufferVer := int(math.Floor(float64(c.height)/2) - 10)
	for i := range 12 {
		for j := range 21 {
			c.grid[i+bufferHor][j+bufferVer] = c.blocks[i][j]
		}
	}

	if c.currentBlock != nil {
		for a := range len(c.currentBlock.shape) {
			for b := range len(c.currentBlock.shape[0]) {
				if !c.currentBlock.shape[a][b] {
					continue
				}
				c.grid[c.currentBlock.x+a+bufferHor][c.currentBlock.y+b+bufferVer] = c.currentBlock.color
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
	c.width = w / 2
	c.height = h
	c.grid = make([][]uint8, w)
	for i := range c.grid {
		c.grid[i] = make([]uint8, h)
	}
	c.blocks = make([][]uint8, 12)
	for i := range c.blocks {
		c.blocks[i] = make([]uint8, 21)
		c.blocks[i][20] = 1
		if i == 0 || i == len(c.blocks)-1 {
			for j := range 21 {
				c.blocks[i][j] = 1
			}
		}
	}

	b := Blocks[rand.IntN(len(Blocks))]

	c.currentBlock = &b
	c.WriteState()
}

// Update the game grid
// Core game logic. Update core.grid accordingly for this frame
func (c *Core) Update() {
	if !c.Ready() {
		return
	}
	defer c.WriteState()

	c.tick++
	if c.tick%10 != 0 {
		return
	}
	c.Drop()
	c.ProcessRows()
}

func (c Core) Ready() bool {
	return len(c.grid) > 0
}
