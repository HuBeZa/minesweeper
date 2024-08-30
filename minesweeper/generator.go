package minesweeper

import (
	"fmt"
)

var instance = &generator{}

type Generator interface {
	Custom(width, height, minesCount int) (Minefield, error)
	Beginner() (Minefield, error)
	Intermediate() (Minefield, error)
	Expert() (Minefield, error)
}

type generator struct {
}

func GameGenerator() Generator {
	return instance
}

func (g *generator) Custom(width, height, minesCount int) (Minefield, error) {
	if err := g.validate(width, height, minesCount); err != nil {
		return nil, err
	}

	mines := NewMineList(width, height, minesCount)
	mines.Randomize(minesCount)
	

	return NewMinefield(width, height, mines), nil
}

func (g *generator) Beginner() (Minefield, error) {
	return g.Custom(9, 9, 10)
}

func (g *generator) Intermediate() (Minefield, error) {
	return g.Custom(16, 16, 40)
}

func (g *generator) Expert() (Minefield, error) {
	return g.Custom(30, 16, 99)
}

func (g *generator) validate(width, height, minesCount int) error {
	if width < 2 || height < 2 {
		return fmt.Errorf("minefield dimensions should be at least 2x2")
	}
	if minesCount <= 0 {
		return fmt.Errorf("a minefield should have at least one mine")
	}
	if minesCount > width*height {
		return fmt.Errorf("the number of mines cannot exceed the size of the minefield")
	}
	return nil
}
