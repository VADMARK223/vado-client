package common

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type Indicator struct {
	widget.BaseWidget
	Circle        *canvas.Circle
	indicatorSize fyne.Size
}

func NewIndicator(fill color.Color, size fyne.Size) *Indicator {
	circle := canvas.NewCircle(fill)
	circle.StrokeColor = color.Gray{Y: 0x99}
	circle.StrokeWidth = 1

	ind := &Indicator{
		Circle:        circle,
		indicatorSize: size,
	}
	ind.ExtendBaseWidget(ind)
	return ind
}

func (i *Indicator) MinSize() fyne.Size {
	return i.indicatorSize
}

func (i *Indicator) CreateRenderer() fyne.WidgetRenderer {
	return &indicatorRenderer{i}
}

type indicatorRenderer struct {
	indicator *Indicator
}

func (r *indicatorRenderer) Layout(_ fyne.Size) {
	r.indicator.Circle.Resize(r.indicator.indicatorSize)
}

func (r *indicatorRenderer) MinSize() fyne.Size {
	return r.indicator.indicatorSize
}

func (r *indicatorRenderer) Refresh() {
	r.indicator.Circle.Refresh()
}

func (r *indicatorRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *indicatorRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.indicator.Circle}
}

func (r *indicatorRenderer) Destroy() {}

// SetFillColor меняет цвет индикатора и перерисовывает его
func (i *Indicator) SetFillColor(c color.Color) {
	i.Circle.FillColor = c
	i.Refresh()
}
