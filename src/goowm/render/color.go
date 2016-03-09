package render

type Color struct {
	R uint32
	G uint32
	B uint32
	A uint32
}

func (c Color) RGBA() (r, g, b, a uint32) {
	return c.R, c.G, c.B, c.A
}
