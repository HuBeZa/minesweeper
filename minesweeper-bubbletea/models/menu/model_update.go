package menu

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/HuBeZa/minesweeper/minesweeper"
	"github.com/HuBeZa/minesweeper/minesweeper-bubbletea/models/messages"
)

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.MenuErrorMsg:
		m.inputError = msg.Err
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit
		case "down":
			if m.selected < len(options)-1 {
				m.selected++
				return m.onFocusChanged()
			}
		case "up":
			if m.selected > 0 {
				m.selected--
				return m.onFocusChanged()
			}
		case "tab", "right":
			return m.changeFocus(+1)
		case "shift+tab", "left":
			return m.changeFocus(-1)
		case "enter":
			return m, m.generateMinefield
		default:
			return m.updateInputs(msg)
		}
	}

	return m, nil
}

func (m model) changeFocus(focusChange int) (tea.Model, tea.Cmd) {
	if m.selectedOption() != custom {
		return m, nil
	}

	m.focusedInput = (m.focusedInput + focusChange) % len(m.inputs)
	if m.focusedInput < 0 {
		m.focusedInput += len(m.inputs)
	}
	return m.onFocusChanged()
}

func (m model) onFocusChanged() (tea.Model, tea.Cmd) {
	for i := range m.inputs {
		m.inputs[i].Blur()
	}

	if m.selectedOption() == custom {
		if m.focusedInput == -1 {
			m.focusedInput = 0
		}
		m.inputs[m.focusedInput].Focus()
	} else {
		m.focusedInput = -1
	}

	return m, nil
}

func (m model) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.selectedOption() != custom {
		return m, nil
	}

	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) generateMinefield() tea.Msg {
	switch m.selectedOption() {
	case beginner:
		return messages.StartNewGameMsg{Minefield: minesweeper.GameGenerator().Beginner()}
	case intermediate:
		return messages.StartNewGameMsg{Minefield: minesweeper.GameGenerator().Intermediate()}
	case expert:
		return messages.StartNewGameMsg{Minefield: minesweeper.GameGenerator().Expert()}
	case custom:
		values, err := m.getInputValues()
		if err != nil {
			return messages.MenuErrorMsg{Err: err}
		}

		minefield, err := minesweeper.GameGenerator().Custom(values[width], values[height], values[mineCount])
		if err != nil {
			return messages.MenuErrorMsg{Err: err}
		}

		return messages.StartNewGameMsg{Minefield: minefield}
	default:
		return messages.MenuErrorMsg{Err: fmt.Errorf("selection unknown")}
	}
}

func (m model) getInputValues() ([]int, error) {
	values := make([]int, len(m.inputs))
	for i := range m.inputs {
		val, err := strconv.Atoi(strings.TrimSpace(m.inputs[i].Value()))
		if err != nil {
			return nil, fmt.Errorf("%v must be a number", inputLabels[i])
		}

		values[i] = val
	}
	
	return values, nil
}
