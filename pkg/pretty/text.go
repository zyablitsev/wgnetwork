package pretty

import (
	"strings"
	"unicode/utf8"
)

// TextAlign enumerations
const (
	TextAlignLeft    int = iota // "left        "
	TextAlignCenter             // "   center   "
	TextAlignJustify            // "justify   it"
	TextAlignRight              // "       right"
)

// AlignText aligns the text as directed. For ex.:
//  * AlignText("Jon Snow", 12,    TextAlignLeft) returns "Jon Snow    "
//  * AlignText("Jon Snow", 12,  TextAlignCenter) returns "  Jon Snow  "
//  * AlignText("Jon Snow", 12, TextAlignJustify) returns "Jon     Snow"
//  * AlignText("Jon Snow", 12, TextAlignRight) returns "    Jon Snow"
func AlignText(text string, size, align int) string {
	text = strings.Trim(text, " ")
	if len(text) == 0 {
		return strings.Repeat(" ", size)
	}

	sLen := utf8.RuneCountInString(text)

	if sLen >= size {
		return text
	}

	if align == TextAlignRight {
		return strings.Repeat(" ", size-sLen) + text
	} else if align == TextAlignCenter {
		paddingLeft := ""
		paddingRight := ""
		if sLen < size {
			padding := strings.Repeat(" ", (size-sLen)/2)
			paddingLeft = padding
			paddingRight = padding
			if (size-sLen)%2 > 0 {
				paddingRight += " "
			}
		}
		return paddingLeft + text + paddingRight
	} else if align == TextAlignJustify {
		return justifyText(text, sLen, size)
	}

	return text + strings.Repeat(" ", size-sLen)
}

func justifyText(text string, sLen int, size int) string {
	// split the text into individual words
	a := strings.Split(text, " ")
	words := make([]string, 0, len(a))
	for _, s := range a {
		if s == "" {
			continue
		}

		words = append(words, s)
	}

	if len(words) == 0 {
		return strings.Repeat(" ", size)
	}

	// get the number of spaces to insert into the text
	numSpacesNeeded := size - sLen + strings.Count(text, " ")
	numSpacesNeededBetweenWords := 0
	if len(words) > 1 {
		numSpacesNeededBetweenWords = numSpacesNeeded / (len(words) - 1)
	}

	// create the output string word by word with spaces in between
	var w strings.Builder
	w.Grow(size)
	for idx, word := range words {
		if idx > 0 {
			// insert spaces only after the first word
			if idx == len(words)-1 {
				// insert all the remaining space before the last word
				w.WriteString(strings.Repeat(" ", numSpacesNeeded))
				numSpacesNeeded = 0
			} else {
				// insert the determined number of spaces between each word
				w.WriteString(strings.Repeat(" ", numSpacesNeededBetweenWords))
				// and reduce the number of spaces needed after this
				numSpacesNeeded -= numSpacesNeededBetweenWords
			}
		}

		w.WriteString(word)

		if idx == len(words)-1 && numSpacesNeeded > 0 {
			w.WriteString(strings.Repeat(" ", numSpacesNeeded))
		}
	}

	return w.String()
}
