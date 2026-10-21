package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"bgtracker.ndr54.com/internal/models"
	"bgtracker.ndr54.com/internal/validator"
)

type gameCreateForm struct {
	Title					string `form:"title"`
	Bgg_id					int `form:"bgg_id"`
	Year_pub				int `form:"year_pub"`
	Description				string	`form:"description"`
	Price					float64 `form:"price"`
	Player_min				int `form:"player_min"`
	Player_max				int `form:"player_max"`
	Playtime_min			int `form:"playtime_min"`
	Playtime_max			int `form:"playtime_max"`
	Age_min					int `form:"age_min"`
	Scoring					bool `form:"scoring"`
	State					string `form:"state"`
	validator.Validator		`form:"-"`
}

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

	data.Form = gameCreateForm{
		Year_pub:	time.Now().Year(),
		Scoring: 	true,
		State:		"owned",
	}

	err := app.render(w, r, http.StatusOK, "addGame.tmpl.html", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) gameCreatePost(w http.ResponseWriter, r *http.Request) {
	var form gameCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// err := app.formDecoder.Decode(&form, r.PostForm)
	// if err !=  nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	// bgg_id, err := strconv.Atoi(r.PostForm.Get("bgg_id"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
    //     return
	// }

	// year_pub, err := strconv.Atoi(r.PostForm.Get("year_pub"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
    //     return
	// }

	// price, err := strconv.ParseFloat(r.PostForm.Get("price"), 64)
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
    //     return
	// }

	// player_min, err := strconv.Atoi(r.PostForm.Get("player_min"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
    //     return
	// }

	// player_max, err := strconv.Atoi(r.PostForm.Get("player_max"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
    //     return
	// }

	// playtime_min, err := strconv.Atoi(r.PostForm.Get("playtime_min"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
    //     return
	// }

	// playtime_max, err := strconv.Atoi(r.PostForm.Get("playtime_max"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
    //     return
	// }

	// age_min, err := strconv.Atoi(r.PostForm.Get("age_min"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
    //     return
	// }

	// // Defaults to 'no/false' if it returns blank as well
	// scoring := false
	// if r.PostForm.Get("scoring") == "yes" {
	// 	scoring = true
	// }

	// form := gameCreateForm{
	// 	Title:			r.PostForm.Get("title"),
	// 	Bgg_id:			bgg_id,
	// 	Year_pub:		year_pub,
	// 	Description:	r.PostForm.Get("description"),
	// 	Price:			price,
	// 	Player_min:		player_min,
	// 	Player_max:		player_max,
	// 	Playtime_min:	playtime_min,
	// 	Playtime_max:	playtime_max,
	// 	Age_min:		age_min,
	// 	Scoring:		scoring,
	// 	State:			r.PostForm.Get("state"),
	// }

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be longer than 100 characters")
	form.CheckField(validator.PermittedValue(form.State, "owned", "wanted"), "state", "This field must be either owned or wanted")
	// Add check for scoring?
	// Add sensible checks for other fields (published player and time min < max)


	if !form.Valid() {
		data  := app.newTemplateData(r)
		data.Form = form
		
		err := app.render (w, r, http.StatusUnprocessableEntity, "addGame.tmpl.html", data)
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	id, err := app.games.Insert(form.Title, form.Bgg_id, form.Price,
		form.State, form.Year_pub, form.Description, form.Player_min,
		form.Player_max, form.Playtime_min, form.Playtime_max, form.Age_min,
		form.Scoring)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/games/%d", id), http.StatusSeeOther)
}