package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"mmaschedule-go/event"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

var host = flag.String("host", "127.0.0.1:8000", "Set the web server host address.")
var debug = flag.Bool("debug", false, "Enable debug mode.")

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

	if *debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
        slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	f, err := os.OpenFile("mmaschedule.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	defer f.Close()
	log.SetOutput(f)

	cmd := flag.Args()[0]

	client := NewScraperClient()
	queries, err := InitDb("db.sqlite")
	if err != nil {
		slog.Error("Error initializing database", "error", err)
		return
	}

	switch cmd {
	case "runserver":
		go ScrapeEventsLoop(queries, client, true)
		err = RunServer(*host, queries)
		if err != nil {
			slog.Error("Error starting web server:", "error", err)
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
