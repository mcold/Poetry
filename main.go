package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"poetry/db_repo"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"

	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

type Config struct {
	App           fyne.App
	MainWindow    fyne.Window
	InfoLog       *log.Logger
	DB            db_repo.Repository
	LinesArr      []string
	LinesArrDef   []string // text by default
	ListLinesData binding.ExternalStringList
	ListTransData binding.ExternalStringList
}

const (
	difHeight = 50
	winWidth  = 800
)

func main() {
	var myApp Config
	var elemHeight float32

	elemHeight = 400

	fyneApp := app.NewWithID("poetry")
	myApp.App = fyneApp
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	sqlDB, err := myApp.connectSQL()
	if err != nil {
		log.Panic(err)
	}

	myApp.setupDB(sqlDB)

	myApp.MainWindow = fyneApp.NewWindow("poetry")

	listItems, slide, btns := myApp.makeUI()

	listItems.Resize(fyne.Size{Width: winWidth, Height: elemHeight})
	listItems.Move(fyne.Position{X: 0, Y: 0})

	elemHeight += difHeight

	slide.Resize(fyne.Size{Width: winWidth, Height: difHeight})
	slide.Move(fyne.Position{X: 0, Y: elemHeight - difHeight})

	elemHeight += difHeight

	btns.Resize(fyne.Size{Width: winWidth, Height: difHeight})
	btns.Move(fyne.Position{X: 0, Y: elemHeight - difHeight})

	elemHeight += difHeight

	c1 := container.NewWithoutLayout(listItems, slide, btns)

	myApp.MainWindow.SetContent(c1)

	myApp.MainWindow.Resize(fyne.Size{Width: winWidth, Height: elemHeight})

	myApp.MainWindow.ShowAndRun()
}

func (app *Config) connectSQL() (*sql.DB, error) {
	path := ""

	if os.Getenv("DB_PATH") != "" {
		path = os.Getenv("DB_PATH")
	} else {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		path = filepath.Join(dir, "DB.db")
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *Config) setupDB(sqlDB *sql.DB) {
	app.DB = db_repo.NewSQLiteRepository(sqlDB)

	dir, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	path := filepath.Join(dir, "DB.db")

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := app.DB.Migrate()

		if err != nil {
			os.Exit(1)
		}
	}

}
