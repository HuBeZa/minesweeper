package minesweeper

var drawHeader = true

type Minefield interface {
	Flag(row, col int) (Coordinates, error)
	Unflag(row, col int) (Coordinates, error)
	ToggleFlag(row, col int) (Coordinates, error)
	Dig(row, col int) ([]Coordinates, error)
	GameStatus() GameStatus
	CellStatus(row, col int) CellStatus
	AllCellStatus() [][]CellStatus
	FlagsLeft() int
	Width() int
	Height() int
}

type minefield struct {
	width  int
	height int
	cells  [][]*cell
	mines  []Coordinates
	flags map[Coordinates]struct{}
	dugCount int
	status   GameStatus
}

func NewMinefield(width, height int, mines MineList) Minefield {
	f := &minefield{
		width:  width,
		height: height,
		mines:  mines.Coordinates(),
		flags: make(map[Coordinates]struct{}, mines.Len()),
		cells: make([][]*cell, height),
	}

	// init cells & mines
	for row := range f.cells {
		f.cells[row] = make([]*cell, width)
		for col := range f.cells[row] {
			isMine := mines.IsMine(row, col)
			f.cells[row][col] = newCell(isMine)
		}
	}

	// init minesAround counter
	f.initMinesAround()

	return f
}

func (f *minefield) initMinesAround() {
	for _, coord := range f.mines {
		f.incrementMinesAround(coord)
	}
}

func (f *minefield) incrementMinesAround(mineCoord Coordinates) {
	for _, cell := range f.getSurroundingCells(mineCoord) {
		f.cells[cell.Row][cell.Col].minesAround++
	}
}

func (f *minefield) getSurroundingCells(mineCoord Coordinates) []Coordinates {
	surroundingCells := []Coordinates{
		// row above
		{Row: mineCoord.Row - 1, Col: mineCoord.Col - 1},
		{Row: mineCoord.Row - 1, Col: mineCoord.Col},
		{Row: mineCoord.Row - 1, Col: mineCoord.Col + 1},
		// same row
		{Row: mineCoord.Row, Col: mineCoord.Col - 1},
		{Row: mineCoord.Row, Col: mineCoord.Col + 1},
		// row below
		{Row: mineCoord.Row + 1, Col: mineCoord.Col - 1},
		{Row: mineCoord.Row + 1, Col: mineCoord.Col},
		{Row: mineCoord.Row + 1, Col: mineCoord.Col + 1},
	}

	res := make([]Coordinates, 0, len(surroundingCells))
	for _, cell := range surroundingCells {
		// skip out of bounds cells
		if cell.Row >= 0 && cell.Row < f.height &&
			cell.Col >= 0 && cell.Col < f.width {
			res = append(res, cell)
		}
	}

	return res
}

func (f *minefield) Flag(row, col int) (Coordinates, error) {
	if f.status != GameOn {
		return Coordinates{-1, -1}, &GameOverError{}
	}
	if len(f.flags) == len(f.mines) {
		// number of flags cannot exceed the number of mines
		return Coordinates{-1, -1}, &OutOfFlagsError{}
	}
	exists, cell := f.getCell(row, col)
	if !exists {
		return Coordinates{-1, -1}, &InvalidCoordinatesError{}
	}
	if cell.isFlagged {
		return Coordinates{-1, -1}, &AlreadyFlaggedError{}
	}
	if cell.isDug {
		return Coordinates{-1, -1}, &AlreadyDuggedError{}
	}

	cell.isFlagged = true
	coord := Coordinates{row, col}
	f.flags[coord] = struct{}{}
	return coord, nil
}

func (f *minefield) Unflag(row, col int) (Coordinates, error) {
	if f.status != GameOn {
		return Coordinates{-1, -1}, &GameOverError{}
	}
	exists, cell := f.getCell(row, col)
	if !exists {
		return Coordinates{-1, -1}, &InvalidCoordinatesError{}
	}
	if !cell.isFlagged {
		return Coordinates{-1, -1}, &AlreadyUnflaggedError{}
	}
	if cell.isDug {
		return Coordinates{-1, -1}, &AlreadyDuggedError{}
	}

	cell.isFlagged = false
	coord := Coordinates{row, col}
	delete(f.flags, coord)
	return coord, nil
}

func (f *minefield) ToggleFlag(row, col int) (Coordinates, error) {
	if f.status != GameOn {
		return Coordinates{-1, -1}, &GameOverError{}
	}
	exists, cell := f.getCell(row, col)
	if !exists {
		return Coordinates{-1, -1}, &InvalidCoordinatesError{}
	}
	if !cell.isFlagged {
		return f.Flag(row, col)
	}
	return f.Unflag(row, col)
}

