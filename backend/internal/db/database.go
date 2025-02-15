package db

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"os"
)

type Database struct {
	Name string
	DB   *sql.DB
}

func (d *Database) Open() (err error) {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: d.Name,
	}

	d.DB, err = sql.Open("mysql", cfg.FormatDSN())
	return err
}

func (d *Database) Close() (err error) {
	return d.DB.Close()
}
