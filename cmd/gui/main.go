package main

import (
	"fmt"
	"time"
	"vado-client/internal/appcontext"
	"vado-client/internal/constants/code"
	"vado-client/internal/logger"
	"vado-client/internal/tab"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const Port = "50051"

func main() {
	a := app.NewWithID("vado-client")
	w := a.NewWindow("Vado client")

	client, err := createClient(Port)
	if err != nil {
		fmt.Printf("Fail create gRPC client: %s", err.Error())
	}

	appCtx := appcontext.NewAppContext(initLogger(), client, w)
	appCtx.Log.Infow("Start vado-client.", "time", time.Now().Format("2006-01-02 15:04:05"))

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			closeWindow(appCtx, w)
		}
	})

	preferences := a.Preferences()
	token := preferences.String(code.JwtToken)
	w.SetContent(tab.NewTab(appCtx, a, token))
	w.Resize(fyne.NewSize(400, 600))

	w.SetCloseIntercept(func() {
		closeWindow(appCtx, w)
	})

	w.ShowAndRun()
}

func initLogger() *zap.SugaredLogger {
	zapLogger, zapLoggerInitErr := logger.Init(true)
	if zapLoggerInitErr != nil {
		panic(zapLoggerInitErr)
	}
	defer func() { _ = zapLogger.Sync() }()

	return zapLogger
}

func createClient(port string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient("localhost:"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	return conn, err
}

func closeWindow(ctx *appcontext.AppContext, w fyne.Window) {
	_ = ctx.GRPC.Close()
	w.Close()
}
