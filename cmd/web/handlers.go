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

func (app *application) gameCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	fieldErrors := make(map[string]string)

	title := r.PostForm.Get("title")
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be longer than 100 characters"
	}

	bgg_idValue := r.PostForm.Get("bgg_id")
	bgg_id := 0
	if bgg_idValue != "" {
		bgg_id, err = strconv.Atoi(bgg_idValue)
			if err != nil {
			fieldErrors["bgg_id"] = "This field must be an integer"
		}
	}
	
	// bgg_id, err := strconv.Atoi(r.PostForm.Get("bgg_id"))
	// if err != nil {
	// 	fieldErrors["bgg_id"] = "This field must be an integer"
	// }

	year_pub, err := strconv.Atoi(r.PostForm.Get("year_pub"))
	if err != nil {
		fieldErrors["year_pub"] = "This field must be an integer"
	}

	description := r.PostForm.Get("description")

	price, err := strconv.ParseFloat(r.PostForm.Get("price"), 64)
	if err != nil {
		fieldErrors["price"] = "This field must be a number"
	}

	player_min, err := strconv.Atoi(r.PostForm.Get("player_min"))
	if err != nil {
		fieldErrors["player_min"] = "This field must be an integer"
	}

	player_max, err := strconv.Atoi(r.PostForm.Get("player_max"))
	if err != nil {
		fieldErrors["player_max"] = "This field must be an integer"
	}

	playtime_min, err := strconv.Atoi(r.PostForm.Get("playtime_min"))
	if err != nil {
		fieldErrors["playtime_min"] = "This field must be an integer"
	}

	playtime_max, err := strconv.Atoi(r.PostForm.Get("playtime_max"))
	if err != nil {
		fieldErrors["playtime_max"] = "This field must be an integer"
	}

	age_min, err := strconv.Atoi(r.PostForm.Get("age_min"))
	if err != nil {
		fieldErrors["age_min"] = "This field must be an integer"
	}

	scoringValue := r.PostForm.Get("scoring")
	// fmt.Printf("scoring value: %s", scoringValue)
	scoreSet := map[string]bool{"yes": true, "no": true}
	if !scoreSet[scoringValue] {
		fieldErrors["scoring"] = "This field must either be Yes or No"
	}
	scoring := true
	if scoringValue == "yes" {
		scoring = true
	} else if scoringValue == "no" {
		scoring = false
	}

	// score_int, err := strconv.Atoi(r.PostForm.Get("scoring"))
	// if err != nil {
	// 	app.logger.Error(err.Error())
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }
	// scoring := true
	// if score_int == 1 {
	// 	scoring = true
	// } else if score_int == 0 {
	// 	scoring = false
	// } else {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	stateSet := map[string]bool{"owned": true, "wanted": true}
	state := r.PostForm.Get("state")
	if !stateSet[state] {
		fieldErrors["state"] = "This field must be either owned or wanted"
	}
	
	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}

	id, err := app.games.Insert(title, bgg_id, price, state, year_pub,
	description, player_min, player_max, playtime_min, playtime_max,
	age_min, scoring)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/games/%d", id), http.StatusSeeOther)
}