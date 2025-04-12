package crossword

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestVocabularyAdd(t *testing.T) {
	voc := new(Vocabulary)
	voc.Add("apple")
	voc.Add("banana")
	voc.Add("orange")

	if string(voc.Entries[0].Word) != "apple" {
		t.Fail()
	}
	if string(voc.Entries[1].Word) != "banana" {
		t.Fail()
	}
	if string(voc.Entries[2].Word) != "orange" {
		t.Fail()
	}

}

type wordPart struct {
	word      string
	partIndex int
}

func sortWordParts(w []wordPart) {
	sort.Slice(w, func(a, b int) bool {
		return w[a].word < w[b].word || (w[a].word == w[b].word && w[a].partIndex < w[b].partIndex)
	})
}

func (w wordPart) String() string {
	return fmt.Sprintf("{%v, %v}", w.word, w.partIndex)
}

func TestVocabularySearchByPattern(t *testing.T) {
	allWords := []string{
		"apple",
		"art",
		"avocado",
		"banana",
		"copper",
		"orange",
		"tuple",
		"zoo",
	}

	var expectedAllWords []wordPart
	for _, w := range allWords {
		expectedAllWords = append(expectedAllWords, wordPart{w, len(w)})
	}

	voc := new(Vocabulary)
	for _, word := range allWords {
		voc.Add(word)
	}

	tests := map[string]struct {
		pattern       string
		requiredIndex int
		expectedWords []wordPart
	}{
		"all": {
			pattern:       strings.Repeat("_", voc.MaxWordLength),
			requiredIndex: voc.MaxWordLength,
			expectedWords: expectedAllWords,
		},
		"no match": {
			pattern:       "___xyz___",
			requiredIndex: 3,
			expectedWords: []wordPart{},
		},
		"no match because of length": {
			pattern:       "appl",
			requiredIndex: 0,
			expectedWords: []wordPart{},
		},
		"no match because of optional": {
			pattern:       "app_z",
			requiredIndex: 0,
			expectedWords: []wordPart{},
		},
		"no match because of optional sticking at the end": {
			pattern:       "ap___xyz",
			requiredIndex: 0,
			expectedWords: []wordPart{},
		},
		"no match because of optional sticking at the beginning": {
			pattern:       "xyz_pp__",
			requiredIndex: 4,
			expectedWords: []wordPart{},
		},
		"by letter count": {
			pattern:       "___",
			requiredIndex: 3,
			expectedWords: []wordPart{{"art", 3}, {"zoo", 3}},
		},
		"by exact match": {
			pattern:       "art",
			requiredIndex: 0,
			expectedWords: []wordPart{{"art", 0}},
		},
		"by first letter": {
			pattern:       "a____",
			requiredIndex: 0,
			expectedWords: []wordPart{{"apple", 0}, {"art", 0}},
		},
		"by first letter with optionals": {
			pattern:       "a__l_",
			requiredIndex: 0,
			expectedWords: []wordPart{{"apple", 0}},
		},
		"by first letter with skipped optionals at the end": {
			pattern:       "a____xyz",
			requiredIndex: 0,
			expectedWords: []wordPart{{"art", 0}},
		},
		"by first letter with skipped optionals at the beginning": {
			pattern:       "xyz_a____",
			requiredIndex: 4,
			expectedWords: []wordPart{{"apple", 0}, {"art", 0}},
		},
		"by first letter with skipped optionals at both sides": {
			pattern:       "xyz_a_____xyz",
			requiredIndex: 4,
			expectedWords: []wordPart{{"apple", 0}, {"art", 0}},
		},
		"by first letter with skipped optionals at both sides and matched optional": {
			pattern:       "xyz_a__l__xyz",
			requiredIndex: 4,
			expectedWords: []wordPart{{"apple", 0}},
		},
		"by first 2 letters": {
			pattern:       "ap___",
			requiredIndex: 0,
			expectedWords: []wordPart{{"apple", 0}},
		},
		"by first 3 letters": {
			pattern:       "app__",
			requiredIndex: 0,
			expectedWords: []wordPart{{"apple", 0}},
		},
		"by middle letter": {
			pattern:       "__p__",
			requiredIndex: 2,
			expectedWords: []wordPart{{"apple", 2}, {"tuple", 2}},
		},
		"by middle letter repeated": {
			pattern:       "___p___",
			requiredIndex: 3,
			expectedWords: []wordPart{{"apple", 1}, {"apple", 2}, {"tuple", 2}, {"copper", 2}, {"copper", 3}},
		},
		"by middle letter with optionals": {
			pattern:       "__p_e",
			requiredIndex: 2,
			expectedWords: []wordPart{{"apple", 2}, {"tuple", 2}},
		},
		"by middle 2 letters": {
			pattern:       "_pp__",
			requiredIndex: 1,
			expectedWords: []wordPart{{"apple", 1}},
		},
		"by middle 3 letters": {
			pattern:       "_ppl_",
			requiredIndex: 1,
			expectedWords: []wordPart{{"apple", 1}},
		},
		"by last letter": {
			pattern:       "____e",
			requiredIndex: 4,
			expectedWords: []wordPart{{"apple", 4}, {"tuple", 4}},
		},
		"by last letter with optionals": {
			pattern:       "__p_e",
			requiredIndex: 4,
			expectedWords: []wordPart{{"apple", 4}, {"tuple", 4}},
		},
		"by last 2 letters": {
			pattern:       "___le",
			requiredIndex: 3,
			expectedWords: []wordPart{{"apple", 3}, {"tuple", 3}},
		},
		"by last 2 letters with optionals": {
			pattern:       "t__le",
			requiredIndex: 3,
			expectedWords: []wordPart{{"tuple", 3}},
		},
		"by last 3 letters": {
			pattern:       "__ple",
			requiredIndex: 2,
			expectedWords: []wordPart{{"apple", 2}, {"tuple", 2}},
		},
		"by last 3 letters with optionals": {
			pattern:       "a_ple",
			requiredIndex: 2,
			expectedWords: []wordPart{{"apple", 2}},
		},
		"by repeated required part": {
			pattern:       "___an___",
			requiredIndex: 3,
			expectedWords: []wordPart{{"banana", 1}, {"banana", 3}, {"orange", 2}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			queryByPattern := voc.newQueryByPattern(
				Pattern{
					Sequence:          []rune(strings.ReplaceAll(tc.pattern, "_", string(Empty))),
					RequiredPartIndex: tc.requiredIndex,
				})

			actualWords := []wordPart{}
			for ep := queryByPattern(); ep.Entry != nil; ep = queryByPattern() {
				actualWords = append(actualWords, wordPart{string(ep.Entry.Word), ep.PartIndex})
			}

			sortWordParts(actualWords)
			sortWordParts(tc.expectedWords)

			if !reflect.DeepEqual(actualWords, tc.expectedWords) {
				t.Fatalf("expected: %v, got: %v", tc.expectedWords, actualWords)
			}

		})
	}

}
