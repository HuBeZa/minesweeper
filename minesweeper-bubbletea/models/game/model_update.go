package game

import (
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/HuBeZa/minesweeper/minesweeper"
	"github.com/HuBeZa/minesweeper/minesweeper-bubbletea/models/messages"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case stopwatch.TickMsg, stopwatch.StartStopMsg:
		var cmd tea.Cmd
		m.sw, cmd = m.sw.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+n":
			return m, messages.ShowMenu
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit
		}
	case tea.MouseMsg:
		if m.field.GameStatus() != minesweeper.GameOn ||
			(msg.Action != tea.MouseActionMotion && msg.Action != tea.MouseActionPress && msg.Action != tea.MouseActionRelease) ||
			(msg.Button != tea.MouseButtonLeft && msg.Button != tea.MouseButtonRight) {
			return m, nil
		}

		row, col := m.getClickedCell(msg)
		if msg.Action == tea.MouseActionPress || msg.Action == tea.MouseActionMotion {
			return m.handleMousePressed(row, col)
		}
		return m.handleMouseRelease(msg, row, col)
	}
	return m, nil
}

func (m model) getClickedCell(msg tea.MouseMsg) (int, int) {
	for row := 0; row < m.field.Height(); row++ {
		for col := 0; col < m.field.Width(); col++ {
			if m.zone.Get(cellId(row, col)).InBounds(msg) {
				return row, col
			}
		}
	}
	return -1, -1
}

func (m model) handleMousePressed(row, col int) (tea.Model, tea.Cmd) {
	m.pressedCell = nil
	if row >= 0 {
		m.pressedCell = &minesweeper.Coordinates{Row: row, Col: col}
	}
	return m, nil
}

func (m model) handleMouseRelease(msg tea.MouseMsg, row, col int) (tea.Model, tea.Cmd) {
	m.pressedCell = nil
	if row < 0 {
		return m, nil
	}

	if msg.Button == tea.MouseButtonLeft {
		m.field.Dig(row, col)
	} else if msg.Button == tea.MouseButtonRight {
		m.field.ToggleFlag(row, col)
	}

	if m.field.GameStatus() != minesweeper.GameOn {
		// stop stopwatch on game over
		return m, m.sw.Stop()
	}

	if !m.started {
		// start stopwatch on first move
		m.started = true
		return m, m.sw.Init()
	}
	return m, nil
}
