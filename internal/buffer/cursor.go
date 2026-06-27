package buffer

type Cursor struct {
	col    int
	row    int
	idealX int
	buf    *Buffer
}

func NewCursor(buf *Buffer) *Cursor {
	return &Cursor{
		col:    0,
		row:    0,
		idealX: 0,
		buf:    buf,
	}
}

// Col returns the column of the cursor.
func (c *Cursor) Col() int {
	return c.col
}

// Row returns the row of the cursor.
func (c *Cursor) Row() int {
	return c.row
}

// SetPos sets the position of the cursor.
func (c *Cursor) SetPos(col, row int) {
	c.col = col
	c.row = row
	c.idealX = col
}

// MoveLeft moves the cursor to the left.
func (c *Cursor) MoveLeft() {
	if c.col > 0 {
		c.col--
		c.idealX--
	}
}

// MoveRight moves the cursor to the right.
func (c *Cursor) MoveRight() {
	lineLen := c.buf.LineLength(c.row)
	if c.col < lineLen {
		c.col++
		c.idealX = c.col
	}
}

// MoveUp moves the cursor to the up.
func (c *Cursor) MoveUp() {
	if c.row > 0 {
		c.row--
		c.snapToLineLength()
	}
}

// MoveToStartOfLine moves the cursor to the start of the line.
func (c *Cursor) MoveToStartOfLine() {
	c.col = 0
	c.idealX = 0
}

// MoveToEndOfLine moves the cursor to the end of the line.
func (c *Cursor) MoveToEndOfLine() {
	lineLen := c.buf.LineLength(c.row)
	c.col = lineLen
	c.idealX = lineLen
}

// MoveDown moves the cursor to the down.
func (c *Cursor) MoveDown() {
	if c.row < c.buf.LineCount()-1 {
		c.row++
		c.snapToLineLength()
	}
}

// snapToLineLength snaps the cursor to the end of the line.
func (c *Cursor) snapToLineLength() {
	lineLen := c.buf.LineLength(c.row)

	if c.idealX > lineLen {
		c.col = lineLen
	} else {
		c.col = c.idealX
	}
}

// MoveToStartOfFile jumps to the very first line
func (c *Cursor) MoveToStartOfFile() {
	c.SetPos(0, 0)
}

// MoveToEndOfFile jumps to the very last line
func (c *Cursor) MoveToEndOfFile() {
	lastRow := c.buf.LineCount() - 1
	if lastRow < 0 {
		lastRow = 0
	}
	c.SetPos(0, lastRow)
}

// MoveToNextWord moves to the start of the next wrod
func (c *Cursor) MoveToNextWord() {
	line := c.buf.GetLine(c.row)
	col := c.col

	// Move past current alphanumeric characters
	for col < len(line) && isWordChar(line[col]) {
		col++
	}
	// Move past whitespace/punctuation
	for col < len(line) && !isWordChar(line[col]) {
		col++
	}

	if col < len(line) {
		c.SetPos(col, c.row)
	} else if c.row < c.buf.LineCount()-1 {
		c.SetPos(0, c.row+1)
	}
}

func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' || r == '.' || r == '!' || r == ':' || r == '+'
}

// MoveToBackwardWord moves to the cursor backward to the start of the current or previous word.
func (c *Cursor) MoveToBackwardWord() {
	col := c.col
	row := c.row

	// If we're already at the start of the file, do nothing
	if row == 0 && col == 0 {
		return
	}

	getLineRunes := func(r int) []rune {
		if r < 0 || r >= c.buf.LineCount() {
			return nil
		}
		return c.buf.GetLine(r)
	}

	line := getLineRunes(row)

	// If at start of line, wrap to end of previous line
	if col == 0 {
		row--
		line = getLineRunes(row)
		col = len(line)
	}

	// 1. Move backward past whitespace/punctuation to find a word character
	for col > 0 && !isWordChar(line[col-1]) {
		col--
	}

	// 2. Move backward past aplhanumeric word characters to find start of the word
	for col > 0 && isWordChar(line[col-1]) {
		col--
	}

	c.SetPos(col, row)

}

// MoveToEndOfWord move the cursor forward to end of the current or next wrod.
func (c *Cursor) MoveToEndOfWord() {
	row := c.row
	col := c.col
	line := c.buf.GetLine(row)

	wrapToNextLine := func() bool {
		if row < c.buf.LineCount()-1 {
			row++
			line = c.buf.GetLine(row)
			col = 0
			return true
		}
		return false
	}

	if len(line) == 0 || col >= len(line)-1 {
		if !wrapToNextLine() {
			return
		}
	} else {
		// Increment by 1 first to prevent cursor from getting stuck on current end-of-word
		col++
	}

	// 1.Move forward past whitespace/punctuation of find a wrod character
	for col < len(line) && !isWordChar(line[col]) {
		col++
		if col >= len(line) {
			if !wrapToNextLine() {
				// put cursor at the last character of the previous line on EOF boundary
				row--
				line = c.buf.GetLine(row)
				c.SetPos(len(line)-1, row)
				return
			}
		}
	}
	//

}
