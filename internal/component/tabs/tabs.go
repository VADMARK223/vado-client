package tabs

import (
	"vado-client/internal/app"
	"vado-client/internal/component/tabs/chatTab"
	"vado-client/internal/component/tabs/hello"
	"vado-client/internal/component/tabs/tabItem"
	"vado-client/internal/component/tabs/testTab"

	"fyne.io/fyne/v2/container"
)

func New(appCtx *app.Context) *container.AppTabs {
	factories := map[*container.TabItem]func() tabItem.TabContent{}

	tabs := container.NewAppTabs(
		container.NewTabItem("Чат", chatTab.New(appCtx)),
		container.NewTabItem("Проверка", hello.NewHelloBox(appCtx)),
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

		//if !ok {
		//	appCtx.Log.Debugw("Tab is not loaded")
		//}
		//if item == nil || loaded[item] {
		//	return
		//}

		/*if factory, ok := factories[item]; ok {
			go func() {
				// Обновляем UI безопасно
				fyne.Do(func() {
					item.Content = factory()
					loaded[item] = true
					delete(factories, item)
				})
			}()
		}*/
	}

	return tabs
}
