package minesweeper

import (
	"fmt"
)

var instance = &generator{}

type Generator interface {
	Custom(width, height, minesCount int) (Minefield, error)
	Beginner() Minefield
	Intermediate() Minefield
	Expert() Minefield
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

func (g *generator) Beginner() Minefield {
	f, _ := g.Custom(9, 9, 10)
	return f
}

func (g *generator) Intermediate() Minefield {
	f, _ := g.Custom(16, 16, 40)
	return f
}

func (g *generator) Expert() Minefield {
	f, _ := g.Custom(30, 16, 99)
	return f
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
