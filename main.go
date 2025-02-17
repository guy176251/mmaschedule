package main

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"mmaschedule-go/event"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	q, err := InitDb("db.sqlite")
	if err != nil {
		fmt.Println("Error initializing database: ", err)
	} else {
		client := NewScraperClient()
		ScrapeEvents(q, client)
	}
}

func InitDb(name string) (*event.Queries, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return nil, err
	}

	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return nil, err
	}

	queries := event.New(db)

	return queries, nil
}
