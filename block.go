package tetris

import "math/rand/v2"

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
func (c *Core) Rotate() {
	if c.gameOver {
		return
	}
	b := c.currentBlock
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
	if c.CanPlace(tmpBlock) {
		b.shape = tmp2
		return
	}
	tmpBlock.x = b.x + 1
	if c.CanPlace(tmpBlock) {
		b.shape = tmp2
		b.x += 1
		return
	}

	tmpBlock.x = b.x + 2
	if c.CanPlace(tmpBlock) {
		b.shape = tmp2
		b.x += 2
		return
	}

	tmpBlock.x = b.x - 1
	if c.CanPlace(tmpBlock) {
		b.shape = tmp2
		b.x -= 1
		return
	}
}

func (c *Core) MoveBlock(d string) {
	if c.gameOver {
		return
	}
	switch d {
	case "h":
		tmp := *c.currentBlock
		tmp.x -= 1
		if c.CanPlace(tmp) {
			c.currentBlock.x--
		}
	case "l":
		tmp := *c.currentBlock
		tmp.x += 1
		if c.CanPlace(tmp) {
			c.currentBlock.x++
		}
	}
}

func (c *Core) CanPlace(b Block) bool {
	// Collision with placed blocks and walls
	for i := range len(b.shape) {
		for j := range len(b.shape[0]) {
			if !b.shape[i][j] {
				continue
			}
			if i+b.x < 0 {
				return false
			}
			if c.blocks[i+b.x][j+b.y] != 0 {
				return false
			}
		}
	}
	return true
}

func (c *Core) NextBlock() {
	b := Blocks[rand.IntN(len(Blocks))]
	c.currentBlock = c.nextBlock
	c.nextBlock = &b
	if !c.CanPlace(*c.nextBlock) {
		c.gameOver = true
	}
}

func (c *Core) Drop() {
	if c.gameOver {
		return
	}
	b := *c.currentBlock
	b.y++
	canDrop := c.CanPlace(b)

	switch true {
	case canDrop:
		c.pendingPlacement = false
		c.currentBlock.y++
	case c.pendingPlacement && !canDrop:
		c.PlaceCurrentBlock()
		c.NextBlock()
		c.pendingPlacement = false
	case !c.pendingPlacement && !canDrop:
		c.pendingPlacement = true
	}
}

func (c *Core) ProcessRows() {
	shift := 0
	for yr := range len(c.blocks[0]) - 1 {
		y := len(c.blocks[0]) - yr - 2
		score := true
		for x := range len(c.blocks) - 2 {
			if shift != 0 {
				c.blocks[x+1][y+shift] = c.blocks[x+1][y]
			}
			if c.blocks[x+1][y] == 0 {
				score = false
			}
		}
		if score {
			shift++
		}
	}
	c.score += uint64(shift * shift)
}

func (c *Core) PlaceCurrentBlock() {
	for i := range len(c.currentBlock.shape) {
		for j := range len(c.currentBlock.shape[0]) {
			if !c.currentBlock.shape[i][j] {
				continue
			}
			c.blocks[i+c.currentBlock.x][j+c.currentBlock.y] = c.currentBlock.color
		}
	}
}
