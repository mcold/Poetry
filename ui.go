package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) makeUI() *widget.List {

	Lines, _ := app.DB.LinesByItem(1) // TODO: change arg value
	for _, line := range Lines {
		app.LinesArr = append(app.LinesArr, line.Ttext)
		fmt.Println(line.Ttext)
	}

	app.ListLinesData = binding.BindStringList(&app.LinesArr)

	l_lines := widget.NewListWithData(app.ListLinesData,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	return l_lines
}
