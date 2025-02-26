package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"mmaschedule-go/event"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

var host = flag.String("host", "127.0.0.1:8000", "Set the web server host address.")

const POSITIONAL_ARGS_HELP = `
Valid commands:
    runserver: Runs the web server.
    scrape: Runs the web scraper.
`

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Print(POSITIONAL_ARGS_HELP)
		os.Exit(1)
	}

	cmd := flag.Args()[0]

	client := NewScraperClient()
	queries, err := InitDb("db.sqlite")
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	switch cmd {
	case "runserver":
		go ScrapeEventsLoop(queries, client, true)
		err = RunServer(*host, queries)
		if err != nil {
			fmt.Println("Error starting web server:", err)
		}
	case "scrape":
		ScrapeEvents(queries, client, true)
	default:
		fmt.Print(POSITIONAL_ARGS_HELP)
		os.Exit(1)
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
