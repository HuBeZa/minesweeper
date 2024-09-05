package menu

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type gameType string

const (
	beginner     gameType = "Beginner"
	intermediate gameType = "Intermediate"
	expert       gameType = "Expert"
	custom       gameType = "Custom"
)

var options = []gameType{beginner, intermediate, expert, custom}

type model struct {
	selected int

	inputs       []textinput.Model
	focusedInput int
	inputError error
}

func NewModel() tea.Model {
	return model{
		inputs:       initInputs(),
		focusedInput: -1,
	}
}

func (m model) selectedOption() gameType {
	return options[m.selected]
}
