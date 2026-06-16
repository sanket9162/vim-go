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
		if c.col > len(c.buf.Lines[c.row]) {
			c.col = len(c.buf.Lines[c.row])
		}
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
	if c.row < len(c.buf.Lines)-1 {
		c.row++
		if c.col > len(c.buf.Lines[c.row]) {
			c.col = len(c.buf.Lines[c.row])
		}
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
