package crossword

import (
	"errors"
	"math/rand"
)

type Segment struct {
	Start    *Square
	Length   int
	Vertical bool
}

type Square struct {
	Left   *Square
	Top    *Square
	Right  *Square
	Bottom *Square

	// not nil
	RowStart    *Square
	RowEnd      *Square
	ColumnStart *Square
	ColumnEnd   *Square

	// updatable
	Letter rune
	WordH  bool
	WordV  bool
}

type Crossword struct {
	Width int

	Height int

	Grid [][]Square

	Vocabulary *Vocabulary

	Version int

	random *rand.Rand
}

type crosswordVersionRef struct {
	crossword *Crossword
	version   int
}

func (ref crosswordVersionRef) isValid() bool {
	return ref.crossword != nil && ref.version == ref.crossword.Version
}

type Editable interface {
}

type WordPlacement struct {
	Square   *Square
	Vertical bool
	Word     *Entry
}

type WordProposal struct {
	crosswordRef crosswordVersionRef
	parent       *WordProposal
	children     []WordProposal
	Placement    WordPlacement
	Complete     bool
	Score        int
}

const (
	Ok = iota
	InvalidPlace
)

func New(width, height int, voc *Vocabulary, seed int64) *Crossword {
	if width < 2 || height < 2 {
		return nil
	}

	cw := Crossword{
		Width:      width,
		Height:     height,
		Grid:       make([][]Square, height),
		Vocabulary: voc,
		random:     rand.New(rand.NewSource(seed)),
	}
	// init left/right/top/bottom links
	for y := range cw.Grid {
		row := make([]Square, width)
		cw.Grid[y] = row
		for x, current := range row {
			if x != 0 {
				left := &row[x-1]
				current.Left = left
				left.Right = &current
			}
			if y != 0 {
				top := &cw.Grid[y-1][x]
				current.Top = top
				top.Bottom = &current
			}

			// init row/column links
			current.ColumnStart = &cw.Grid[0][x]
			current.ColumnEnd = &cw.Grid[height-1][x]
			current.RowStart = &row[0]
			current.RowEnd = &row[width-1]
		}
	}

	return &cw
}

var (
	errArgNil                    = errors.New("required argument is nil")
	errInvalidCrosswordReference = errors.New("proposal is for the different crossword")
	errOldProposalVersion        = errors.New("proposal is outdated")
	errIncompleteProposal        = errors.New("proposal is incomplete")
)

func (from *Square) RayIterator(vertical bool, backward bool) func() *Square {
	if vertical && backward {
		return func() *Square {
			if from == nil {
				return nil
			}
			result := from.Top
			from = result
			return result
		}
	} else if vertical && !backward {
		return func() *Square {
			if from == nil {
				return nil
			}
			result := from.Bottom
			from = result
			return result
		}
	} else if !vertical && backward {
		return func() *Square {
			if from == nil {
				return nil
			}
			result := from.Left
			from = result
			return result
		}
	} else {
		return func() *Square {
			if from == nil {
				return nil
			}
			result := from.Right
			from = result
			return result
		}
	}
}

func (sq *Square) LineIterator(vertical bool) func() *Square {
	if sq == nil {
		panic("receiver is nil")
	}
	if vertical {
		return sq.ColumnStart.RayIterator(vertical, false)
	} else {
		return sq.RowStart.RayIterator(vertical, false)
	}
}

func (sq *Square) PatternSequenceIterator(vertical bool) func() Segment {
	it := sq.LineIterator(vertical)

	var start *Square
	length := 0

	return func() Segment {
		for it != nil {
			if sq := it(); sq == nil {
				it = nil
				break
			}

			if start == nil {
				if sq.Letter != End {
					start = sq
					length = 1
				}
			} else {
				if sq.Letter == End {
					break
				} else {
					length++
				}

			}
		}

		result := Segment{Start: start, Length: length, Vertical: vertical}
		start = nil
		length = 0
		return result
	}
}

func (sq *Square) PatternIterator(vertical bool) func() Pattern {
	it := sq.PatternSequenceIterator(vertical)
	return func() Pattern {

	}
}

func (cw *Crossword) GetAvailablePatterns(lineIndex int, vertical bool) func() Pattern {
	// returns all possible patterns, but only with required parts

	length := cw.Width
	x, y := 0, lineIndex
	i := &x
	if vertical {
		length = cw.Height
		x, y = y, x
		i = &y
	}

	var currentSequence []rune
	var currentRequiredPartStarted bool
	currentRequiredPartIndex := -1

	return func() Pattern {

		for *i < length {
			sq := cw.Grid[y][x]
			*i++

			if sq.Letter == End {
				// TODO "scroll forwand" until the end of End blocks

				break
			} else if sq.Letter == Empty {
				currentSequence = append(currentSequence, sq.Letter)
				if currentRequiredPartStarted {
					// "look forward" to send pattern with this current required part

					currentRequiredPartStarted = false
				}
			} else {
				if currentRequiredPartStarted {
					currentSequence = append(currentSequence, sq.Letter)
				} else {
					currentRequiredPartStarted = true
					currentRequiredPartIndex = *i - 1
				}

			}
		}

		result := Pattern{
			Sequence:          currentSequence,
			RequiredPartIndex: currentRequiredPartIndex,
		}
		currentSequence = nil
		currentRequiredPartIndex = -1
		currentRequiredPartStarted = false
		return result
	}
}

func NewWordProposal(cw *Crossword, parent *WordProposal, placement WordPlacement) (WordProposal, error) {
	if cw == nil {
		return WordProposal{}, errArgNil
	}
	if parent != nil {
		if cw != parent.crosswordRef.crossword {
			return WordProposal{}, errInvalidCrosswordReference
		}
		if !parent.crosswordRef.isValid() {
			return WordProposal{}, errOldProposalVersion
		}
	}

	// check placement

	// compute score

	// compute completeness

	// update parent

}

func (cw *Crossword) AcceptWordProposal(p *WordProposal) error {
	if cw != p.crosswordRef.crossword {
		return errInvalidCrosswordReference
	}
	if !p.crosswordRef.isValid() {
		return errOldProposalVersion
	}
	if !p.Complete {
		return errIncompleteProposal
	}

	// apply all changes from the hierarchy and increment crossword version

	return nil
}