func (f *minefield) Dig(row, col int) ([]Coordinates, error) {
	if f.status != GameOn {
		return nil, &GameOverError{}
	}
	exists, cell := f.getCell(row, col)
	if !exists {
		return nil, &InvalidCoordinatesError{}
	}
	if cell.isFlagged {
		return nil, &AlreadyFlaggedError{}
	}
	if !f.digOne(cell, row, col) {
		return nil, &AlreadyDuggedError{}
	}

	// game lost
	if cell.isMine {
		f.status = Lost

		// collect changed cells - the dugged cell + wrong flagged cells + unflagged mines
		changes := []Coordinates{{row,col}}
		changes = append(changes, f.getWronglyFlaggedCells()...)
		changes = append(changes, f.getUnflaggedMines()...)
		return changes, nil
	}

	// expose safe cells
	dugged := f.autoDig(cell, Coordinates{Row: row, Col: col})

	// winning condition - all non-mine cell are dug
	if f.dugCount == f.width*f.height-len(f.mines) {
		f.status = Won

		// collect changed cells - the dugged cell + auto dugged cells + unflagged mines
		return append(dugged, f.getUnflaggedMines()...), nil 
	}

	return dugged, nil
}

func (f *minefield) getWronglyFlaggedCells() []Coordinates {
	wrongFlags := make([]Coordinates, 0)
	for flagCoord := range f.flags {
		if found, flagCell := f.getCell(flagCoord.Row, flagCoord.Col); found && !flagCell.isMine {
			wrongFlags = append(wrongFlags, flagCoord)
		}
	}
	return wrongFlags
}

func (f *minefield) getUnflaggedMines() []Coordinates {
	unflagged := make([]Coordinates, 0)
	for _, mineCoord := range f.mines {
		if found, mineCell := f.getCell(mineCoord.Row, mineCoord.Col); found && !mineCell.isDug && !mineCell.isFlagged {
			unflagged = append(unflagged, mineCoord)
		}
	}
	return unflagged
}

func (f *minefield) digOne(cell *cell, row, col int) bool {
	if cell.isDug {
		return false
	}

	f.Unflag(row, col)
	cell.isDug = true
	f.dugCount++
	return true
}

// autoDig - if cell don't contain a mine and not surrounded by mines, continue auto-digging its surroundings recursively
func (f *minefield) autoDig(cell *cell, coord Coordinates) []Coordinates {
	dugged := []Coordinates{coord}
	f.autoDigRecursive(cell, coord, dugged)
	return dugged
}

func (f *minefield) autoDigRecursive(cell *cell, coord Coordinates, dugged []Coordinates) {
	if cell.minesAround > 0 {
		return
	}

	for _, currCoord := range f.getSurroundingCells(coord) {
		currCell := f.cells[currCoord.Row][currCoord.Col]
		if f.digOne(currCell, currCoord.Row, currCoord.Col) {
			dugged = append(dugged, currCoord)
			f.autoDigRecursive(currCell, currCoord, dugged)
		}
	}
}

func (f *minefield) GameStatus() GameStatus {
	return f.status
}

func (f *minefield) AllCellStatus() [][]CellStatus {
	res := make([][]CellStatus, f.height)
	for row := range f.cells {
		res[row] = make([]CellStatus, f.width)
		for col := range f.cells[row] {
			res[row][col] = f.CellStatus(row, col)
		}
	}
	return res
}

func (f *minefield) CellStatus(row, col int) CellStatus {
	exists, c := f.getCell(row, col)
	if !exists {
		return Unknown
	}

	if f.status == Lost {
		if c.isFlagged {
			if c.isMine {
				return Flagged
			}
			return FlaggedWrong
		}
		if c.isMine {
			if c.isDug {
				return Explode
			}
			return Mine
		}
	} else if f.status == Won {
		// auto-flag mines if won
		if c.isMine {
			return Flagged
		}
	}

	if c.isFlagged {
		return Flagged
	}
	if !c.isDug {
		return Undugged
	}
	if c.minesAround == 0 {
		return NoMinesAround
	}
	return CellStatus(c.minesAround)
}

func (f *minefield) FlagsLeft() int {
	return len(f.mines) - len(f.flags)
}

func (f *minefield) Width() int {
	return f.width
}

func (f *minefield) Height() int {
	return f.height
}

func (f *minefield) getCell(row, col int) (bool, *cell) {
	if row < 0 || row >= f.height || col < 0 || col >= f.width {
		return false, nil
	}
	return true, f.cells[row][col]
}
