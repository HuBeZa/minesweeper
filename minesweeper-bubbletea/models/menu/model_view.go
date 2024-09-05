package menu

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

const (
	// inputs enum
	width = iota
	height
	mineCount
)

var (
	headerStyle = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Margin(0, 1).Padding(0, 3).
			Background(lipgloss.Color("#4b0082")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	optionStyle         = lipgloss.NewStyle().Margin(0, 1)
	selectedOptionStyle = optionStyle.Foreground(lipgloss.Color("#008080"))
	inputRowStyle       = lipgloss.NewStyle().MarginLeft(3)
	labelStyle          = lipgloss.NewStyle().MarginRight(1).Foreground(lipgloss.Color("#808080")).Italic(true)
	errorStyle          = inputRowStyle.Foreground(lipgloss.Color("#d70000"))
	helpStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).MarginTop(1).MarginLeft(1)

	inputLabels = []string{"width", "height", "mines"}
)

func initInputs() []textinput.Model {
	inputs := make([]textinput.Model, 3)

	// set inputs common options
	for i := range inputs {
		input := textinput.New()
		input.CharLimit = 2
		input.Prompt = ""
		input.Placeholder = "___"
		input.Width = len(inputLabels[i])
		inputs[i] = input
	}

	// set specific options
	inputs[mineCount].CharLimit = 3

	return inputs
}

func (m model) View() string {
	rows := make([]string, 0, len(options)+5)
	rows = append(rows, m.renderHeader())
	rows = append(rows, m.renderOptions()...)
	rows = append(rows, m.renderCustomOptions()...)
	rows = append(rows, m.renderErrors()...)
	rows = append(rows, m.renderHelp())

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m model) renderHeader() string {
	return headerStyle.Render("Minesweeper")
}

func (m model) renderOptions() []string {
	rows := make([]string, len(options))
	for i := range options {
		format := "  %s"
		style := optionStyle
		if i == m.selected {
			format = "• %s"
			style = selectedOptionStyle
		}
		rows[i] = style.Render(fmt.Sprintf(format, options[i]))
	}
	return rows
}

func (m model) renderCustomOptions() []string {
	if m.selectedOption() != custom {
		return nil
	}

	labels := make([]string, len(m.inputs))
	inputs := make([]string, len(m.inputs))
	for i := range m.inputs {
		labels[i] = labelStyle.Render(inputLabels[i])
		inputs[i] = m.inputs[i].View()
	}

	return []string{
		inputRowStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, labels...)),
		inputRowStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, inputs...)),
	}
}

func (m model) renderErrors() []string {
	if m.selectedOption() != custom || m.inputError == nil {
		return nil
	}

	return []string{
		errorStyle.Render(m.inputError.Error()),
	}
}

func (m model) renderHelp() string {
	return helpStyle.Render("↑↓←→: navigate • enter: select • ctrl-q: exit")
}
