package main

import (
	"context"
	"fmt"
	"os"
	"time"
	pbServer "vado-client/api/pb/server"
	"vado-client/internal/appcontext"
	"vado-client/internal/component/common"
	"vado-client/internal/component/tab"
	"vado-client/internal/component/userInfo"
	"vado-client/internal/constants/color"
	"vado-client/internal/grpc/client"
	"vado-client/internal/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

const PortGRPC = "50051"

func main() {
	envAppID := os.Getenv("APP_ID")
	a := newApp(envAppID)
	w := newWindow(a, envAppID)

	clientGPRC, err := client.CreateClient(PortGRPC)

	if err != nil {
		fmt.Printf("Fail create gRPC client: %s", err.Error())
	}

	appCtx := appcontext.NewAppContext(initLogger(), clientGPRC, a, w)
	appCtx.Log.Infow("Start vado-client.", "time", time.Now().Format("2006-01-02 15:04:05"))

	bottomObjs := []fyne.CanvasObject{userInfo.CreateUserInfo(appCtx), layout.NewSpacer()}
	bottomObjs = append(bottomObjs, createServerStatus(appCtx)...)
	bottomBar := container.NewHBox(bottomObjs...)
	root := container.NewBorder(nil, bottomBar, nil, nil, tab.NewTab(appCtx))
	w.SetContent(root)
	w.ShowAndRun()
}

func newApp(envAppID string) fyne.App {
	var appID string
	if envAppID == "" {
		appID = "vado-client"
	} else {
		appID = fmt.Sprintf("vado-client-%s", envAppID)
	}
	return app.NewWithID(appID)
}

func newWindow(a fyne.App, envAppID string) fyne.Window {
	var title string
	if envAppID == "" {
		title = "Vado client (Single)"
	} else {
		title = fmt.Sprintf("Vado client (%s)", envAppID)
	}
	w := a.NewWindow(title)
	w.Resize(fyne.NewSize(450, 400))
	return w
}

func getStatusServer(appCtx *appcontext.AppContext) bool {
	serverClient := pbServer.NewServerServiceClient(appCtx.GRPC)
	pingResp, errPing := serverClient.Ping(context.Background(), &emptypb.Empty{})

	var result bool
	if errPing != nil {
		result = false
	} else {
		result = pingResp.Run
	}

	return result
}

func updateIndicatorColor(appCtx *appcontext.AppContext, indicator *common.Indicator) {
	if getStatusServer(appCtx) {
		indicator.SetFillColor(color.Green())
	} else {
		indicator.SetFillColor(color.Red())
	}
}

func createServerStatus(appCtx *appcontext.AppContext) []fyne.CanvasObject {
	fastModeTxt := widget.NewRichTextFromMarkdown("Server status:")
	indicator := common.NewIndicator(color.Orange(), fyne.NewSize(10, 10))
	refreshBtn := widget.NewButton("Обновить", func() {
		updateIndicatorColor(appCtx, indicator)
	})
	updateIndicatorColor(appCtx, indicator)

	box := container.NewHBox(fastModeTxt, container.NewCenter(indicator), refreshBtn)
	return []fyne.CanvasObject{box}
}

func initLogger() *zap.SugaredLogger {
	zapLogger, zapLoggerInitErr := logger.Init(true)
	if zapLoggerInitErr != nil {
		panic(zapLoggerInitErr)
	}
	defer func() { _ = zapLogger.Sync() }()

	return zapLogger
}
