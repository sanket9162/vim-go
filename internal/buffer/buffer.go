package buffer

type Buffer struct {
	Lines [][]rune
}

func NewBuffer() *Buffer {
	return &Buffer{
		Lines: [][]rune{{}},
	}
}

func (b *Buffer) LineCount() int {
	return len(b.Lines)
}

func (b *Buffer) LineLength(row int) int {
	if row < 0 || row >= len(b.Lines) {
		return 0
	}
	return len(b.Lines[row])
}

func (b *Buffer) InsertChar(row, col int, r rune) {
	if row >= len(b.Lines) {
		return
	}

	line := b.Lines[row]
	line = append(line, 0)
	copy(line[col+1:], line[col:])
	line[col] = r
	b.Lines[row] = line
}

func (b *Buffer) DeleteChar(row, col int) (int, int) {
	if col > 0 {
		line := b.Lines[row]
		b.Lines[row] = append(line[:col-1], line[col:]...)
		return row, col - 1
	} else if row > 0 {
		prevLine := b.Lines[row-1]
		newCol := len(prevLine)
		b.Lines[row-1] = append(prevLine, b.Lines[row]...)

		b.Lines = append(b.Lines[:row], b.Lines[row+1:]...)
		return row - 1, newCol
	}
	return row, col
}

func (b *Buffer) InsertNewline(row, col int) {
	line := b.Lines[row]
	remaining := append([]rune{}, line[col:]...) // Copy text after cursor
	b.Lines[row] = line[:col]                    // Trim current line

	// Insert the 'remaining' slice as a new line
	newLines := append(b.Lines[:row+1], nil)
	copy(newLines[row+2:], b.Lines[row+1:])
	newLines[row+1] = remaining
	b.Lines = newLines
}
