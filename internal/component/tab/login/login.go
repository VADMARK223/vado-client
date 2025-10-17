package login

import (
	"context"
	"vado-client/internal/appcontext"
	"vado-client/internal/constants/code"
	"vado-client/internal/grpc/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	pb "vado-client/api/pb/auth"
)

func NewLogin(appCtx *appcontext.AppContext, a fyne.App) *fyne.Container {
	usernameInput := widget.NewEntry()
	usernameInput.SetPlaceHolder("Введите имя")

	passwordInput := widget.NewEntry()
	passwordInput.SetPlaceHolder("Введите пароль")

	authClient := pb.NewAuthServiceClient(appCtx.GRPC)

	enterBtn := widget.NewButton("Вход", func() {
		resp, err := authClient.Login(context.Background(), &pb.LoginRequest{
			Username: usernameInput.Text,
			Password: passwordInput.Text,
		})

		if err != nil {
			dialog.ShowInformation("Ошибка входа", err.Error(), appCtx.Win)
			return
		}

		appCtx.Log.Debugf("JWT: %s", resp.Token)
		prefs := a.Preferences()
		prefs.SetString(code.JwtToken, resp.Token)
	})

	quitBtn := widget.NewButton("Выход", func() {
		client.ResetToken(a)
	})

	return container.NewVBox(usernameInput, passwordInput, enterBtn, quitBtn)
}
