package tetris

import (
	"math/rand/v2"
)

// Init the game state
// Setup tetris bucket
func (c *Core) Init(w, h int) {
	if w < 0 {
		return
	}
	c.width = w
	c.height = h

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

	c.currentBlock = &b
}

// Update the game grid
// Core game logic. Update core.grid accordingly for this frame
func (c *Core) Update() {
	if !c.Ready() {
		return
	}

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
