package minesweeper

type cell struct {
	isMine      bool
	minesAround int
	isFlagged   bool
	isDug       bool
}

func newCell(isMine bool) *cell {
	return &cell{
		isMine: isMine,
	}
}
