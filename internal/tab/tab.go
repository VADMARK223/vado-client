package tab

import (
	"vado-client/internal/appcontext"
	"vado-client/internal/tab/chat"
	"vado-client/internal/tab/hello"

	"fyne.io/fyne/v2/container"
)

func NewTab(ctx *appcontext.AppContext) *container.AppTabs {
	tabs := container.NewAppTabs(
		container.NewTabItem("Hello", hello.NewHelloBox(ctx)),
		container.NewTabItem("Chat", chat.NewChat(ctx)),
	)
	return tabs
}
