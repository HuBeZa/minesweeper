package messages

import (
	"github.com/HuBeZa/minesweeper/minesweeper"
	tea "github.com/charmbracelet/bubbletea"
)

type ShowMenuMsg struct{}

func ShowMenu() tea.Msg {
	return ShowMenuMsg{}
}


type StartNewGameMsg struct {
	Minefield minesweeper.Minefield
}

type MenuErrorMsg struct {
	Err error
}
