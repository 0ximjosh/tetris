package tetris

// Block is the data for a single tetris block
// x and y are the top left coord
// Each block is in a grid
type Block struct {
	x     uint8
	y     uint8
	shape [][]bool
	color uint8
}

// Rotate rotates the shape 90 degrees
func (b *Block) Rotate() {
	tmp := make([][]bool, len(b.shape))
	for i := range len(b.shape) {
		tmp[i] = make([]bool, len(b.shape[0]))
		for j := range len(b.shape[0]) {
			tmp[i][len(b.shape[0])-j-1] = b.shape[i][j]
		}
	}
	b.shape = make([][]bool, len(tmp[0]))
	for i := range len(tmp[0]) {
		b.shape[i] = make([]bool, len(tmp))
		for j := range len(tmp) {
			b.shape[i][j] = tmp[j][i]
		}
	}
	if b.y+uint8(len(b.shape[0])) >= 20 {
		b.y = uint8(20 - len(b.shape[0]))
	}
}

func (b *Block) Move(d string) {
	switch d {
	case "h":
		if b.x == 0 {
			return
		}
		b.x--
	case "j":
		if b.y+uint8(len(b.shape[0])) >= 20 {
			return
		}
		b.y++
	case "l":
		if b.x+uint8(len(b.shape)) >= 10 {
			return
		}
		b.x++
	}
}
