package mainTab

import (
	"vado-client/internal/app"
	"vado-client/internal/component/common/userInfo"
	"vado-client/internal/component/tabs/tabItem"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
)

type Tab struct {
	log *zap.SugaredLogger
	*fyne.Container
	btn *widget.Button
}

func New(ctx *app.Context) tabItem.TabContent {
	c := container.NewVBox(
		widget.NewLabel("Main page"),
		userInfo.NewUserInfo(ctx),
	)

	return &Tab{
		log:       ctx.Log,
		Container: c,
	}
}

func (t Tab) Open() {
	t.log.Debugw("Main tab opened")
}

func (t Tab) Close() {
	t.log.Debugw("Main tab closed")
}

func (t Tab) Canvas() fyne.CanvasObject {
	return t.Container
}
