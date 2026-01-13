package tetris

import "time"

type Core struct {
	grid   [][]uint8
	width  int
	height int
	paused bool
	Fps    int
	// hideHelp         bool TODO
	currentBlock     *Block
	nextBlock        *Block
	pendingPlacement bool
	blocks           [][]uint8
	tickSpeed        int
	tick             int
	score            uint64
	gameOver         bool
	startTime        time.Time
	endTime          time.Time
}
