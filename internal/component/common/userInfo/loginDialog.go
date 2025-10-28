package userInfo

import (
	"context"
	"strings"
	"time"
	pb "vado-client/api/pb/auth"
	"vado-client/internal/app"
	"vado-client/internal/grpc/middleware"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ShowLoginDialog(ctx *app.Context, f *func(token string)) {
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Введите логин")

	passwordEntry := widget.NewEntry()
	passwordEntry.Password = true
	passwordEntry.SetPlaceHolder("Введите пароль")

	authClient := pb.NewAuthServiceClient(ctx.GRPC)
	var dlg dialog.Dialog
	doneBtn := widget.NewButton("Войти", func() {
		resp, err := authClient.Login(context.Background(), &pb.LoginRequest{
			Username: usernameEntry.Text,
			Password: passwordEntry.Text,
		})

		if err != nil {
			dialog.ShowInformation("Ошибка входа", err.Error(), ctx.Win)
			return
		}

		ctx.Prefs.SetUserID(resp.Id)
		ctx.Prefs.SetUsername(resp.Username)
		ctx.Prefs.SetAccessToken(resp.Token)
		ctx.Prefs.SetRefreshToken(resp.RefreshToken)
		ctx.Prefs.SetExpiresAt(time.Now().Add(middleware.TokenAliveMinutes * time.Minute).Unix())

		if f != nil {
			(*f)(resp.Token)
		}

		dlg.Hide()
	})
	doneBtn.Importance = widget.HighImportance
	doneBtn.Disable()

	cancelBtn := widget.NewButton("Отмена", func() {
		dlg.Hide()
	})

	usernameEntry.OnChanged = func(username string) {
		updateDoneBtnEnable(doneBtn, username, passwordEntry.Text)
	}

	passwordEntry.OnChanged = func(password string) {
		updateDoneBtnEnable(doneBtn, usernameEntry.Text, password)
	}

	form := widget.NewForm(
		widget.NewFormItem("Логин", usernameEntry),
		widget.NewFormItem("Пароль", passwordEntry),
	)

	content := container.NewVBox(form, container.NewHBox(layout.NewSpacer(), cancelBtn, doneBtn))

	dlg = dialog.NewCustomWithoutButtons("Вход", content, ctx.Win)
	dlg.Resize(fyne.NewSize(400, 180))
	dlg.Show()

	// Через короткое время после показа диалога — установить фокус
	time.AfterFunc(100*time.Millisecond, func() {
		fyne.Do(func() {
			fyne.CurrentApp().Driver().CanvasForObject(usernameEntry).Focus(usernameEntry)
		})
	})
}

func updateDoneBtnEnable(btn *widget.Button, username string, password string) {
	if getEnableDoneBtn(username, password) {
		btn.Enable()
	} else {
		btn.Disable()
	}
}

func getEnableDoneBtn(username string, password string) bool {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return false
	}

	return true
}
