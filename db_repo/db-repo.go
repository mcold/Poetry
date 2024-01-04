package db_repo

type Repository interface {
	Migrate() error
	InsertAuthor(author Author) (*Author, error)
	InsertItem(item Item) (*Item, error)
	InsertLine(line Line) (*Line, error)
	InsertTrans(trans Trans) (*Trans, error)
	AllAuthors() ([]Author, error)
	ItemsByAuthor(id_author int) ([]Item, error)
	LinesByItem(id_item int) ([]Line, error)
	TransByLine(id_line int) (*Trans, error)
}

type Author struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Descr string `json:"descr"`
}

type Item struct {
	ID        int64  `json:"id"`
	ID_Author int64  `json:"id_author"`
	Name      string `json:"name"`
	Next      int64  `json:"next"`
}

type Line struct {
	ID      int64  `json:"id"`
	ID_Item int64  `json:"id_item"`
	Num     int64  `json:"num"`
	Ttext   string `json:"ttext"`
	Stop    int64  `json:"stop"`
}

type Trans struct {
	ID      int64  `json:"id"`
	ID_Line int64  `json:"id_line"`
	Ttext   string `json:"ttext"`
	Lang    string `json:"lang"`
}
