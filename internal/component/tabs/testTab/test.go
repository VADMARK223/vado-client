package testTab

import (
	"fmt"
	"vado-client/internal/app"
	"vado-client/internal/app/keyman"
	"vado-client/internal/component/tabs/tabItem"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Tab struct {
	*fyne.Container
	btn    *widget.Button
	canvas fyne.Canvas
	unsub  func() // Отписка от события
	keyman *keyman.KeyManager
}

func (t *Tab) Open() {
	t.unsub = t.keyman.Subscribe(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {
			t.btn.OnTapped()
		}
	})
}

func (t *Tab) Close() {
	if t.unsub != nil {
		t.unsub()
	}
}

func (t *Tab) Canvas() fyne.CanvasObject {
	return t.Container
}

func New(appCtx *app.Context) tabItem.TabContent {
	btn := widget.NewButton("Click", func() {
		fmt.Println("Clicked!")
	})

	c := container.NewVBox(
		widget.NewLabel("Test Tab"),
		btn,
	)

	return &Tab{
		Container: c,
		btn:       btn,
		canvas:    appCtx.Win.Canvas(),
		keyman:    appCtx.KeyMan,
	}
}
