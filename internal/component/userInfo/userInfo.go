package userInfo

import (
	"fmt"
	"vado-client/internal/appcontext"
	"vado-client/internal/grpc/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateUserInfo(appCtx *appcontext.AppContext, a fyne.App) *fyne.Container {
	userInfo := widget.NewRichTextFromMarkdown(fmt.Sprintf("Пользователь: **%s**", client.GetUsername(a)))
	a.Preferences().AddChangeListener(func() {
		userInfo.ParseMarkdown(fmt.Sprintf("Пользователь: **%s**", client.GetUsername(a)))
		userInfo.Refresh()
	})

	enterBtn := widget.NewButton("Вход", func() {
		ShowLoginDialog(appCtx, a, func(token string) {
			appCtx.Log.Debugw("Create token", "token", token)
		})

	})

	quitBtn := widget.NewButton("Выход", func() {
		client.Logout(a)
	})

	updateVisibility(a, enterBtn, quitBtn)

	a.Preferences().AddChangeListener(func() {
		updateVisibility(a, enterBtn, quitBtn)
	})

	return container.NewHBox(userInfo, enterBtn, quitBtn)
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
