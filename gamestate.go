package tetris

import "math"

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
			c.grid[i+bufferHor][j+bufferVer] = c.blocks[i][j]
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
							c.grid[i+a+bufferHor][j+b+bufferVer] = c.currentBlock.color
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

	c.currentBlock = &Block{
		x:     2,
		y:     2,
		color: 2,
		shape: [][]bool{{false, true, false}, {true, true, true}, {false, false, false}},
	}
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
	c.ProcessRows()
	c.Drop()
}

func (c Core) Ready() bool {
	return len(c.grid) > 0
}
