package game

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"

	"github.com/HuBeZa/minesweeper/minesweeper"
)

var (
	leftHeaderStyle  = lipgloss.NewStyle().AlignHorizontal(lipgloss.Left).PaddingLeft(1)
	rightHeaderStyle = lipgloss.NewStyle().AlignHorizontal(lipgloss.Right).PaddingRight(1)

	winMessageStyle  = lipgloss.NewStyle().Background(lipgloss.Color("#4aa45b")).Foreground(lipgloss.Color("#FFFFFF")).Bold(true).AlignHorizontal(lipgloss.Center)
	loseMessageStyle = winMessageStyle.Background(lipgloss.Color("#ff0000"))
	helpStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).MarginTop(1).MarginLeft(1)

	fieldStyle       = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).Padding(0, 1)
	cellStyle        = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, false, true).Padding(0, 1)
	topCellStyle     = cellStyle.UnsetBorderTop()
	leftCellStyle    = cellStyle.UnsetBorderLeft()
	topLeftCellStyle = leftCellStyle.UnsetBorderTop()

	boldRedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Bold(true)

	cellStatusToString = map[minesweeper.CellStatus]string{
		minesweeper.Undugged:      "■",
		minesweeper.NoMinesAround: colorString("□", "#808080"),
		minesweeper.MinesAround1:  colorString("1", "#1e90ff"),
		minesweeper.MinesAround2:  colorString("2", "#008000"),
		minesweeper.MinesAround3:  colorString("3", "#b22222"),
		minesweeper.MinesAround4:  colorString("4", "#4b0082"),
		minesweeper.MinesAround5:  colorString("5", "#a52a2a"),
		minesweeper.MinesAround6:  colorString("6", "#008080"),
		minesweeper.MinesAround7:  colorString("7", "#afeeee"),
		minesweeper.MinesAround8:  colorString("8", "#ffa500"),
		minesweeper.Flagged:       colorString("?", "#ff6347"),
		minesweeper.FlaggedWrong:  boldRedStyle.Render("X"),
		minesweeper.Mine:          colorString("M", "#ff0000"),
		minesweeper.Explode:       boldRedStyle.Render("Ж"),
	}
)

func (m model) View() string {
	field := m.renderField()
	fieldWidth := realWidthOf(field)
	header := m.renderHeader(fieldWidth)
	footer := m.renderFooter(fieldWidth)

	return m.zone.Scan(
		lipgloss.JoinVertical(lipgloss.Left, header, field, footer))
}

func (m model) renderField() string {
	cells := m.field.AllCellStatus()
	tableView := make([]string, len(cells))
	for row := range cells {
		tableView[row] = m.renderRow(row, cells[row])
	}

	return fieldStyle.Render(lipgloss.JoinVertical(lipgloss.Left, tableView...))
}

func (m model) renderRow(row int, cells []minesweeper.CellStatus) string {
	cellsStr := make([]string, len(cells))
	for col := range cells {
		cellsStr[col] = m.renderCell(row, col, cells[col])
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, cellsStr...)
}

func (m model) renderCell(row, col int, cellStatus minesweeper.CellStatus) string {
	style := getCellStyle(row, col)
	if m.pressedCell.Equals(row, col) && (cellStatus == minesweeper.Undugged || cellStatus == minesweeper.Flagged) {
		style = style.Reverse(true)
	}
	cellStr := cellStatusToString[cellStatus]
	cellStr = style.Render(cellStr)
	cellZoneId := cellId(row, col)
	return m.zone.Mark(cellZoneId, cellStr)
}

func getCellStyle(row, col int) lipgloss.Style {
	if col == 0 {
		if row == 0 {
			return topLeftCellStyle
		}
		return leftCellStyle
	}
	if row == 0 {
		return topCellStyle
	}
	return cellStyle
}

func (m model) renderHeader(width int) string {
	leftHeader := fmt.Sprintf("Flags: %v", m.field.FlagsLeft())
	rightHeader := m.sw.View()

	leftHeader = leftHeaderStyle.Width(width / 2).Render(leftHeader)
	if width%2 != 0 {
		width++
	}
	rightHeader = rightHeaderStyle.Width(width / 2).Render(rightHeader)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftHeader, rightHeader)
}

func (m model) renderFooter(width int) string {
	rows := make([]string, 0, 2)
	if m.field.GameStatus() == minesweeper.Won {
		rows = append(rows, winMessageStyle.Width(width).Render("YOU WON"))
	} else if m.field.GameStatus() == minesweeper.Lost {
		rows = append(rows, loseMessageStyle.Width(width).Render("YOU LOST"))
	}

	rows = append(rows, helpStyle.Render("ctrl-n: new game • ctrl-q: exit"))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func realWidthOf(s string) int {
	width := strings.IndexRune(s, '\n')
	if width < 0 {
		return utf8.RuneCountInString(s)
	}

	return utf8.RuneCountInString(s[:width])
}

func colorString(s, color string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(s)
}
