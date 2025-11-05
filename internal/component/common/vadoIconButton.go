package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type VadoIconButton struct {
	widget.BaseWidget
	icon     fyne.Resource
	onTapped func()
	disabled bool
}

func NewVadoIconButton(icon fyne.Resource, tapped func()) *VadoIconButton {
	b := &VadoIconButton{icon: icon, onTapped: tapped}
	b.ExtendBaseWidget(b)
	return b
}

func (b *VadoIconButton) CreateRenderer() fyne.WidgetRenderer {
	img := widget.NewIcon(b.icon)
	return widget.NewSimpleRenderer(img)
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
	b.Refresh()
}

func (b *VadoIconButton) Enable() {
	b.disabled = false
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

func (b *VadoIconButton) MouseIn(*desktop.MouseEvent)    {}
func (b *VadoIconButton) MouseOut()                      {}
func (b *VadoIconButton) MouseMoved(*desktop.MouseEvent) {}
