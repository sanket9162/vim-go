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

func (c *Cursor) Col() int {
	return c.col
}

func (c *Cursor) Row() int {
	return c.row
}

func (c *Cursor) SetPos(col, row int) {
	c.col = col
	c.row = row
	c.idealX = col
}

func (c *Cursor) MoveLeft() {
	if c.col > 0 {
		c.col--
		c.idealX--
	}
}

func (c *Cursor) MoveRight() {
	lineLen := c.buf.LineLength(c.row)
	if c.col < lineLen {
		c.col++
		c.idealX = c.col
	}
}

func (c *Cursor) MoveUp() {
	if c.row > 0 {
		c.row--
		if c.col > len(c.buf.Lines[c.row]) {
			c.col = len(c.buf.Lines[c.row])
		}
	}
}

func (c *Cursor) MoveToStartOfLine() {
	c.col = 0
	c.idealX = 0
}

func (c *Cursor) MoveToEndOfLine() {
	lineLen := c.buf.LineLength(c.row)
	c.col = lineLen
	c.idealX = lineLen
}

func (c *Cursor) MoveDown() {
	if c.row < len(c.buf.Lines)-1 {
		c.row++
		if c.col > len(c.buf.Lines[c.row]) {
			c.col = len(c.buf.Lines[c.row])
		}
	}
}

func (c *Cursor) snapToLineLength() {
	lineLen := c.buf.LineLength(c.row)

	if c.idealX > lineLen {
		c.col = lineLen
	} else {
		c.col = c.idealX
	}
}
