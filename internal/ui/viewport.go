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
		Height: h,
	}
}

func (v *Viewport) SetSize(w, h int) {
	v.Width = w
	v.Height = h
}
