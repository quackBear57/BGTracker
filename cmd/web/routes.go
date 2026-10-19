package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./assets/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("GET /games", app.gamesView)
	mux.HandleFunc("GET /games/{id}", app.gameView)
	mux.HandleFunc("GET /games/new", app.gameCreate)
	mux.HandleFunc("POST /games/new", app.gameCreatePost)

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}