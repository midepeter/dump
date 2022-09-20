package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/xuri/excelize/v2"
)

var (
	sourceFile  string
	databaseUrl string
)

var (
	ErrSourceFile  = errors.New("Invalid Source file")
	ErrDatabaseUrl = errors.New("Invalid Database Url")
)

func main() {
	flag.StringVar(&sourceFile, "s", "", "this is specifies the directory to your source files that your are hoping to dump into your database")
	flag.StringVar(&databaseUrl, "d", "", "this specifies the databaseURL you are want to dump into")

	flag.Parse()

	if *sourceFile == "" {
		panic(ErrSourceFile)
	}

	if *databaseUrl == "" {
		panic(ErrDatabaseUrl)
	}
}

type Dumper struct {
	conn   *sql.Conn
	tables []string
}

//Dump dumps the excel file into the target database
func (d Dump) Dump() error {
}

type File struct {
	f *excelize.File
}

func (f File) ReadFile(filename string) error {
	if file != "" {
		return fmt.Errorf("Empty file name: Please provide a valid file name")
	}

	ext := strings.Split(filename, ".")
	if len(ext) != 2 {
		return fmt.Errorf("Invalid file format")
	}

	if ext[1] == "xlsx" || ext[1] == "xls" {
		openFile, err := excelize.OpenFile(filename)
		if err != nil {
			log.Fatalln("Unable to open excel file ", err)
			return err
		}

		f.f = openFile
	}

	return nil
}

func (f File) sheet() string {
	list := f.f.GetSheetList()
	if list < 1 {
		log.Println("Empty sheet list")
		return nil
	}

	return list[0]
}

func (f File) createTables() []string {
	tables := make([]string, 0)
	for _, v := range f.sheet() {
		cols := f.f.GetCols(v)
		if len(cols) > 1 {
			return nil
		}

		for _, v := range cols {
			tables = v
		}
	}

	return tables
}

func (f File) createValues(cols int) ([][]string, error) {
	values := make([][cols]string, 0)
	sheet := f.sheet()
	values, err := f.f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("Unable to get sheet values from rows")
	}

	return values
}

type DB struct {
	conn *sql.Conn
	db   *sql.DB
}

func NewDB(ctx context.Context, dsn, url string) *DB {
	if url == "" {
		return nil
	}

	db, err := sql.Open(dsn, url)
	if err != nil {
		log.Println("Error setting up database")
		return nil
	}

	dbConn, _ := db.Conn(ctx)

	return &DB{
		conn: dbConn,
		db:   db,
	}
}

func (d DB) CreateTable(ctx context.Context, tableCols []string, tablename string) error {
	if len(tableCols) < 1 {
		return errors.New("Invalid table column length")
	}

	stmt := fmt.Sprintf("CREATE TABLE %s IF NOT EXISTS", tablename)
	s, err := d.conn.PrepareContext(ctx, stmt)
	if err != nil {
		return fmt.Errorf("Invalid query: %v\n", err)
	}

	res, err := s.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("Error executing query: %v\n", err)
	}

	if res.RowsAffected > 0 {
		log.Printf("%s successfully created!", tablename)
	}
	return nil
}

func (d DB) Ping(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	d.conn.PingContext(ctx)
}

func (d DB) Migrate(tableCols []string, url string) error {
	m, err := migrate.New("", url)
	if err != nil {
		log.Fatalf("Error occured during migration: %v", err)
		return err
	}

	err = m.Up()
	if err != nil {
		log.Fatalf("Error occured during migration: %v", err)
		return err
	}

	return nil
}
