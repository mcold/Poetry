package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) makeUI() (*widget.List, *widget.Slider, *fyne.Container, *fyne.Container, *fyne.Container) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	nPcntDef, err := strconv.Atoi(os.Getenv("PCNT_DEF"))
	if err != nil {
		log.Fatal("Error converting PCNT_DEF value")
	}

	f := float64(nPcntDef / 100)

	idItem, err := strconv.Atoi(os.Getenv("ID_ITEM"))
	if err != nil {
		log.Fatal("Error converting ID_ITEM value")
	}
	id_item := idItem

	Lines, err := app.DB.LinesByItem(id_item, app.PageNum, app.PageSize) // TODO: change arg value
	if err != nil {
		log.Fatal("Error loading data from DB")
	}

	for _, line := range Lines {
		app.LinesArr = append(app.LinesArr, line.Ttext)
	}

	app.TransArr, err = app.DB.TransArrByLine(id_item, app.PageNum, app.PageSize)
	if err != nil {
		log.Fatal("Error loading data from DB")
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

	app.ListLines = l_lines

	data := binding.BindFloat(&f)
	slide := widget.NewSliderWithData(0, 1, data)
	// nStepSide, err := strconv.Atoi(os.Getenv("PCNT_STEP"))
	// if err != nil {
	// 	log.Fatal("Error converting PCNT_STEP value")
	// }
	l_lines.OnSelected = app.refreshData

	slide.Step = 0.1
	slide.OnChanged = func(v float64) {
		data.Set(v)
		app.LinesArr = hide(app.LinesArrDef, int(v*100))
		app.ListLinesData.Reload()
	}

	// data.AddListener(binding.NewDataListener(func() {
	// 	fmt.Println(data.Get())
	// }))

	btnPcnt := container.NewGridWithColumns(4,
		widget.NewButton("0%", func() {
			data.Set(0)
			app.LinesArr = hide(app.LinesArrDef, 0)
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

	btnPage := container.NewGridWithColumns(4,
		widget.NewButton("<", func() {
			newPageNum := app.PageNum - 1
			if newPageNum <= 0 {
				newPageNum = 1
			}
			if newPageNum != app.PageNum {
				app.PageNum = newPageNum
				app.LinesArr = nil
				app.TransArr = nil

				Lines, err := app.DB.LinesByItem(app.ItemID, app.PageNum, app.PageSize) // TODO: change arg value
				if err != nil {
					log.Fatal("Error loading data from DB")
				}
				for _, line := range Lines {
					app.LinesArr = append(app.LinesArr, line.Ttext)
				}
				app.TransArr, err = app.DB.TransArrByLine(app.ItemID, app.PageNum, app.PageSize)
				if err != nil {
					log.Fatal("Error loading data from DB")
				}

				app.LinesArrDef = app.LinesArr
				app.ListLinesData.Reload()
			}
		}),
		widget.NewButton(">", func() {
			newPageNum := app.PageNum + 1

			Lines, err := app.DB.LinesByItem(app.ItemID, newPageNum, app.PageSize) // TODO: change arg value
			if err != nil {
				log.Fatal("Error loading data from DB")
			}
			app.TransArr = nil
			app.TransArr, err = app.DB.TransArrByLine(app.ItemID, app.PageNum, app.PageSize)
			if err != nil {
				log.Fatal("Error loading data from DB")
			}

			if len(Lines) > 0 {
				app.PageNum = newPageNum
				app.LinesArr = nil
				for _, line := range Lines {
					app.LinesArr = append(app.LinesArr, line.Ttext)
				}
				app.LinesArrDef = app.LinesArr
				app.ListLinesData.Reload()
			}
		}))

	btnItem := container.NewGridWithColumns(2,
		widget.NewButton("👆", func() {
			newItemID := app.ItemID - 1
			newLines, err := app.DB.LinesByItem(newItemID, 1, app.PageSize) // TODO: change arg value
			if err != nil {
				log.Fatal("Error loading data from DB")
			}
			if len(newLines) > 0 {
				app.ItemID = newItemID
				app.LinesArr = nil
				app.TransArr = nil

				for _, line := range newLines {
					app.LinesArr = append(app.LinesArr, line.Ttext)
				}
				app.PageNum = 1
				app.LinesArrDef = app.LinesArr
				app.ListLinesData.Reload()
			}
		}),
		widget.NewButton("👇", func() {
			newItemID := app.ItemID + 1
			newLines, err := app.DB.LinesByItem(newItemID, 1, app.PageSize) // TODO: change arg value
			if err != nil {
				log.Fatal("Error loading data from DB")
			}
			if len(newLines) > 0 {
				app.ItemID = newItemID
				app.LinesArr = nil
				app.TransArr = nil

				for _, line := range newLines {
					app.LinesArr = append(app.LinesArr, line.Ttext)
				}
				app.PageNum = 1
				app.LinesArrDef = app.LinesArr
				app.ListLinesData.Reload()
			}
		}))

	return l_lines, slide, btnPcnt, btnPage, btnItem
}

func (app *Config) refreshData(id int) {
	app.NumLineActive = id
	for i := 0; i < app.PageSize; i++ {
		if i != id {
			app.ListLines.Unselect(i)
		}
	}
	app.ListLines.Refresh()
}
