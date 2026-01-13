package tetris

type Core struct {
	grid             [][]uint8
	width            int
	height           int
	paused           bool
	Fps              int
	hideHelp         bool
	currentBlock     *Block
	nextBlock        *Block
	pendingPlacement bool
	blocks           [][]uint8
	tickSpeed        int
	tick             int
	score            uint64
	gameOver         bool
}
