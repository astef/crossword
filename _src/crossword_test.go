package crossword

import (
	"testing"
)

func TestNewCrossword(t *testing.T) {

	voc := new(Vocabulary)
	voc.Add("asdasd")
	voc.Add("qweqwe")
	voc.Add("wqfwqf")

	New(20, 20, voc, 123)
}
