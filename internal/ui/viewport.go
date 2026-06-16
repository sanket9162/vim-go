package ui

type Viewport struct {
	Width   int
	Height  int
	OffsetX int
	OffsetY int
}

func NewViewport(w, h int) *Viewport {
	return &Viewport{
		Width:  w,
		Height: h - 1,
	}
}

// SetSize updates the viewport dimensions.
func (v *Viewport) SetSize(w, h int) {
	v.Width = w
	v.Height = h
}

// ScrollTo adjusts the viewport offset to keep the cursor in view.
func (v *Viewport) ScrollTo(col, row int) {
	if row < v.OffsetY {
		v.OffsetY = row
	} else if row >= v.OffsetY+v.Height {
		v.OffsetY = row - v.Height + 1
	}

	if col < v.OffsetX {
		v.OffsetX = col
	} else if col >= v.OffsetX+v.Width {
		v.OffsetX = col - v.Width + 1
	}
}
