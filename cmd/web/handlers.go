package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"bgtracker.ndr54.com/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	games, err := app.games.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, game := range games {
		fmt.Fprintf(w, "%+v\n", game)
	}

	// files := []string{
	// 	"./assets/html/pages/home.tmpl.html",
	// 	"./assets/html/partials/nav.tmpl.html",
	// 	"./assets/html/base.tmpl.html",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// }
}

func (app *application) gamesView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Page to view all of the games..."))
}

func (app *application) gameView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w,r)
		return
	}

	game, err := app.games.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoRecord):
			http.NotFound(w,r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	files := []string{
		"./assets/html/base.tmpl.html",
		"./assets/html/partials/nav.tmpl.html",
		"./assets/html/pages/view.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", game)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) gameCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Page to create a new game..."))
}

func (app *application) gameCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "Test Game 1"
	bgg_id := 12345
	price := 20.99
	state := "owned"
	year_pub := 1999
	description := "test desc"
	player_min := 1
	player_max := 99
	playtime_min := 10
	playtime_max := 120
	age_min := 10
	scoring := true

	id, err := app.games.Insert(title, bgg_id, price, state, year_pub,
	description, player_min, player_max, playtime_min, playtime_max,
	age_min, scoring)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/games/%d", id), http.StatusSeeOther)
}