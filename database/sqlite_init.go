package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Init() *sql.DB {
	db, err := sql.Open("sqlite3", "./demo.db")
	if err != nil {
		log.Fatal(err)
	}

	//DEMO PURPOSES: creating two tables, "urls" and "stats" and index
	sqlStmt := `
	CREATE TABLE urls (short_url text NOT NULL PRIMARY KEY, long_url text NOT NULL, expire_at timestamp NOT NULL);
	CREATE TABLE stats (id INTEGER PRIMARY KEY AUTOINCREMENT, short_url text NOT NULL, created_at timestamp NOT NULL);
	CREATE INDEX idx_short_url_created_at ON stats(short_url, created_at)
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Print(err)
	}

	return db
}
