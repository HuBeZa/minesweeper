package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/HuBeZa/minesweeper/minesweeper-bubbletea/models/game"
	"github.com/HuBeZa/minesweeper/minesweeper-bubbletea/models/menu"
	"github.com/HuBeZa/minesweeper/minesweeper-bubbletea/models/messages"
)

type model struct {
	activeModel tea.Model
}

func newModel() tea.Model {
	return &model{activeModel: menu.NewModel()}
}

func (m model) Init() tea.Cmd {
	return m.activeModel.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.StartNewGameMsg:
		return m.switchModel(game.NewModel(msg.Minefield))
	case messages.ShowMenuMsg:
		return m.switchModel(menu.NewModel())
	default:
		var cmd tea.Cmd
		m.activeModel, cmd = m.activeModel.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	return m.activeModel.View()
}

func (m model) switchModel(newModel tea.Model) (tea.Model, tea.Cmd) {
	m.activeModel = newModel
	return m, m.activeModel.Init()
}

func main() {
	m := newModel()
	_, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
