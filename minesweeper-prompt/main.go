package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/HuBeZa/minesweeper/minesweeper"
	"github.com/fatih/color"
	"github.com/inancgumus/screen"
)

var (
	drawHeader  = true
	headerColor = color.New(color.FgHiBlack, color.Bold)
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		printHelp()
		os.Exit(0)
	}

	field, err := generateMinefield()
	if err != nil {
		color.HiRed("error: %v\n", err)
		printHelp()
		os.Exit(1)
	}

	screen.Clear()
	draw(field)

	reader := bufio.NewReader(os.Stdin)
	for field.GameStatus() == minesweeper.GameOn {
		fmt.Print("Enter command and coordinates in this format: <command row col>\n - flag: 'f 2 1'\n - unflag: 'u 7 0'\n - dig: d 3 5\nYour command: ")
		text, _ := reader.ReadString('\n')
		err := runCommand(field, text)
		screen.Clear()
		draw(field)
		if err != nil {
			fmt.Println()
			color.HiRed("%s\n", err)
		}
	}
}

func generateMinefield() (minesweeper.Minefield, error) {
	if len(os.Args) < 2 {
		return minesweeper.GameGenerator().Beginner()
	}

	switch cmd := strings.ToLower(os.Args[1]); cmd {
	case "beginner", "b":
		return minesweeper.GameGenerator().Beginner()
	case "intermediate", "i":
		return minesweeper.GameGenerator().Intermediate()
	case "expert", "e":
		return minesweeper.GameGenerator().Expert()
	case "custom", "c":
		if len(os.Args) < 5 {
			return nil, fmt.Errorf("not enough arguments for 'custom' command")
		}

		strArgs := os.Args[2:]
		intArgs := make([]int, 3)
		for i, argName := range []string{"width", "height", "mines-count"} {
			intArg, err := strconv.Atoi(strArgs[i])
			if err != nil {
				return nil, fmt.Errorf("wrong type arguments '%v', expect number'", argName)
			}
			intArgs[i] = intArg
		}
		return minesweeper.GameGenerator().Custom(intArgs[0], intArgs[1], intArgs[2])
	default:
		return nil, fmt.Errorf("unknown command '%v'", cmd)
	}
}

func printHelp() {
	fmt.Println("usage: minesweeper-prompt.exe [-h | --help]\n" +
		"                              <command> [<args>]\n" +
		"commands:\n" +
		"\tbeginner | b\n" +
		"\tintermediate | i\n" +
		"\texpert | e\n" +
		"\tcustom | c <width> <height> <mines-count>")
}

func runCommand(field minesweeper.Minefield, input string) error {
	input = strings.TrimSpace(input)

	parts := strings.Split(input, " ")
	if len(parts) != 3 {
		return fmt.Errorf("wrong format")
	}

	cmd := parts[0]
	if cmd != "f" && cmd != "u" && cmd != "d" {
		return fmt.Errorf("wrong command name '%s', use 'f', 'u' or 'd'", parts[0])
	}

	row, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("wrong row coordinate format '%s', use a number", parts[1])
	}
	col, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("wrong col coordinate format '%s', use a number", parts[2])
	}

	switch cmd {
	case "f":
		field.Flag(row, col)
	case "u":
		field.Unflag(row, col)
	case "d":
		field.Dig(row, col)
	}
	return nil
}

func draw(f minesweeper.Minefield) {
	var sb strings.Builder
	cells := f.AllCellStatus()

	if drawHeader {
		sb.WriteString(fmt.Sprintf(" 🚩 = %v\n", f.FlagsLeft()))

		sb.WriteString(headerColor.Sprintf("%8v", "|\t"))
		for col := range cells[0] {
			if col > 0 {
				sb.WriteString(headerColor.Sprint("\t"))
			}

			sb.WriteString(headerColor.Sprintf("%v", col))
		}
		sb.WriteString(headerColor.Sprint("\t| \n"))
	}

	for row := range cells {
		if drawHeader {
			sb.WriteString(headerColor.Sprintf(" %-4v ", row))
		}

		sb.WriteString("|\t")
		for col := range cells[row] {
			if col > 0 {
				sb.WriteString("\t")
			}

			cellStr := cellStatusToString(cells[row][col])
			sb.WriteString(fmt.Sprintf("%v", cellStr))
		}
		sb.WriteString("\t| \n")
	}

	if f.GameStatus() == minesweeper.Lost {
		c := color.New(color.BgRed, color.FgWhite, color.Bold)
		sb.WriteString(c.Sprint("GAME OVER!!! GAME OVER!!! GAME OVER!!! GAME OVER!!! GAME OVER!!! GAME OVER!!!"))
	} else if f.GameStatus() == minesweeper.Won {
		c := color.New(color.FgMagenta, color.Bold)
		sb.WriteString("😎 😎 😎 😎 😎 😎 😎 😎 😎 😎 😎 " + c.Sprint("YOU WON") + " 😎 😎 😎 😎 😎 😎 😎 😎 😎 😎 😎")
	}

	fmt.Println(sb.String())
}

func cellStatusToString(status minesweeper.CellStatus) string {
	switch status {
	case minesweeper.Undugged:
		return "■"
	case minesweeper.NoMinesAround:
		return "□"
	case minesweeper.MinesAround1, minesweeper.MinesAround2, minesweeper.MinesAround3, minesweeper.MinesAround4,
		minesweeper.MinesAround5, minesweeper.MinesAround6, minesweeper.MinesAround7, minesweeper.MinesAround8:
		return strconv.Itoa(int(status))
	case minesweeper.Flagged:
		return "🚩"
	case minesweeper.FlaggedWrong:
		return color.HiRedString("✗")
	case minesweeper.Mine:
		return "💣"
	case minesweeper.Explode:
		return "💥"
	}

	return ""
}
