package tab

import (
	"vado-client/internal/appcontext"
	"vado-client/internal/component/tab/chat"
	"vado-client/internal/component/tab/hello"
	"vado-client/internal/component/tab/login"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func NewTab(ctx *appcontext.AppContext, a fyne.App, token string) *container.AppTabs {
	tabs := container.NewAppTabs(
		container.NewTabItem("Логин", login.NewLogin(ctx, a)),
		container.NewTabItem("Проверка", hello.NewHelloBox(ctx, token)),
		container.NewTabItem("Чат", chat.NewChat(ctx)),
	)
	return tabs
}
