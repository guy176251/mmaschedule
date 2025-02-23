package main

import (
	"mmaschedule-go/event"
	"net/http"
	"time"
)

func RunServer(addr string, queries *event.Queries) error {
	mux := http.NewServeMux()
	state := ServerState{queries: queries}
	mux.HandleFunc("GET /", state.HandleFunc(RouteIndex))
	mux.HandleFunc("GET /events/{slug}/", state.HandleFunc(RouteEvent))
	err := http.ListenAndServe(addr, mux)
	return err
}

func RouteIndex(w http.ResponseWriter, r *http.Request, q *event.Queries) {
	slug, err := q.GetUpcomingEvent(r.Context(), time.Now().Unix())
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
		http.NotFound(w, r)
		return
	}

	upcoming, err := q.ListUpcomingEvents(r.Context(), time.Now().Unix())
	if err != nil {
		http.NotFound(w, r)
		return
	}

	TemplEventPage(event, upcoming).Render(r.Context(), w)
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
