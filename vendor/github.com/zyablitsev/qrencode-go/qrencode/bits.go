package qrencode

import (
	"bytes"
	"image"
	"image/color"
	"io"
)

// The test benchmark shows that encoding with boolBitVector/boolBitGrid is
// twice as fast as byteBitVector/byteBitGrid and uin32BitVector/uint32BitGrid.

type BitVector struct {
	boolBitVector
}

type BitGrid struct {
	boolBitGrid
}

func (v *BitVector) AppendBits(b BitVector) {
	v.boolBitVector.AppendBits(b.boolBitVector)
}

func NewBitGrid(width, height int) *BitGrid {
	return &BitGrid{newBoolBitGrid(width, height)}
}

func (v *BitVector) String() string {
	b := bytes.Buffer{}
	for i, l := 0, v.Length(); i < l; i++ {
		if v.Get(i) {
			b.WriteString("1")
		} else {
			b.WriteString("0")
		}
	}
	return b.String()
}

func (g *BitGrid) String() string {
	b := bytes.Buffer{}
	for y, w, h := 0, g.Width(), g.Height(); y < h; y++ {
		for x := 0; x < w; x++ {
			if g.Empty(x, y) {
				b.WriteString(" ")
			} else if g.Get(x, y) {
				b.WriteString("#")
			} else {
				b.WriteString("_")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

// Encode the Grid in ANSI escape sequences and set the background according
// to the values in the BitGrid surrounded by a white frame
func (g *BitGrid) TerminalOutput(w io.Writer) {
	white := "\033[40;37;1m"
	reset := "\033[0m"

	empty := " "
	lowhalf := "\342\226\204"
	uphalf := "\342\226\200"
	full := "\342\226\210"

	newline := "\n"

	invert := true
	if invert {
		empty, full = full, empty
		lowhalf, uphalf = uphalf, lowhalf
	}

	margin := 2
	width := g.Width() + margin*2
	for y := 0; y < margin/2; y++ {
		w.Write([]byte(white))
		for x := 0; x < width; x++ {
			w.Write([]byte(empty))
		}
		w.Write([]byte(reset))
		w.Write([]byte(newline))
	}

	for y := 0; y < g.Height(); y += 2 {
		w.Write([]byte(white))
		for x := 0; x < margin; x++ {
			w.Write([]byte(empty))
		}
		w.Write([]byte(reset))

		w.Write([]byte(white))
		for x := 0; x < g.Width(); x++ {
			if g.Get(x, y) && y+1 < g.Height() && g.Get(x, y+1) {
				w.Write([]byte(full))
			} else if g.Get(x, y) {
				w.Write([]byte(uphalf))
			} else if y+1 < g.Height() && g.Get(x, y+1) {
				w.Write([]byte(lowhalf))
			} else {
				w.Write([]byte(empty))
			}
		}
		w.Write([]byte(reset))

		w.Write([]byte(white))
		for x := 0; x < margin; x++ {
			w.Write([]byte(empty))
		}
		w.Write([]byte(reset))

		w.Write([]byte(newline))
	}

	margin = 2
	width = g.Width() + margin*2
	for y := 0; y < margin/2; y++ {
		w.Write([]byte(white))
		for x := 0; x < width; x++ {
			w.Write([]byte(empty))
		}
		w.Write([]byte(reset))
		w.Write([]byte(newline))
	}
}

// Return an image of the grid, with black blocks for true items and
// white blocks for false items, with the given block size and a
// default margin.
func (g *BitGrid) Image(blockSize int) image.Image {
	return g.ImageWithMargin(blockSize, 4)
}

// Return an image of the grid, with black blocks for true items and
// white blocks for false items, with the given block size and margin.
func (g *BitGrid) ImageWithMargin(blockSize, margin int) image.Image {
	width := blockSize * (2*margin + g.Width())
	height := blockSize * (2*margin + g.Height())
	i := image.NewGray16(image.Rect(0, 0, width, height))
	for y := 0; y < blockSize*margin; y++ {
		for x := 0; x < width; x++ {
			i.Set(x, y, color.White)
			i.Set(x, height-1-y, color.White)
		}
	}
	for y := blockSize * margin; y < height-blockSize*margin; y++ {
		for x := 0; x < blockSize*margin; x++ {
			i.Set(x, y, color.White)
			i.Set(width-1-x, y, color.White)
		}
	}
	for y, w, h := 0, g.Width(), g.Height(); y < h; y++ {
		for x := 0; x < w; x++ {
			x0 := blockSize * (x + margin)
			y0 := blockSize * (y + margin)
			c := color.White
			if g.Get(x, y) {
				c = color.Black
			}
			for dy := 0; dy < blockSize; dy++ {
				for dx := 0; dx < blockSize; dx++ {
					i.Set(x0+dx, y0+dy, c)
				}
			}
		}
	}
	return i
}
