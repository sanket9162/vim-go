package buffer

// GapBuffer implements an efficient 1D text buffer using a movable memory gap.
type GapBuffer struct {
	data     []rune
	gapLeft  int
	gapRight int
}

// NewGapBuffer creates a new buffer with an initial gap size.
func NewGapBuffer(initialSize int) *GapBuffer {
	if initialSize < 16 {
		initialSize = 1024
	}
	return &GapBuffer{
		data:     make([]rune, initialSize),
		gapLeft:  0,
		gapRight: initialSize,
	}
}

// Length returns the total number of actual characters in the buffer.
func (gb *GapBuffer) Length() int {
	return len(gb.data) - (gb.gapRight - gb.gapLeft)
}

// moveGap shifts the gap to the target index.
func (gb *GapBuffer) moveGap(targetIndex int) {
	if targetIndex == gb.gapLeft {
		return
	}

	if targetIndex < gb.gapLeft {
		amount := gb.gapLeft - targetIndex
		copy(gb.data[gb.gapRight-amount:gb.gapRight], gb.data[targetIndex:gb.gapLeft])
		gb.gapLeft -= amount
		gb.gapRight -= amount
	} else {
		amount := targetIndex - gb.gapLeft
		copy(gb.data[gb.gapLeft:gb.gapLeft+amount], gb.data[gb.gapRight:gb.gapRight+amount])
		gb.gapLeft += amount
		gb.gapRight += amount
	}
}

// expandGap grows the backing array when the gap is full.
func (gb *GapBuffer) expandGap() {
	newSize := len(gb.data) * 2
	newData := make([]rune, newSize)
	newGapRight := newSize - (len(gb.data) - gb.gapRight)

	copy(newData[:gb.gapLeft], gb.data[:gb.gapLeft])
	copy(newData[newGapRight:], gb.data[gb.gapRight:])

	gb.data = newData
	gb.gapRight = newGapRight
}

// Insert adds a character at the specified index.
func (gb *GapBuffer) Insert(index int, r rune) {
	gb.moveGap(index)
	if gb.gapLeft == gb.gapRight {
		gb.expandGap()
	}
	gb.data[gb.gapLeft] = r
	gb.gapLeft++
}

// Delete removes the character immediately before the specified index.
func (gb *GapBuffer) Delete(index int) {
	gb.moveGap(index)
	if gb.gapLeft > 0 {
		gb.gapLeft--
	}
}

// GetRune returns the character at the specified logical index.
func (gb *GapBuffer) GetRune(index int) rune {
	if index < gb.gapLeft {
		return gb.data[index]
	}
	return gb.data[index+(gb.gapRight-gb.gapLeft)]
}

// Text returns the entire buffer as a string.
func (gb *GapBuffer) Text() string {
	runes := make([]rune, gb.Length())
	copy(runes[:gb.gapLeft], gb.data[:gb.gapLeft])
	copy(runes[gb.gapLeft:], gb.data[gb.gapRight:])
	return string(runes)
}
