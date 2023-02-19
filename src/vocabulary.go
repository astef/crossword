package crossword

const (
	Empty rune = 0
	End   rune = 3
)

type Vocabulary struct {
	MaxWordLength int
	Entries       []*Entry
	indexes       []index
}

type Entry struct {
	Word []rune
}

type index = map[string][]EntryPart

type EntryPart struct {
	Entry     *Entry
	PartIndex int
}

type Pattern struct {
	Sequence          Segment
	RequiredPartIndex int
}

func (voc *Vocabulary) Add(word string) {
	wordRunes := []rune(word)

	// single char words are ignored
	if len(wordRunes) < 2 {
		return
	}

	voc.MaxWordLength = maxInt(voc.MaxWordLength, len(wordRunes))
	voc.Entries = append(voc.Entries, &Entry{Word: wordRunes})
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func lettersMatch(letters []rune, pattern []rune, reverse bool) bool {
	if len(letters) > len(pattern) {
		return false
	}

	for i := 0; i < len(letters); i++ {
		var letterIndex, patternIndex int
		if reverse {
			letterIndex, patternIndex = len(letters)-1-i, len(pattern)-1-i
		} else {
			letterIndex, patternIndex = i, i
		}

		var patternRune, lettersRune = pattern[patternIndex], letters[letterIndex]

		if patternRune != Empty && patternRune != lettersRune {
			return false
		}

		// make sure words are not "sticking", e.g. "APPLE" should not match "F_PPLE" pattern
		if i == len(letters)-1 && len(pattern) > len(letters) {
			var nextPatternIndex int
			if reverse {
				nextPatternIndex = patternIndex - 1
			} else {
				nextPatternIndex = patternIndex + 1
			}
			if pattern[nextPatternIndex] != Empty {
				return false
			}
		}
	}

	return true
}

func newIndex(entries []*Entry, keyLength int) index {
	index := make(index)

	for _, entry := range entries {

		for i := 0; i < len(entry.Word)-keyLength+1; i++ {
			key := entry.Word[i : i+keyLength]
			index[string(key)] = append(index[string(key)], EntryPart{Entry: entry, PartIndex: i})

		}

	}

	return index
}

func (voc *Vocabulary) newQueryBySubstring(requiredSubstring []rune) func() EntryPart {
	if len(requiredSubstring) == 0 {
		// all entries match
		i := 0
		return func() EntryPart {
			if i == len(voc.Entries) {
				return EntryPart{}
			}
			entry := voc.Entries[i]
			i++
			return EntryPart{Entry: entry, PartIndex: len(entry.Word)}
		}
	}

	// lazy index initialization
	for l := len(voc.indexes); l < len(requiredSubstring); l++ {
		voc.indexes = append(voc.indexes, newIndex(voc.Entries, l+1))
	}

	i := 0
	indexEntries := voc.indexes[len(requiredSubstring)-1][string(requiredSubstring)]

	return func() EntryPart {
		if len(indexEntries) == i {
			return EntryPart{}
		}
		result := indexEntries[i]
		i++
		return result
	}
}

/*
Creates a generator, which will return all words, matching the pattern.

Returning nil indicates the end of results.

Pattern may consist of Empty rune and letters.
When multiple letters come in a row, we call this sequence a word part.

In a pattern, there may be a required word part.
It is required to be present in a target word in order to match the pattern.
requiredIndex specifies the index of the required word part in the pattern.

Special case:

	requiredIndex == len(pattern) // means no required part

Empty rune in the pattern means "any character".

All other letters in pattern are optional parts.
They are only required to match if the word is "crossing" them.
Optional parts should not "stick" to the word from any side
(there should be an Empty rune between them).

In order to match, the word must conform to the length of the pattern
(there should be enough empty space), all crossed optional parts, the required part.
*/
func (voc *Vocabulary) newQueryByPattern(pattern Pattern) func() EntryPart {
	if len(pattern.Sequence) < 2 {
		panic("pattern is too short")
	}

	if pattern.RequiredPartIndex < 0 || pattern.RequiredPartIndex > len(pattern.Sequence) {
		panic("requiredIndex is out of bounds")
	}

	if pattern.RequiredPartIndex > 0 && pattern.RequiredPartIndex < len(pattern.Sequence) && pattern.Sequence[pattern.RequiredPartIndex-1] != Empty {
		panic("requiredIndex is not pointing to the beginning of the required part")
	}

	leftPattern := pattern.Sequence[0:pattern.RequiredPartIndex]
	requiredSubstring := []rune{}
	requiredSubstringFinalized := false
	rightPattern := []rune{}
	for i := pattern.RequiredPartIndex; i < len(pattern.Sequence); i++ {
		r := pattern.Sequence[i]
		if r == Empty {
			requiredSubstringFinalized = true
		}
		if requiredSubstringFinalized {
			rightPattern = append(rightPattern, r)
		} else {
			requiredSubstring = append(requiredSubstring, r)
		}
	}

	substringQuery := voc.newQueryBySubstring(requiredSubstring)
	return func() EntryPart {
		for {
			indexHit := substringQuery()

			if indexHit.Entry == nil {
				return indexHit
			}

			leftPart := indexHit.Entry.Word[0:indexHit.PartIndex]
			rightPart := indexHit.Entry.Word[indexHit.PartIndex+len(requiredSubstring):]

			if lettersMatch(leftPart, leftPattern, true) &&
				lettersMatch(rightPart, rightPattern, false) {
				return indexHit
			}
		}
	}
}
