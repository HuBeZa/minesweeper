package game

import (
	"fmt"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"

	"github.com/HuBeZa/minesweeper/minesweeper"
)

type model struct {
	field       minesweeper.Minefield
	sw          stopwatch.Model
	zone        *zone.Manager
	started     bool
	pressedCell *minesweeper.Coordinates
}

func NewModel(field minesweeper.Minefield) tea.Model {
	return model{
		field: field,
		sw:    stopwatch.New(),
		zone:  zone.New(),
	}
}

// cellId is used by bubblezone to corelate between mouse clicks to minefield cells
func cellId(row, col int) string {
	return fmt.Sprintf("%v.%v", row, col)
}
