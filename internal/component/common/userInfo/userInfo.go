package userInfo

import (
	"fmt"
	"vado-client/internal/app"
	"vado-client/internal/app/prefs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewUserInfo(ctx *app.Context) *fyne.Container {
	usernameTxt := widget.NewRichTextFromMarkdown("")
	updateUsernameText(ctx.Prefs, usernameTxt)
	ctx.App.Preferences().AddChangeListener(func() {
		updateUsernameText(ctx.Prefs, usernameTxt)
	})

	enterBtn := widget.NewButton("Log in", func() {
		callBack := func(token string) {
			ctx.Log.Debugw("Create token", "token", token)
		}
		ShowLoginDialog(ctx, &callBack)
	})

	quitBtn := widget.NewButton("Log out", func() {
		ctx.Prefs.Reset()
	})

	updateVisibility(ctx, enterBtn, quitBtn)

	ctx.App.Preferences().AddChangeListener(func() {
		updateVisibility(ctx, enterBtn, quitBtn)
	})

	return container.NewHBox(usernameTxt, enterBtn, quitBtn)
}

func updateUsernameText(prefs *prefs.Prefs, txt *widget.RichText) {
	txt.ParseMarkdown(fmt.Sprintf("User: **%s**", func() string {
		if prefs.Username() != "" {
			return prefs.Username()
		}
		return "Guest"
	}()))
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
