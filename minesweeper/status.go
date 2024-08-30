package minesweeper

type GameStatus int

const (
	GameOn GameStatus = iota
	Lost
	Won
)

type CellStatus int

const (
	Undugged      CellStatus = 0
	MinesAround1             = 1
	MinesAround2             = 2
	MinesAround3             = 3
	MinesAround4             = 4
	MinesAround5             = 5
	MinesAround6             = 6
	MinesAround7             = 7
	MinesAround8             = 8
	NoMinesAround            = 9
	Flagged                  = 10
	FlaggedWrong             = 11
	Mine                     = 12
	Explode                  = 13
	Unknown                  = 14
)
