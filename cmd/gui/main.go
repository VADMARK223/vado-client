package main

import (
	"context"
	"fmt"
	"os"
	"time"
	pbPing "vado-client/api/pb/ping"
	"vado-client/internal/app"
	"vado-client/internal/app/logger"
	"vado-client/internal/component/common"
	"vado-client/internal/component/common/userInfo"
	"vado-client/internal/component/tabs"
	"vado-client/internal/config/color"
	"vado-client/internal/grpc/client"
	"vado-client/internal/utils"

	"fyne.io/fyne/v2"
	f "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/niemeyer/pretty"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	envAppID := os.Getenv("APP_ID")
	a := newApp(envAppID)
	w := newWindow(a, envAppID)

	clientGPRC, err := client.CreateClient()

	if err != nil {
		fmt.Printf("Fail create gRPC client: %s", err.Error())
	}

	zapLogger := logger.Init(true)
	defer func() { _ = zapLogger.Sync() }()

	appCtx := app.NewAppCtx(zapLogger, clientGPRC, a, w)
	appCtx.Log.Infow("Start vado-client.", "time", utils.FormatTime(time.Now()))

	bottomObjs := []fyne.CanvasObject{userInfo.NewUserInfo(appCtx), layout.NewSpacer()}
	bottomObjs = append(bottomObjs, createServerStatus(appCtx)...)
	bottomBar := container.NewHBox(bottomObjs...)
	root := container.NewBorder(nil, bottomBar, nil, nil, tabs.New(appCtx))
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
	return f.NewWithID(appID)
}

func newWindow(a fyne.App, envAppID string) fyne.Window {
	var title string
	if envAppID == "" {
		title = "Vado client (Single)"
	} else {
		title = fmt.Sprintf("Vado client (%s)", envAppID)
	}
	w := a.NewWindow(title)
	w.Resize(fyne.NewSize(455, 703))
	return w
}

func getStatusServer(appCtx *app.Context) (result bool) {
	serverClient := pbPing.NewPingServiceClient(appCtx.GRPC)
	pingResp, errPing := serverClient.Ping(context.Background(), &emptypb.Empty{})

	if errPing != nil {
		dialog.ShowInformation("Ошибка", pretty.Sprintf("%s", errPing), appCtx.Win)
	} else {
		result = pingResp.Run
	}
	appCtx.Log.Infow("Get status server.", "errPing", errPing)
	return
}

func updateIndicatorColor(appCtx *app.Context, indicator *common.Indicator) {
	if getStatusServer(appCtx) {
		indicator.SetFillColor(color.Green())
	} else {
		indicator.SetFillColor(color.Red())
	}
}

func createServerStatus(appCtx *app.Context) []fyne.CanvasObject {
	fastModeTxt := widget.NewRichTextFromMarkdown("Server:")
	indicator := common.NewIndicator(color.Orange(), fyne.NewSize(10, 10))
	refreshBtn := widget.NewButton("Refresh", func() {
		updateIndicatorColor(appCtx, indicator)
	})
	updateIndicatorColor(appCtx, indicator)

	box := container.NewHBox(fastModeTxt, container.NewCenter(indicator), refreshBtn)
	return []fyne.CanvasObject{box}
}
