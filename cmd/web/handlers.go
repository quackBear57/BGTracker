package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"bgtracker.ndr54.com/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	games, err := app.games.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Games = games

	err = app.render(w, r, http.StatusOK, "home.tmpl.html", data)
	if err != nil {
		app.serverError(w, r, err)
	}
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

	data := app.newTemplateData(r)
	data.Game = game

	err = app.render(w, r, http.StatusOK, "view.tmpl.html", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) gameCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := app.render(w, r, http.StatusOK, "addGame.tmpl.html", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

type gameCreateForm struct {
	Title			string
	Bgg_id			int
	Year_pub		int
	Description		string
	Price			float64
	Player_min		int
	Player_max		int
	Playtime_min	int
	Playtime_max	int
	Age_min			int
	Scoring			bool
	State			string
	FieldErrors		map[string]string
}

func (app *application) gameCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	bgg_id, err := strconv.Atoi(r.PostForm.Get("bgg_id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
        return
	}

	year_pub, err := strconv.Atoi(r.PostForm.Get("year_pub"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
        return
	}

	price, err := strconv.ParseFloat(r.PostForm.Get("price"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
        return
	}

	player_min, err := strconv.Atoi(r.PostForm.Get("player_min"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
        return
	}

	player_max, err := strconv.Atoi(r.PostForm.Get("player_max"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
        return
	}

	playtime_min, err := strconv.Atoi(r.PostForm.Get("playtime_min"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
        return
	}

	playtime_max, err := strconv.Atoi(r.PostForm.Get("playtime_max"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
        return
	}

	age_min, err := strconv.Atoi(r.PostForm.Get("age_min"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
        return
	}

	// Defaults to 'no/false' if it returns blank as well
	scoring := false
	if r.PostForm.Get("scoring") == "yes" {
		scoring = true
	}

	form := gameCreateForm{
		Title:			r.PostForm.Get("title"),
		Bgg_id:			bgg_id,
		Year_pub:		year_pub,
		Description:	r.PostForm.Get("description"),
		Price:			price,
		Player_min:		player_min,
		Player_max:		player_max,
		Playtime_min:	playtime_min,
		Playtime_max:	playtime_max,
		Age_min:		age_min,
		Scoring:		scoring,
		State:			r.PostForm.Get("state"),
		FieldErrors:	map[string]string{},
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be longer than 100 characters"
	}

	if strings.TrimSpace(form.State) == "" {
		form.FieldErrors["state"] = "This field must not be blank"
	}
	
	if len(form.FieldErrors) > 0 {
		data  := app.newTemplateData(r)
		data.Form = form
		
		err := app.render (w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	id, err := app.games.Insert(form.Title, bgg_id, price, form.State, year_pub,
	form.Description, player_min, player_max, playtime_min, playtime_max,
	age_min, scoring)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/games/%d", id), http.StatusSeeOther)
}