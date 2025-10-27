package color

import "image/color"

var (
	red    = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	green  = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	orange = color.RGBA{R: 255, G: 165, B: 0, A: 255}
)

func Red() color.Color    { return red }
func Green() color.Color  { return green }
func Orange() color.Color { return orange }
