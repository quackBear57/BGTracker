package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", home)
	mux.HandleFunc("GET /games", gamesView)
	mux.HandleFunc("GET /games/{id}", gameView)
	mux.HandleFunc("GET /games/new", gameCreate)
	mux.HandleFunc("POST /games/new", gameCreatePost)

	log.Print("starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}