package main

import (
	"os"
	"strings"
)

type Generator struct {
}

// ┌───┬───┬───┬───┐
// │   │ F │ U │ N │
// ├───┼───┼───┼───┤
// │ G │ U │ N │   │
// ├───┼───┼───┼───┤
// │   │ N │ O │   │
// ├───┼───┼───┼───┤
// │   │ D │   │   │
// └───┴───┴───┴───┘

func main() {
	// print(unicode.IsPrint(LetterEndOfWord))
	// print(unicode.IsPrint(LetterNull))

	grid := NewGrid(Config{Width: 16, Height: 9})

	// curPos := grid.Pos(Position{})

	for {
		// p := grid.GetPlacementPattern

		// voc.SelectWordByPattern

		// theoryGrid := grid.AddWord

		// unmatched := theoryGrid.GetUnmatchedPatterns

		// for u := range unmatched {
		//  voc.SelectWordByPattern(u)
		//
		// }

	}

	printGrid(grid)

	// curPos := grid.Pos(Position{})
	// for {
	// }
}

func printGrid(grid *Grid) {
	// clear screen
	os.Stdout.WriteString("\033c")

	cfg := grid.Config()

	size := (cfg.Width*4 + 2) * (cfg.Height*2 + 1)

	sb := strings.Builder{}
	sb.Grow(size)

	// body
	for y := 0; y < cfg.Height; y++ {
		// top line
		if y == 0 {
			sb.WriteString("┌───")
			for x := 1; x < cfg.Width; x++ {
				sb.WriteString("┬───")
			}
			sb.WriteString("┐\n")
		} else {
			sb.WriteString("├───")
			for x := 1; x < cfg.Width; x++ {
				sb.WriteString("┼───")
			}
			sb.WriteString("┤\n")
		}

		// middle line
		for x := 0; x < cfg.Width; x++ {
			pos := grid.PosXY(x, y)
			var letter rune = ' '
			if pos.IsLetter() {
				letter = rune(pos.Letter())
			}

			sb.WriteRune('│')
			sb.WriteRune(' ')    // left mark
			sb.WriteRune(letter) // letter
			sb.WriteRune(' ')    // right mark
		}
		sb.WriteString("│\n")
	}

	// bottom line
	sb.WriteString("└───")
	for x := 1; x < cfg.Width; x++ {
		sb.WriteString("┴───")
	}
	sb.WriteString("┘\n")

	// output
	os.Stdout.WriteString(sb.String())
}
