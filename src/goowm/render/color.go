package render

type Color struct {
	R uint32
	G uint32
	B uint32
	A uint32
}

func NewColor(r, g, b uint32) Color {
	return Color{r, g, b, 255}
}

func (c Color) ToFloat64s() (r, g, b float64) {
	r = float64(c.R) / 255.0
	g = float64(c.G) / 255.0
	b = float64(c.B) / 255.0

	return
}

func (c Color) RGBA() (r, g, b, a uint32) {
	return c.R, c.G, c.B, c.A
}
