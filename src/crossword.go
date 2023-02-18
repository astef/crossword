package crossword

type Square struct {
	Letter rune

	Left   *Square
	Top    *Square
	Right  *Square
	Bottom *Square
}

type Crossword struct {
	Grid [][]Square

	Vocabulary *Vocabulary
}

func New(width, height int, voc *Vocabulary, seed int64) *Crossword {
	// init & link squares
	cw := Crossword{Grid: make([][]Square, height), Vocabulary: voc}
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
		}
	}

	// init random
	// r := rand.New(rand.NewSource(seed))

	// // main loop
	// // place
	// currentWord, currentLetter := 0

	// for {

	// }

	// entry := voc.Entries[r.Intn(len(cw.Vocabulary.Entries))]

	return &cw
}
