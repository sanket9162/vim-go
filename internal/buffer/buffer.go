package buffer

import (
	"bufio"
	"os"
)

// Buffer represents the text content being edited as a slice of rune slices.
type Buffer struct {
	Lines [][]rune
}

// NewBuffer creates and returns a new empty Buffer.
func NewBuffer() *Buffer {
	return &Buffer{
		Lines: [][]rune{{}},
	}
}

// Load reads the contents of a file into the buffer.
func (b *Buffer) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, []rune(scanner.Text()))
	}

	if len(lines) == 0 {
		lines = append(lines, []rune{})
	}

	b.Lines = lines
	return scanner.Err()
}

// Save writes the buffer contents to the specified file.
func (b *Buffer) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i, line := range b.Lines {
		_, err := writer.WriteString(string(line))
		if err != nil {
			return err
		}
		if i < len(b.Lines)-1 {
			writer.WriteByte('\n')
		}
	}
	return writer.Flush()
}

// LineCount returns the number of lines currently in the buffer.
func (b *Buffer) LineCount() int {
	return len(b.Lines)
}

// LineLength returns the number of characters in the specified row.
func (b *Buffer) LineLength(row int) int {
	if row < 0 || row >= len(b.Lines) {
		return 0
	}
	return len(b.Lines[row])
}

// InsertChar inserts a rune at the specified row and column.
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

// DeleteChar removes a character at the specified row and column.
// If the column is 0, it merges the current line with the previous one.
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

// InsertNewline splits the line at the specified row and column into two lines.
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
