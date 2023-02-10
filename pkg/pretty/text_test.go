package pretty

import (
	"fmt"
	"testing"
)

func TestAlignText(t *testing.T) {
	var s string

	s = AlignText("John Snow", 12, TextAlignLeft)
	fmt.Printf("%q\n", s)

	s = AlignText("John Snow", 12, TextAlignCenter)
	fmt.Printf("%q\n", s)

	s = AlignText("John Snow", 12, TextAlignJustify)
	fmt.Printf("%q\n", s)

	s = AlignText("John Snow", 12, TextAlignRight)
	fmt.Printf("%q\n", s)

	s = AlignText("John Doe", 12, TextAlignLeft)
	fmt.Printf("%q\n", s)

	s = AlignText("John Doe", 12, TextAlignCenter)
	fmt.Printf("%q\n", s)

	s = AlignText("John Doe", 12, TextAlignJustify)
	fmt.Printf("%q\n", s)

	s = AlignText("John Doe", 12, TextAlignRight)
	fmt.Printf("%q\n", s)
}
