package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./assets/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /", dynamic.ThenFunc(app.home))
	mux.Handle("GET /games", dynamic.ThenFunc(app.gamesView))
	mux.Handle("GET /games/{id}", dynamic.ThenFunc(app.gameView))
	mux.Handle("GET /games/new", dynamic.ThenFunc(app.gameCreate))
	mux.Handle("POST /games/new", dynamic.ThenFunc(app.gameCreatePost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}