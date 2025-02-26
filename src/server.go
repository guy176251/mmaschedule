package main

import (
	"embed"
	"fmt"
	"io/fs"
	"mmaschedule-go/event"
	"net/http"
	"os"
	"strconv"
	"time"
)

//go:embed static/*
var staticfiles embed.FS

func RunServer(addr string, queries *event.Queries) error {
	mux := http.NewServeMux()
	state := ServerState{queries: queries}
	static, err := StaticHandler("/static/")
	if err != nil {
		return err
	}
	mux.HandleFunc("GET /", state.HandleFunc(RouteIndex))
	mux.HandleFunc("GET /events/{slug}/", state.HandleFunc(RouteEvent))
	mux.Handle("GET /static/", static)
	err = http.ListenAndServe(addr, mux)
	return err
}

func StaticHandler(prefix string) (http.Handler, error) {
	static, err := fs.Sub(staticfiles, "static")
	if err != nil {
		return nil, err
	}
	return http.StripPrefix(prefix, http.FileServer(http.FS(static))), nil
}

func RouteIndex(w http.ResponseWriter, r *http.Request, q *event.Queries) {
	slug, err := q.GetUpcomingEvent(r.Context(), EventTime())
	if err != nil {
		http.NotFound(w, r)
	} else {
		http.Redirect(w, r, "/events/"+slug+"/", http.StatusSeeOther)
	}
}

func RouteEvent(w http.ResponseWriter, r *http.Request, q *event.Queries) {
	slug := r.PathValue("slug")

	if len(slug) == 0 {
		http.NotFound(w, r)
		return
	}

	event, err := q.GetEvent(r.Context(), slug)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	upcoming, err := q.ListUpcomingEvents(r.Context(), EventTime())
	if err != nil {
		http.NotFound(w, r)
		return
	}

	TemplEventPage(event, upcoming).Render(r.Context(), w)
}

func EventTime() int64 {
	timestamp, err := EnvTime()
	if err != nil {
		t := time.Now()
		return t.Add(time.Duration(-6) * time.Hour).Unix()
	} else {
		return timestamp
	}
}

func EnvTime() (int64, error) {
	timestr := os.Getenv("CURRENT_TIME")
	if len(timestr) == 0 {
		return 0, fmt.Errorf("Time env is empty.")
	}
	return strconv.ParseInt(timestr, 10, 64)
}

type StateHandler func(w http.ResponseWriter, r *http.Request, q *event.Queries)

type ServerState struct {
	queries *event.Queries
}

func (state *ServerState) HandleFunc(handler StateHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, state.queries)
	}
}
