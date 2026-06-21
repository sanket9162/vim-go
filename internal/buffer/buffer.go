package buffer

import (
	"os"
)

// Buffer wraps a GapBuffer and translates between 2D (row, col) coordinates
// used by the UI and the 1D logical indices used by the GapBuffer.
type Buffer struct {
	gb        *GapBuffer
	lineStart []int // Stores the 1D index where each line begins
}

// Load reads the contents of a file into the buffer.
func (b *Buffer) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	
	// Create a new GapBuffer large enough to hold the file + some extra space
	size := len(data) * 2
	if size < 1024 {
		size = 1024
	}
	b.gb = NewGapBuffer(size)
	
	// Insert all data at once
	for i, r := range string(data) {
		b.gb.data[i] = r
	}
	b.gb.gapLeft = len(string(data))
	
	b.recomputeLineStarts()
	return nil
}

// Save writes the buffer contents to the specified file.
func (b *Buffer) Save(filename string) error {
	text := b.gb.Text()
	return os.WriteFile(filename, []byte(text), 0644)
}

// NewBuffer creates a new Buffer utilizing a GapBuffer backend.
func NewBuffer() *Buffer {
	b := &Buffer{
		gb:        NewGapBuffer(1024),
		lineStart: []int{0}, // Line 0 starts at index 0
	}
	return b
}

// LineCount returns the number of lines in the buffer.
func (b *Buffer) LineCount() int {
	return len(b.lineStart)
}

// LineLength returns the number of characters in the specified row, excluding the newline.
func (b *Buffer) LineLength(row int) int {
	if row < 0 || row >= b.LineCount() {
		return 0
	}
	
	start := b.lineStart[row]
	var end int
	if row == b.LineCount()-1 {
		end = b.gb.Length()
	} else {
		end = b.lineStart[row+1] - 1 // -1 to exclude the '\n'
	}
	return end - start
}

// getIndex converts a 2D (row, col) position into a 1D GapBuffer index.
func (b *Buffer) getIndex(row, col int) int {
	if row < 0 {
		return 0
	}
	if row >= b.LineCount() {
		return b.gb.Length()
	}
	
	start := b.lineStart[row]
	length := b.LineLength(row)
	if col > length {
		col = length
	}
	return start + col
}

// recomputeLineStarts scans the buffer and updates the lineStart cache.
// This is a naive implementation; optimized versions only scan from the edit point.
func (b *Buffer) recomputeLineStarts() {
	starts := []int{0}
	length := b.gb.Length()
	for i := 0; i < length; i++ {
		if b.gb.GetRune(i) == '\n' {
			starts = append(starts, i+1)
		}
	}
	b.lineStart = starts
}

// InsertChar inserts a rune at the specified row and column.
func (b *Buffer) InsertChar(row, col int, r rune) {
	idx := b.getIndex(row, col)
	b.gb.Insert(idx, r)
	b.recomputeLineStarts()
}

// InsertNewline inserts a newline character at the cursor position.
func (b *Buffer) InsertNewline(row, col int) {
	idx := b.getIndex(row, col)
	b.gb.Insert(idx, '\n')
	b.recomputeLineStarts()
}

// DeleteChar removes a character before the specified row and column.
func (b *Buffer) DeleteChar(row, col int) (int, int) {
	idx := b.getIndex(row, col)
	if idx > 0 {
		b.gb.Delete(idx)
		b.recomputeLineStarts()
		
		// Calculate new 2D position to return to the cursor
		if col > 0 {
			return row, col - 1
		}
		// If we deleted a newline, cursor moves to end of previous line
		return row - 1, b.LineLength(row - 1)
	}
	return row, col
}

// DeleteRange deletes all characters between (startRow, startCol) and (endRow, endCol) inclusive.
func (b *Buffer) DeleteRange(startRow, startCol, endRow, endCol int) (int, int) {
	if startRow > endRow || (startRow == endRow && startCol > endCol) {
		startRow, endRow = endRow, startRow
		startCol, endCol = endCol, startCol
	}

	startIdx := b.getIndex(startRow, startCol)
	endIdx := b.getIndex(endRow, endCol) + 1
	if endIdx > b.gb.Length() {
		endIdx = b.gb.Length()
	}

	amount := endIdx - startIdx
	if amount <= 0 {
		return startRow, startCol
	}

	for i := 0; i < amount; i++ {
		b.gb.Delete(startIdx + 1)
	}

	b.recomputeLineStarts()
	return startRow, startCol
}

// GetLine returns a specific line as a slice of runes for rendering.
func (b *Buffer) GetLine(row int) []rune {
	if row < 0 || row >= b.LineCount() {
		return nil
	}
	
	start := b.lineStart[row]
	length := b.LineLength(row)
	line := make([]rune, length)
	
	for i := 0; i < length; i++ {
		line[i] = b.gb.GetRune(start + i)
	}
	return line
}
