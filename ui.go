package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) makeUI() (*widget.List, *widget.Slider, *fyne.Container) {
	f := 0.2

	Lines, _ := app.DB.LinesByItem(1) // TODO: change arg value
	for _, line := range Lines {
		app.LinesArr = append(app.LinesArr, line.Ttext)
	}

	app.LinesArrDef = app.LinesArr
	app.ListLinesData = binding.BindStringList(&app.LinesArr)

	l_lines := widget.NewListWithData(app.ListLinesData,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	// btnLeft := widget.NewButtonWithIcon("Left", fyne.NewStaticResource())

	data := binding.BindFloat(&f)
	slide := widget.NewSliderWithData(0, 1, data)
	slide.Step = 0.01
	slide.OnChanged = func(v float64) {
		fmt.Println(v)
		data.Set(v)
		app.LinesArr = hide(app.LinesArrDef, int(v*100))
		app.ListLinesData.Reload()
	}

	buttons := container.NewGridWithColumns(4,
		widget.NewButton("0%", func() {
			data.Set(0)
			app.LinesArr = hide(app.LinesArrDef, 0)
			app.ListLinesData.Reload()
		}),
		widget.NewButton("30%", func() {
			data.Set(0.3)
			app.LinesArr = hide(app.LinesArrDef, 30)
			app.ListLinesData.Reload()
		}),
		widget.NewButton("50%", func() {
			data.Set(0.5)
			app.LinesArr = hide(app.LinesArrDef, 50)
			app.ListLinesData.Reload()
		}),
		widget.NewButton("70%", func() {
			data.Set(0.7)
			app.LinesArr = hide(app.LinesArrDef, 70)
			app.ListLinesData.Reload()
		}),
		widget.NewButton("100%", func() {
			data.Set(1)
			app.LinesArr = hide(app.LinesArrDef, 100)
			app.ListLinesData.Reload()
		}))

	return l_lines, slide, buttons
}
