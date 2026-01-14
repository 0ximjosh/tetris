package tetris

func (m *Model) UpdateDims(w, h int) {
	m.width = w
	m.height = h
}

func (m Model) Ready() bool {
	return len(m.grid) > 0
}
