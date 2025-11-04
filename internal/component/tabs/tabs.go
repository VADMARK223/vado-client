package tabs

import (
	"vado-client/internal/app"
	"vado-client/internal/component/tabs/chatTab"
	"vado-client/internal/component/tabs/hello"
	"vado-client/internal/component/tabs/kafkaTab"
	"vado-client/internal/component/tabs/mainTab"
	"vado-client/internal/component/tabs/tabItem"

	"fyne.io/fyne/v2/container"
)

func New(ctx *app.Context) *container.AppTabs {
	defaultTabIndex := ctx.Prefs.LastTabIndex()

	factories := map[*container.TabItem]func() tabItem.TabContent{}

	tabs := container.NewAppTabs(
		tabItem.New("Чат", func() tabItem.TabContent {
			return chatTab.New(ctx)
		}, factories),
		container.NewTabItem("SeyHello", hello.NewHelloBox(ctx)),
		tabItem.New("Kafka", func() tabItem.TabContent {
			return kafkaTab.New(ctx)
		}, factories),
		tabItem.New("Главная", func() tabItem.TabContent {
			return mainTab.New(ctx)
		}, factories),
	)

	loaded := map[*container.TabItem]tabItem.TabContent{}
	var active tabItem.TabContent
	tabs.OnSelected = func(item *container.TabItem) {
		if item == nil {
			return
		}

		// закрыть предыдущую вкладку
		if active != nil {
			active.Close()
		}

		// загрузка новой вкладки
		content, ok := loaded[item]
		if !ok {
			if factory, ok := factories[item]; ok {
				content = factory()
				loaded[item] = content
				delete(factories, item)
				item.Content = content.Canvas()
				tabs.Refresh()
			}
		}

		// открыть новую вкладку
		if content != nil {
			content.Open()
			active = content
		}
	}

	ctx.Log.Debugw("Last index", "index", defaultTabIndex)
	if len(tabs.Items) > 0 {
		if defaultTabIndex == 0 {
			item := tabs.Items[defaultTabIndex]
			if tabs.OnSelected != nil {
				tabs.OnSelected(item)
			}
		} else {
			tabs.SelectIndex(defaultTabIndex)
		}
	}

	ctx.AddCloseHandler(func() {
		ctx.Log.Debugw("Last index", "index", tabs.SelectedIndex())
		ctx.Prefs.SetLastTabIndex(tabs.SelectedIndex())
	})

	return tabs
}
