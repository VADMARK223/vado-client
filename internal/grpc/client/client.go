package client

import (
	"fmt"
	"strconv"
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
	token := preferences.String(code.AccessToken)
	return token
}

func Logout(a fyne.App) {
	preferences := a.Preferences()
	preferences.RemoveValue(code.AccessToken)
	preferences.RemoveValue(code.RefreshToken)
	preferences.RemoveValue(code.ExpiresAt)
	preferences.RemoveValue(code.Username)
	preferences.RemoveValue(code.LastInput)
	preferences.RemoveValue(code.Id)
}

func GetUserID(a fyne.App) uint64 {
	preferences := a.Preferences()
	userID, parseUintErr := strconv.ParseUint(preferences.String(code.Id), 10, 64)
	if parseUintErr != nil {
		return 0
	}
	return userID
}

func GetUsername(a fyne.App) string {
	preferences := a.Preferences()
	username := preferences.String(code.Username)
	if username != "" {
		return username
	}
	return "Гость"
}

func GetLastInput(a fyne.App) string {
	preferences := a.Preferences()
	lastInput := preferences.String(code.LastInput)
	return lastInput
}

func SetLastInput(a fyne.App, value string) {
	preferences := a.Preferences()
	preferences.SetString(code.LastInput, value)
}

func IsAuth(a fyne.App) bool {
	return GetToken(a) != ""
}
