package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type VadoClearableEntry struct {
	widget.BaseWidget
	entry     *widget.Entry
	clearBtn  *widget.Button
	container *fyne.Container

	userOnChanged func(string)
}

func NewVadoEntry() *VadoClearableEntry {
	e := &VadoClearableEntry{}
	e.ExtendBaseWidget(e)

	e.entry = widget.NewEntry()
	e.entry.SetPlaceHolder("Введите текст...")

	e.clearBtn = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		e.entry.SetText("")
	})
	e.clearBtn.Importance = widget.LowImportance
	e.clearBtn.Disable()

	e.entry.OnChanged = func(s string) {
		if s == "" {
			e.clearBtn.Disable()
		} else {
			e.clearBtn.Enable()
		}

		if e.userOnChanged != nil {
			e.userOnChanged(s)
		}
	}

	e.container = container.NewBorder(nil, nil, nil, e.clearBtn, e.entry)
	return e
}

func (b *VadoClearableEntry) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(b.container)
}

func (b *VadoClearableEntry) Text() string {
	return b.entry.Text
}

func (b *VadoClearableEntry) SetText(input string) {
	b.entry.SetText(input)
}

func (b *VadoClearableEntry) SetPlaceHolder(s string) {
	b.entry.SetPlaceHolder(s)
}

func (b *VadoClearableEntry) OnChanged(fn func(string)) {
	b.userOnChanged = fn
}

func (b *VadoClearableEntry) OnSubmitted(fn func(string)) {
	b.entry.OnSubmitted = fn
}
