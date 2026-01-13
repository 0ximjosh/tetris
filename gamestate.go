package tetris

import (
	"math/rand/v2"
	"time"
)

// Init the game state
// Setup tetris bucket
func (c *Core) Init(w, h int) {
	if w < 0 {
		return
	}
	c.width = w
	c.height = h
	c.tickSpeed = 20
	c.startTime = time.Now()

	c.grid = make([][]uint8, 12)
	c.blocks = make([][]uint8, 12)
	for i := range c.blocks {
		c.blocks[i] = make([]uint8, 21)
		c.grid[i] = make([]uint8, 21)
		c.blocks[i][20] = 1
		if i == 0 || i == len(c.blocks)-1 {
			for j := range 21 {
				c.blocks[i][j] = 1
			}
		}
	}

	b := Blocks[rand.IntN(len(Blocks))]
	c.nextBlock = &b
	b2 := Blocks[rand.IntN(len(Blocks))]
	c.currentBlock = &b2
}

// Update the game grid
// Core game logic. Update core.grid accordingly for this frame
func (c *Core) Update() {
	if !c.Ready() {
		return
	}
	if c.gameOver {
		return
	}
	c.tick++
	if c.tick%c.tickSpeed != 0 {
		return
	}
	if c.tick%100 == 0 && c.tickSpeed != 1 {
		c.tickSpeed -= 1
	}
	c.Drop()
	c.ProcessRows()
}

func (c Core) Ready() bool {
	return len(c.grid) > 0
}
