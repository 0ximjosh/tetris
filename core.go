package tetris

import (
	"math/rand/v2"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
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
	Renderer         *lipgloss.Renderer
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, m.HandleKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.UpdateDims(msg.Width, msg.Height)
		return m, nil
	case FrameMsg:
		m.Tick()
		return m, m.Animate()
	default:
		return m, nil
	}
}

func (m Model) View() string {
	return m.String()
}

func (m Model) Init() tea.Cmd {
	return m.Animate()
}

// Reset the game state
// Setup tetris bucket
func (m *Model) Reset() {
	if m.Renderer == nil {
		m.Renderer = lipgloss.DefaultRenderer()
	}
	m.tickSpeed = 20
	m.startTime = time.Now()

	m.grid = make([][]uint8, 12)
	m.blocks = make([][]uint8, 12)
	for i := range m.blocks {
		m.blocks[i] = make([]uint8, 21)
		m.grid[i] = make([]uint8, 21)
		m.blocks[i][20] = 1
		if i == 0 || i == len(m.blocks)-1 {
			for j := range 21 {
				m.blocks[i][j] = 1
			}
		}
	}

	b := Blocks[rand.IntN(len(Blocks))]
	m.nextBlock = &b
	b2 := Blocks[rand.IntN(len(Blocks))]
	m.currentBlock = &b2
}

// Tick updates the game grid
// Core game logic. Update core.grid accordingly for this frame
func (m *Model) Tick() {
	if !m.Ready() {
		return
	}
	if m.gameOver {
		return
	}
	m.tick++
	if m.tick%m.tickSpeed != 0 {
		return
	}
	if m.tick%100 == 0 && m.tickSpeed != 1 {
		m.tickSpeed -= 1
	}
	m.Drop()
	m.ProcessRows()
}
