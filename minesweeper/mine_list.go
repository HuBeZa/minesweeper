package minesweeper

import (
	"fmt"
	"maps"
	"slices"

	"math/rand"
)

type Coordinates struct {
	Row, Col int
}

func indexToCoordinates(index, width int) Coordinates {
	return Coordinates{
		Row: index / width,
		Col: index % width,
	}
}

type MineList interface {
	Randomize(minesCount int)
	Add(row, col int) error
	Len() int
	IsMine(row, col int) bool
	Coordinates() []Coordinates
}

type mineList struct {
	width   int
	height  int
	mineSet map[Coordinates]struct{}
}

func NewMineList(width, height, capacity int) MineList {
	return &mineList{
		width:   width,
		height:  height,
		mineSet: make(map[Coordinates]struct{}, capacity),
	}
}

func (l *mineList) Randomize(minesCount int) {
	initialLen := l.Len()
	maxIndex := l.width * l.height

	for l.Len()-initialLen < minesCount {
		index := rand.Intn(maxIndex)
		l.mineSet[indexToCoordinates(index, l.width)] = struct{}{}
	}
}

func (l *mineList) Add(row, col int) error {
	if row < 0 || row >= l.height || col < 0 || col >= l.width {
		return fmt.Errorf("invalid mine coordinates")
	}

	l.mineSet[Coordinates{Row: row, Col: col}] = struct{}{}
	return nil
}

func (l *mineList) Len() int {
	return len(l.mineSet)
}

func (l *mineList) IsMine(row, col int) bool {
	_, isMine := l.mineSet[Coordinates{Row: row, Col: col}]
	return isMine
}

func (l *mineList) Coordinates() []Coordinates {
	return slices.Collect(maps.Keys(l.mineSet))
}
