package common

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type VadoIconButton struct {
	widget.BaseWidget
	icon     fyne.Resource
	onTapped func()

	disabled bool
	hovered  bool

	bg   *canvas.Circle
	img  *canvas.Image
	root *fyne.Container
}

func NewVadoIconButton(icon fyne.Resource, tapped func()) *VadoIconButton {
	b := &VadoIconButton{icon: icon, onTapped: tapped}
	b.ExtendBaseWidget(b)
	return b
}

func (b *VadoIconButton) CreateRenderer() fyne.WidgetRenderer {
	// фон
	b.bg = canvas.NewCircle(color.NRGBA{}) // прозрачный по умолчанию
	b.bg.StrokeWidth = 0

	// иконка
	b.img = canvas.NewImageFromResource(b.icon)
	b.img.FillMode = canvas.ImageFillContain
	b.img.Translucency = 0.0

	// кладём иконку поверх круга
	b.root = container.NewStack(b.bg, b.img)

	b.updateVisualState()
	return widget.NewSimpleRenderer(b.root)
}

func (b *VadoIconButton) MinSize() fyne.Size {
	return fyne.NewSize(24, 24)
}

func (b *VadoIconButton) Tapped(*fyne.PointEvent) {
	if b.disabled {
		return
	}
	if b.onTapped != nil {
		b.onTapped()
	}
}

func (b *VadoIconButton) Disable() {
	b.disabled = true
	b.hovered = false
	b.updateVisualState()
	b.Refresh()
}

func (b *VadoIconButton) Enable() {
	b.disabled = false
	b.updateVisualState()
	b.Refresh()
}

func (b *VadoIconButton) Disabled() bool {
	return b.disabled
}

func (b *VadoIconButton) Cursor() desktop.Cursor {
	if b.disabled {
		return desktop.DefaultCursor
	}
	return desktop.PointerCursor // курсор "рука"
}

func (b *VadoIconButton) updateVisualState() {
	if b.bg == nil || b.img == nil {
		return
	}

	switch {
	case b.disabled:
		// тусклая иконка, еле заметный фон
		b.img.Translucency = 0.5
		b.bg.FillColor = color.NRGBA{} // без подсветки
	case b.hovered:
		// яркая иконка, мягкий серый круг
		b.img.Translucency = 0.0
		b.bg.FillColor = color.NRGBA{R: 220, G: 220, B: 220, A: 80}
	default:
		// нормальное состояние
		b.img.Translucency = 0.0
		b.bg.FillColor = color.NRGBA{}
	}

	canvas.Refresh(b.bg)
	canvas.Refresh(b.img)
}

func (b *VadoIconButton) MouseIn(*desktop.MouseEvent) {
	if b.disabled {
		return
	}
	b.hovered = true
	b.updateVisualState()
}

func (b *VadoIconButton) MouseOut() {
	b.hovered = false
	b.updateVisualState()
}
func (b *VadoIconButton) MouseMoved(*desktop.MouseEvent) {}
