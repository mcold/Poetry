package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"poetry/db_repo"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"database/sql"

	"github.com/joho/godotenv"

	_ "github.com/glebarez/go-sqlite"
)

type Config struct {
	App            fyne.App
	AuthorID       int
	AuthorItemsArr []db_repo.Item
	MainWindow     fyne.Window
	InfoLog        *log.Logger
	ItemID         int
	DB             db_repo.Repository
	ItemsNameArr   []string
	ItemName       string
	ItemNameData   binding.ExternalString
	LinesArr       []string
	TransArr       []db_repo.Trans
	LinesArrDef    []string
	ListLines      *widget.List
	ListLinesData  binding.ExternalStringList
	ListTransData  binding.ExternalStringList
	ListEditArr    []string
	ListEditData   binding.ExternalStringList
	NumLineActive  int
	PageNum        int
	PageSize       int
	PcntHide       float64
}

const (
	difHeight = 50
	winWidth  = 800
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var myApp Config
	var elemHeight float32

	myApp.PageNum = 1
	nPageSize, err := strconv.Atoi(os.Getenv("PAGE_SIZE"))
	if err != nil {
		log.Fatal("Error converting PAGE_SIZE value")
	}
	myApp.PageSize = nPageSize

	elemHeight = 400

	fyneApp := app.NewWithID("POETRY")
	myApp.App = fyneApp
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	/*TODO: change to storing flag in DB*/
	myApp.ItemID, err = strconv.Atoi(os.Getenv("ID_ITEM"))
	if err != nil {
		log.Fatal("Error converting ItemID value")
	}

	sqlDB, err := myApp.connectSQL()
	if err != nil {
		log.Panic(err)
	}

	myApp.setupDB(sqlDB)

	myApp.MainWindow = fyneApp.NewWindow("POETRY")

	listItems, slide, selAuthor, lblItem, btnPcnt, btnPage, btnItem := myApp.makeUI()

	selAuthor.Resize(fyne.Size{Width: winWidth / 2, Height: difHeight})
	selAuthor.Move(fyne.Position{X: 0, Y: 0})

	lblItem.Resize(fyne.Size{Width: winWidth / 2, Height: difHeight})
	lblItem.Move(fyne.Position{X: winWidth / 2, Y: 0})

	elemHeight += difHeight

	listItems.Resize(fyne.Size{Width: winWidth, Height: elemHeight})
	listItems.Move(fyne.Position{X: 0, Y: difHeight})

	elemHeight += difHeight

	slide.Resize(fyne.Size{Width: winWidth, Height: difHeight})
	slide.Move(fyne.Position{X: 0, Y: elemHeight - difHeight})

	elemHeight += difHeight

	btnPcnt.Resize(fyne.Size{Width: winWidth, Height: difHeight})
	btnPcnt.Move(fyne.Position{X: 0, Y: elemHeight - difHeight})

	elemHeight += difHeight

	btnPage.Resize(fyne.Size{Width: winWidth, Height: difHeight})
	btnPage.Move(fyne.Position{X: 0, Y: elemHeight - difHeight})

	// elemHeight += difHeight

	btnItem.Resize(fyne.Size{Width: winWidth / 2, Height: difHeight})
	btnItem.Move(fyne.Position{X: winWidth / 2, Y: elemHeight - difHeight})

	elemHeight += difHeight

	c1 := container.NewWithoutLayout(selAuthor, lblItem, listItems, slide, btnPcnt, btnPage, btnItem)

	myApp.MainWindow.SetContent(c1)

	myApp.MainWindow.Resize(fyne.Size{Width: winWidth, Height: elemHeight})

	myApp.MainWindow.ShowAndRun()
}

func (app *Config) connectSQL() (*sql.DB, error) {
	path := ""

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbName := os.Getenv("DB_NAME")

	if os.Getenv("DB_PATH") != "" {
		path = os.Getenv("DB_PATH")
	} else {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		path = filepath.Join(dir, dbName)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *Config) setupDB(sqlDB *sql.DB) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbName := os.Getenv("DB_NAME")

	app.DB = db_repo.NewSQLiteRepository(sqlDB)

	dir, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	path := filepath.Join(dir, dbName)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := app.DB.Migrate()

		if err != nil {
			os.Exit(1)
		}
	}

}
