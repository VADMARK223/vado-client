package tab

import (
	"vado-client/internal/appcontext"
	"vado-client/internal/component/tab/chat"
	"vado-client/internal/component/tab/hello"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func NewTab(ctx *appcontext.AppContext, a fyne.App) *container.AppTabs {
	tabs := container.NewAppTabs(
		container.NewTabItem("Проверка", hello.NewHelloBox(ctx, a)),
		container.NewTabItem("Чат", chat.NewChat(ctx, a)),
	)
	return tabs
}
