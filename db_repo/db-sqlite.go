package db_repo

import (
	"database/sql"
	"os"
)

type SQLiteRepository struct {
	Conn *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		Conn: db,
	}
}

func (repo *SQLiteRepository) tblAuthor() error {
	query := `
	create table if not exists author(id    integer primary key autoincrement,
	     							  name  text,
									  descr text);
	`
	_, err := repo.Conn.Exec(query)
	return err
}

func (repo *SQLiteRepository) tblItem() error {
	query := `
	create table if not exists item(id          integer primary key autoincrement,
							   id_author   integer references author(id) on delete cascade,
							   name        text,
							   next        integer);
	`
	_, err := repo.Conn.Exec(query)
	return err
}

func (repo *SQLiteRepository) tblLine() error {
	query := `
	create table if not exists line(id      integer primary key autoincrement,
							        id_item integer references item(id) on delete cascade,
									num     integer,
									ttext   text,
									stop    integer);
	`
	_, err := repo.Conn.Exec(query)
	return err
}

func (repo *SQLiteRepository) tblTrans() error {
	query := `
	create table if not exists trans(id     integer primary key autoincrement,
									 id_line integer references line(id) on delete cascade,
									 ttext   text,
									 lang    text);
	`
	_, err := repo.Conn.Exec(query)
	return err
}

func (repo *SQLiteRepository) Migrate() error {

	err := repo.tblAuthor()
	if err != nil {
		os.Exit(1)
	}

	err = repo.tblItem()
	if err != nil {
		os.Exit(1)
	}

	err = repo.tblLine()
	if err != nil {
		os.Exit(1)
	}

	err = repo.tblTrans()
	if err != nil {
		os.Exit(1)
	}

	return err
}

func (repo *SQLiteRepository) InsertAuthor(author Author) (*Author, error) {
	stmt := "insert into author(name, descr) values (?, ?)"

	res, err := repo.Conn.Exec(stmt, author.Name, author.Descr)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	author.ID = id
	return &author, nil
}

func (repo *SQLiteRepository) InsertItem(item Item) (*Item, error) {
	stmt := "insert into item(id_author, name, next) values (?, ?, ?)"

	res, err := repo.Conn.Exec(stmt, item.ID_Author, item.Name, item.Next)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	item.ID = id
	return &item, nil
}

func (repo *SQLiteRepository) InsertLine(line Line) (*Line, error) {
	stmt := "insert into line(id_item, num, ttext, stop) values (?, ?, ?, ?)"

	res, err := repo.Conn.Exec(stmt, line.ID_Item, line.Num, line.Ttext, line.Stop)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	line.ID = id
	return &line, nil
}

func (repo *SQLiteRepository) InsertTrans(trans Trans) (*Trans, error) {
	stmt := "insert into trans(id_line, ttext, lang) values (?, ?, ?)"

	res, err := repo.Conn.Exec(stmt, trans.ID_Line, trans.Ttext, trans.Lang)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	trans.ID = id
	return &trans, nil
}

func (repo *SQLiteRepository) AllAuthors() ([]Author, error) {
	query := "select id, name, descr from author order by id asc"
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Author
	for rows.Next() {
		var a Author
		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Descr,
		)
		if err != nil {
			return nil, err
		}
		all = append(all, a)
	}

	return all, nil
}

func (repo *SQLiteRepository) ItemsByAuthor(id_author int) ([]Item, error) {
	query := "select id, id_author, name, next from item where id_author = ? order by id asc"
	rows, err := repo.Conn.Query(query, id_author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Item
	for rows.Next() {
		var it Item
		err := rows.Scan(
			&it.ID,
			&it.ID_Author,
			&it.Name,
			&it.Next,
		)
		if err != nil {
			return nil, err
		}
		all = append(all, it)
	}

	return all, nil
}

func (repo *SQLiteRepository) LinesByItem(id_item int, page_num int, page_size int) ([]Line, error) {
	query := `SELECT id,
	id_item,
	num,
	ttext,
	stop
FROM (SELECT id,
			id_item,
			num,
			ttext,
			stop,
			row_number() over (order by num) as row_num
		FROM line
	   WHERE id_item = ?
	   ORDER BY id ASC)
WHERE row_num > (? - 1) * ? AND row_num <= ? * ?`

	rows, err := repo.Conn.Query(query, id_item, page_num, page_size, page_num, page_size)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Line
	for rows.Next() {
		var l Line
		err := rows.Scan(
			&l.ID,
			&l.ID_Item,
			&l.Num,
			&l.Ttext,
			&l.Stop,
		)
		if err != nil {
			return nil, err
		}
		all = append(all, l)
	}

	return all, nil
}

func (repo *SQLiteRepository) TransByLine(id_line int) (*Trans, error) {
	row := repo.Conn.QueryRow("select id, id_line, ttext, lang from trans where id_line = ? order by id asc", id_line)
	var tr Trans
	err := row.Scan(
		&tr.ID,
		&tr.ID_Line,
		&tr.Ttext,
		&tr.Lang,
	)

	if err != nil {
		return nil, err
	}

	return &tr, nil
}
