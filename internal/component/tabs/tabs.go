package tabs

import (
	"vado-client/internal/app"
	"vado-client/internal/component/tabs/chatTab"
	"vado-client/internal/component/tabs/hello"
	"vado-client/internal/component/tabs/kafkaTab"
	"vado-client/internal/component/tabs/tabItem"
	"vado-client/internal/component/tabs/testTab"

	"fyne.io/fyne/v2/container"
)

const defaultTabIndex = 2

func New(appCtx *app.Context) *container.AppTabs {
	factories := map[*container.TabItem]func() tabItem.TabContent{}

	tabs := container.NewAppTabs(
		tabItem.New("Чат", func() tabItem.TabContent {
			return chatTab.New(appCtx)
		}, factories),
		container.NewTabItem("Проверка", hello.NewHelloBox(appCtx)),
		tabItem.New("Kafka", func() tabItem.TabContent {
			return kafkaTab.New(appCtx)
		}, factories),
		tabItem.New("Test", func() tabItem.TabContent {
			return testTab.New(appCtx)
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

	if len(tabs.Items) > 0 {
		item := tabs.Items[defaultTabIndex]
		if tabs.OnSelected != nil {
			tabs.OnSelected(item)
		}
		tabs.SelectIndex(defaultTabIndex)
	}

	return tabs
}
