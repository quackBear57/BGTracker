package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Write([]byte("Hello from BG Tracker"))
}

func gamesView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Page to view all of the games..."))
}

func gameView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w,r)
		return
	}
	
	fmt.Fprintf(w, "Viewing game ID %d details.", id)
}

func gameCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Page to create a new game..."))
}

func gameCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save the new game..."))
}