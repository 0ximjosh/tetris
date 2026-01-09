package tetris

type Core struct {
	grid             [][]uint8
	width            int
	height           int
	paused           bool
	Fps              int
	hideHelp         bool
	currentBlock     *Block
	pendingPlacement bool
	blocks           [][]uint8
	tick             int
}
