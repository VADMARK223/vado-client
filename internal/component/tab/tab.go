package tab

import (
	"vado-client/internal/appcontext"
	"vado-client/internal/component/tab/chat"
	"vado-client/internal/component/tab/hello"

	"fyne.io/fyne/v2/container"
)

func NewTab(ctx *appcontext.AppContext) *container.AppTabs {
	tabs := container.NewAppTabs(
		container.NewTabItem("Чат", chat.NewChat(ctx)),
		container.NewTabItem("Проверка", hello.NewHelloBox(ctx)),
	)
	return tabs
}
