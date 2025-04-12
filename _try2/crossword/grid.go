package main

const (
	LetterEmpty     Letter = 0
	LetterEndOfWord Letter = 3
)

// unicode.IsPrint'able rune with special values
type Letter rune

// Letters without special values
type Word []Letter

type Question struct {
	word Word
}

type Position struct {
	g     *Grid
	index int
}

type cellFlags byte

const (
	cellFlagWord cellFlags = 1 << iota
	cellFlagHWord
	cellFlagVWord
	cellFlagHWordStart
	cellFlagVWordStart
	cellFlagHole // letter can't be placed
	cellFlagHWordTerminator
	cellFlagVWordTerminator
)

type cell struct {
	letter Letter
	flags  cellFlags
}

type Config struct {
	Width, Height int
}

type Grid struct {
	cfg      *Config
	parent   *Grid
	wordsNum int
	cells    []cell
}

func NewGrid(cfg Config) *Grid {
	if cfg.Width < 2 || cfg.Height < 2 {
		panic("expected width and height to be greater then 1")
	}
	return &Grid{
		cells: make([]cell, cfg.Width*cfg.Height),
		cfg:   &cfg,
	}
}

func (g *Grid) Config() Config {
	return *g.cfg
}

func (g *Grid) PosXY(x, y int) Position {
	if x < 0 || x >= g.cfg.Width {
		panic("x is out of range")
	}
	if y < 0 || y >= g.cfg.Height {
		panic("y is out of range")
	}

	return Position{
		g:     g,
		index: y*g.cfg.Width + x,
	}
}

func (g *Grid) Pos(p Position) Position {
	if p.g.cfg != nil && p.g.cfg != g.cfg {
		panic("re-using position from different hierarchy is not allowed")
	}
	return Position{
		g:     g,
		index: p.index,
	}
}

// Returns number of added words
func (g *Grid) WriteWord(cell Position, vertical bool, word Word) (int, *Grid) {
	// TODO: there should be at least one cross with existing word
	// TODO: write direction is only forward

	// check space
	for i := 0; i < len(word); i++ {
		ok, next := cell.forward(vertical)
		if !ok {
			return 0, nil
		}
		l := next.Letter()
		if l == LetterEndOfWord {
			return 0, nil
		}
		if l != LetterEmpty && l != word[i] {
			return 0, nil
		}
	}

	// check word sticking

	// update

	return 0, nil
}

// func (g *Grid) Word(index int) Word {

// }

// func (g *Grid) WordsNum() int {

// }

func (p Position) Next() (bool, Position) {
	nextIndex := p.index + 1
	if nextIndex == len(p.g.cells) {
		return false, p
	}
	return true, Position{g: p.g, index: nextIndex}
}

func (p Position) Letter() Letter {
	return p.g.cells[p.index].letter
}

func (p Position) IsLetter() bool {
	return Has(p.g.cells[p.index].flags, cellFlagWord)
}

func (c Position) forward(vertical bool) (bool, Position) {
	if vertical {
		return c.down()
	}
	return c.right()
}

func (c Position) right() (bool, Position) {
	nextIndex := c.index + 1
	if nextIndex%c.g.cfg.Width == c.g.cfg.Width {
		return false, c
	}
	return true, Position{g: c.g, index: nextIndex}
}

func (c Position) down() (bool, Position) {
	nextIndex := c.index + c.g.cfg.Height
	if nextIndex%c.g.cfg.Height == c.g.cfg.Height {
		return false, c
	}
	return true, Position{g: c.g, index: nextIndex}
}
