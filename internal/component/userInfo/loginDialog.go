package userInfo

import (
	"context"
	"strconv"
	"strings"
	"time"
	pb "vado-client/api/pb/auth"
	"vado-client/internal/appcontext"
	"vado-client/internal/constants/code"
	"vado-client/internal/grpc/middleware"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ShowLoginDialog(appCtx *appcontext.AppContext, f *func(token string)) {
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Введите логин")

	passwordEntry := widget.NewEntry()
	passwordEntry.SetPlaceHolder("Введите пароль")

	authClient := pb.NewAuthServiceClient(appCtx.GRPC)
	var dlg dialog.Dialog
	doneBtn := widget.NewButton("Войти", func() {
		resp, err := authClient.Login(context.Background(), &pb.LoginRequest{
			Username: usernameEntry.Text,
			Password: passwordEntry.Text,
		})

		if err != nil {
			dialog.ShowInformation("Ошибка входа", err.Error(), appCtx.Win)
			return
		}

		prefs := appCtx.App.Preferences()
		prefs.SetString(code.AccessToken, resp.Token)
		prefs.SetString(code.RefreshToken, resp.RefreshToken)
		prefs.SetInt(code.ExpiresAt, int(time.Now().Add(middleware.TokenAliveMinutes*time.Minute).Unix()))
		prefs.SetString(code.Username, resp.Username)
		prefs.SetString(code.Id, strconv.FormatUint(resp.Id, 10))

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

	dlg = dialog.NewCustomWithoutButtons("Вход", content, appCtx.Win)
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
