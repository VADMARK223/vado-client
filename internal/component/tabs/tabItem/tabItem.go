package tabItem

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TabContent interface {
	Open()
	Close()
	Canvas() fyne.CanvasObject
}

func New(text string, factory func() TabContent, factories map[*container.TabItem]func() TabContent) *container.TabItem {
	tab := container.NewTabItem(text, widget.NewProgressBarInfinite())
	factories[tab] = factory
	return tab
}
