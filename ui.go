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

func (app *Config) makeUI() (*widget.List, *widget.Slider, *widget.Select, *widget.Label, *fyne.Container, *fyne.Container, *fyne.Container) {
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
	l_lines.OnSelected = app.refreshData

	slide.Step = 0.1
	slide.OnChanged = func(v float64) {
		data.Set(v)
		app.PcntHide = v * 100
		app.LinesArr = hide(app.LinesArrDef, int(v*100))
		app.ListLinesData.Reload()
	}

	l_authors, err := app.DB.AllAuthors()
	if err != nil {
		log.Fatal("Error loading data from DB")
	}

	var l_authors_str []string

	for _, author := range l_authors {
		l_authors_str = append(l_authors_str, author.Name)
	}

	app.ItemNameData = binding.BindString(&app.ItemName)
	lblItem := widget.NewLabelWithData(app.ItemNameData)

	selAuthor := widget.NewSelect(l_authors_str, func(selected string) {
		oldAuthorID := app.AuthorID
		for _, author := range l_authors {
			if author.Name == selected {
				app.AuthorID = int(author.ID)

				app.AuthorItemsArr, err = app.DB.ItemsByAuthor(app.AuthorID)
				if err != nil {
					log.Fatal("Error loading data from DB")
				}
				if app.AuthorID != oldAuthorID {
					app.ItemID = int(app.AuthorItemsArr[0].ID)
					app.LinesArr = nil
					app.TransArr = nil

					newLines, err := app.DB.LinesByItem(app.ItemID, 1, app.PageSize) // TODO: change arg value
					if err != nil {
						log.Fatal("Error loading data from DB")
					}

					for _, line := range newLines {
						app.LinesArr = append(app.LinesArr, line.Ttext)
					}

					app.PageNum = 1
					app.LinesArrDef = app.LinesArr
					app.LinesArr = hide(app.LinesArrDef, int(app.PcntHide))
					app.ListLinesData.Reload()

					app.ItemsNameArr = nil
					for _, item := range app.AuthorItemsArr {
						app.ItemsNameArr = append(app.ItemsNameArr, item.Name)
					}
					app.ItemName = app.ItemsNameArr[0]
					app.ItemNameData.Reload()
				}
			}
		}
	})

	selAuthor.SetSelected(l_authors_str[0])

	btnPcnt := container.NewGridWithColumns(4,
		widget.NewButton("0%", func() {
			data.Set(0)
			app.PcntHide = 0
			app.LinesArr = hide(app.LinesArrDef, 0)
			app.ListLinesData.Reload()
		}),
		widget.NewButton("50%", func() {
			data.Set(0.5)
			app.PcntHide = 50
			app.LinesArr = hide(app.LinesArrDef, 50)
			app.ListLinesData.Reload()
		}),
		widget.NewButton("70%", func() {
			data.Set(0.7)
			app.PcntHide = 70
			app.LinesArr = hide(app.LinesArrDef, 70)
			app.ListLinesData.Reload()
		}),
		widget.NewButton("100%", func() {
			data.Set(1)
			app.PcntHide = 100
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
				app.LinesArr = hide(app.LinesArrDef, int(app.PcntHide))
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
				app.LinesArr = hide(app.LinesArrDef, int(app.PcntHide))
				app.ListLinesData.Reload()
			}
		}))

	btnItem := container.NewGridWithColumns(2,
		widget.NewButton("👆", func() {
			prevItemID := 0
			prevItemName := app.ItemName
			for _, item := range app.AuthorItemsArr {
				if item.ID <= int64(app.ItemID) {
					if item.ID == int64(app.ItemID) {
						if prevItemID != 0 {
							app.ItemID = prevItemID
							break
						}
					} else {
						prevItemID = int(item.ID)
						prevItemName = item.Name
					}
				}
			}

			app.ItemName = prevItemName
			app.ItemNameData.Reload()

			app.LinesArr = nil
			app.TransArr = nil
			// fmt.Println("There")
			newLines, err := app.DB.LinesByItem(app.ItemID, 1, app.PageSize) // TODO: change arg value
			if err != nil {
				log.Fatal("Error loading data from DB")
			}

			for _, line := range newLines {
				app.LinesArr = append(app.LinesArr, line.Ttext)
			}

			app.PageNum = 1
			app.LinesArrDef = app.LinesArr
			app.LinesArr = hide(app.LinesArrDef, int(app.PcntHide))
			app.ListLinesData.Reload()

		}),
		widget.NewButton("👇", func() {
			b_found := false
			for _, item := range app.AuthorItemsArr {
				if b_found {
					app.ItemID = int(item.ID)
					app.ItemName = item.Name
					app.ItemNameData.Reload()
					break
				}
				if item.ID == int64(app.ItemID) {
					b_found = true
				}
			}

			// if another author
			if !b_found {
				app.ItemID = int(app.AuthorItemsArr[0].ID)
			}

			app.LinesArr = nil
			app.TransArr = nil

			newLines, err := app.DB.LinesByItem(app.ItemID, 1, app.PageSize) // TODO: change arg value
			if err != nil {
				log.Fatal("Error loading data from DB")
			}

			for _, line := range newLines {
				app.LinesArr = append(app.LinesArr, line.Ttext)
			}

			app.PageNum = 1
			app.LinesArrDef = app.LinesArr
			app.LinesArr = hide(app.LinesArrDef, int(app.PcntHide))
			app.ListLinesData.Reload()
		}))

	return l_lines, slide, selAuthor, lblItem, btnPcnt, btnPage, btnItem
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
