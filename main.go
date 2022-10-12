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

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
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

	if sourceFile == "" {
		panic(ErrSourceFile)
	}

	if databaseUrl == "" {
		panic(ErrDatabaseUrl)
	}

	f, err := ReadFile(sourceFile)
	if err != nil {
		panic(err)
	}

	fileName := strings.Split(sourceFile, ".")

	dsn := "postgres"
	url := "postgres://midepeter:password@localhost:5432/sample"
	db := NewDB(context.Background(), dsn, url, &f)

	d := Dumper{
		f:  f,
		db: db,
	}

	fmt.Println("The db conn", db)
	fmt.Println("The filename is", fileName)
	err = d.Dump(fileName[0], 4)
	if err != nil {
		panic(err)
	}
}

type Dumper struct {
	f  File
	db *DB
}

//Dump dumps the excel file into the target database
func (d Dumper) Dump(name string, indexKey int) error {
	rows, err := d.f.GetRows()
	if err != nil {
		return fmt.Errorf("Unable to fetch file rows ", err)
	}

	err = d.db.CreateTable(context.Background(), rows[indexKey], name)
	if err != nil {
		return fmt.Errorf("Unable to create table: %v", err)
	}

	err = d.db.Import(rows, 5)
	if err != nil {
		return fmt.Errorf("Unable to import rows into the database %v", err)
	}

	return nil
}

type File struct {
	f *excelize.File
}

func ReadFile(filename string) (File, error) {
	var f File
	if filename == "" {
		return f, fmt.Errorf("Empty file name: Please provide a valid file name")
	}

	ext := strings.Split(filename, ".")
	if len(ext) != 2 {
		return f, fmt.Errorf("Invalid file format")
	}

	if ext[1] == "xlsx" || ext[1] == "xls" {
		openFile, err := excelize.OpenFile(filename)
		if err != nil {
			log.Fatalln("Unable to open excel file ", err)
			return f, err
		}
		f.f = openFile
	}

	return f, nil
}

func (f File) Sheet() string {
	list := f.f.GetSheetList()
	if len(list) < 1 {
		log.Println("Empty sheet list")
		return ""
	}

	return list[0]
}

func (f File) GetRows() ([][]string, error) {
	tables := make([][]string, 0)
	sheet := f.Sheet()
	log.Println("The name of the sheet here", len(sheet))
	tables, err := f.f.GetRows(sheet)
	if err != nil {
		return tables, nil
	}

	return tables, nil
}

type DB struct {
	conn *sql.Conn
	db   *sql.DB
	f    *File
}

func NewDB(ctx context.Context, driver, url string, f *File) *DB {
	if url == "" {
		return nil
	}

	db, err := sql.Open(driver, url)
	if err != nil {
		log.Println("Error setting up database")
		panic(err)
	}

	dbConn, err := db.Conn(ctx)
	if err != nil {
		panic(err)
	}

	return &DB{
		conn: dbConn,
		db:   db,
		f:    f,
	}
}

func (d DB) CreateTable(ctx context.Context, tableCols map[string]interface{}, tablename string) error {
	if len(tableCols) < 1 {
		return errors.New("Invalid table column length")
	}

	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", tablename)
	s, err := d.conn.PrepareContext(ctx, stmt)
	if err != nil {
		return fmt.Errorf("Invalid query: %v\n", err)
	}

	res, err := s.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("Error executing query: %v\n", err)
	}

	val, _ := res.RowsAffected()
	if val > 0 {
		log.Printf("%s successfully created!", tablename)
	}
	return nil
}

func (d DB) Ping(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	d.conn.PingContext(ctx)
}

func (d DB) Import(tableValues [][]string, contentIdx int) error {
	//Check each row one after the other
	//Infer the value type
	for idx, row := range tableValues {
		if idx >= contentIdx {
			for _, v := range row {
				fmt.Printf("The rows %s values type %T", idx, v)
				return nil
			}
		}
	}
	return nil
}
