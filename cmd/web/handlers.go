package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	files := []string{
		"./assets/html/pages/home.tmpl.html",
		"./assets/html/partials/nav.tmpl.html",
		"./assets/html/pages/base.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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