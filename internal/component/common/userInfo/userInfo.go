package userInfo

import (
	"fmt"
	"vado-client/internal/app"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewUserInfo(ctx *app.Context) *fyne.Container {
	username := ctx.Prefs.Username()
	userNameText := widget.NewRichTextFromMarkdown(fmt.Sprintf("Пользователь: **%s**", username))
	ctx.App.Preferences().AddChangeListener(func() {
		userNameText.ParseMarkdown(fmt.Sprintf("Пользователь: **%s**", username))
		userNameText.Refresh()
	})

	enterBtn := widget.NewButton("Вход", func() {
		callBack := func(token string) {
			ctx.Log.Debugw("Create token", "token", token)
		}
		ShowLoginDialog(ctx, &callBack)
	})

	quitBtn := widget.NewButton("Выход", func() {
		ctx.Prefs.Reset()
	})

	updateVisibility(ctx, enterBtn, quitBtn)

	ctx.App.Preferences().AddChangeListener(func() {
		updateVisibility(ctx, enterBtn, quitBtn)
	})

	return container.NewHBox(userNameText, enterBtn, quitBtn)
}

func updateVisibility(ctx *app.Context, enterBtn *widget.Button, quitBnt *widget.Button) {
	if ctx.Prefs.IsAuth() {
		enterBtn.Hide()
		quitBnt.Show()
	} else {
		quitBnt.Hide()
		enterBtn.Show()
	}
}
