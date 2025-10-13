package tab

import (
	"vado-client/internal/tab/hello"

	"fyne.io/fyne/v2/container"
)

func NewTab() *container.AppTabs {
	tabs := container.NewAppTabs(
		container.NewTabItem("Hello", hello.NewHelloBox()),
	)
	return tabs
}
