package client

import (
	"fmt"
	"vado-client/internal/constants/code"

	"fyne.io/fyne/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateClient(port string) (*grpc.ClientConn, error) {
	target := "localhost:" + port
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания gRPC клиента: %w", err)
	}

	return conn, nil
}

func GetToken(a fyne.App) string {
	preferences := a.Preferences()
	token := preferences.String(code.JwtToken)
	return token
}

func Logout(a fyne.App) {
	preferences := a.Preferences()
	preferences.RemoveValue(code.JwtToken)
	preferences.RemoveValue(code.Username)
}

func GetUsername(a fyne.App) string {
	preferences := a.Preferences()
	username := preferences.String(code.Username)
	if username != "" {
		return username
	}
	return "Гость"
}

func IsAuth(a fyne.App) bool {
	return GetToken(a) != ""
}
