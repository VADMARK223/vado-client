package userInfo

import (
	"fmt"
	"vado-client/internal/appcontext"
	"vado-client/internal/grpc/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateUserInfo(appCtx *appcontext.AppContext) *fyne.Container {
	userNameText := widget.NewRichTextFromMarkdown(fmt.Sprintf("Пользователь: **%s**", client.GetUsername(appCtx.App)))
	appCtx.App.Preferences().AddChangeListener(func() {
		userNameText.ParseMarkdown(fmt.Sprintf("Пользователь: **%s**", client.GetUsername(appCtx.App)))
		userNameText.Refresh()
	})

	enterBtn := widget.NewButton("Вход", func() {
		callBack := func(token string) {
			appCtx.Log.Debugw("Create token", "token", token)
		}
		ShowLoginDialog(appCtx, &callBack)
	})

	quitBtn := widget.NewButton("Выход", func() {
		client.Logout(appCtx.App)
	})

	updateVisibility(appCtx.App, enterBtn, quitBtn)

	appCtx.App.Preferences().AddChangeListener(func() {
		updateVisibility(appCtx.App, enterBtn, quitBtn)
	})

	return container.NewHBox(userNameText, enterBtn, quitBtn)
}

func updateVisibility(a fyne.App, enterBtn *widget.Button, quitBnt *widget.Button) {
	if client.IsAuth(a) {
		enterBtn.Hide()
		quitBnt.Show()
	} else {
		quitBnt.Hide()
		enterBtn.Show()
	}
}
