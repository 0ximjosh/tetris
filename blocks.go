package tetris

var (
	t  = Block{x: 4, y: 2, color: 2, shape: [][]bool{{false, true, false}, {true, true, true}}}
	l  = Block{x: 4, y: 2, color: 3, shape: [][]bool{{false, false, false, false}, {true, true, true, true}, {false, false, false, false}, {false, false, false, false}}}
	j  = Block{x: 4, y: 2, color: 4, shape: [][]bool{{true, false, false}, {true, true, true}, {false, false, false}}}
	j2 = Block{x: 4, y: 2, color: 5, shape: [][]bool{{false, false, true}, {true, true, true}, {false, false, false}}}
	o  = Block{x: 4, y: 2, color: 6, shape: [][]bool{{true, true}, {true, true}}}
	z  = Block{x: 4, y: 2, color: 7, shape: [][]bool{{false, true, true}, {true, true, false}, {false, false, false}}}
	z2 = Block{x: 4, y: 2, color: 8, shape: [][]bool{{true, true, false}, {false, true, true}, {false, false, false}}}
)

var Blocks = []Block{t, l, j, j2, o, z, z2}
