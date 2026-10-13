package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./assets/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("GET /games", app.gamesView)
	mux.HandleFunc("GET /games/{id}", app.gameView)
	mux.HandleFunc("GET /games/new", app.gameCreate)
	mux.HandleFunc("POST /games/new", app.gameCreatePost)

	return mux
}