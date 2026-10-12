package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./assets/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /", home)
	mux.HandleFunc("GET /games", gamesView)
	mux.HandleFunc("GET /games/{id}", gameView)
	mux.HandleFunc("GET /games/new", gameCreate)
	mux.HandleFunc("POST /games/new", gameCreatePost)

	log.Printf("starting server on %s", *addr)

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}