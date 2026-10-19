package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

func (app *application) gameCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	// fmt.Printf("title: %s \n", title)

	bgg_id, err := strconv.Atoi(r.PostForm.Get("bgg_id"))
	if err != nil {
		app.logger.Error(err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// fmt.Printf("bgg_id: %d \n", bgg_id)

	year_pub, err := strconv.Atoi(r.PostForm.Get("year_pub"))
	if err != nil {
		app.logger.Error(err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// fmt.Printf("year_pub: %d \n", year_pub)

	description := r.PostForm.Get("description")
	// fmt.Printf("description: %s \n", description)

	price, err := strconv.ParseFloat(r.PostForm.Get("price"), 64)
	if err != nil {
		app.logger.Error(err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// fmt.Printf("price: %f \n", price)

	player_min, err := strconv.Atoi(r.PostForm.Get("player_min"))
	if err != nil {
		app.logger.Error(err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// fmt.Printf("player_min: %d \n", player_min)

	player_max, err := strconv.Atoi(r.PostForm.Get("player_max"))
	if err != nil {
		app.logger.Error(err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// fmt.Printf("player_max: %d \n", player_max)

	playtime_min, err := strconv.Atoi(r.PostForm.Get("playtime_min"))
	if err != nil {
		app.logger.Error(err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// fmt.Printf("playtime_min: %d \n", playtime_min)

	playtime_max, err := strconv.Atoi(r.PostForm.Get("playtime_max"))
	if err != nil {
		app.logger.Error(err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// fmt.Printf("playtime_max: %d \n", playtime_max)

	age_min, err := strconv.Atoi(r.PostForm.Get("age_min"))
	if err != nil {
		app.logger.Error(err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// fmt.Printf("age_min: %d \n", age_min)

	score_int, err := strconv.Atoi(r.PostForm.Get("scoring"))
	if err != nil {
		app.logger.Error(err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	scoring := true
	if score_int == 1 {
		scoring = true
	} else if score_int == 0 {
		scoring = false
	} else {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// fmt.Printf("scoring: %t \n", scoring)

	state := r.PostForm.Get("state")
	// fmt.Printf("state: %s \n", state)
	
	id, err := app.games.Insert(title, bgg_id, price, state, year_pub,
	description, player_min, player_max, playtime_min, playtime_max,
	age_min, scoring)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/games/%d", id), http.StatusSeeOther)
}