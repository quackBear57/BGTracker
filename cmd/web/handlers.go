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

type userSignupForm struct {
	Name					string `form:"name"`
	Email					string `form:"email"`
	Password				string `form:"password"`
	validator.Validator		`form:"-"`
}

type userLoginForm struct {
	Email					string `form:"email"`
	Password				string `form:"password"`
	validator.Validator		`form:"-"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

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

	app.sessionManager.Put(r.Context(), "flash", "Game added successfully!")

	http.Redirect(w, r, fmt.Sprintf("/games/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}

	err := app.render(w, r, http.StatusOK, "signup.tmpl.html", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name",
		"This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email",
		"This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email",
		"This field must be a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password",
		"This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password",
		"This field must be at least 8 characters long")
	form.CheckField(validator.MaxBytes(form.Password, 72), "password",
		"This field must not be more than 72 bytes long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		
		err := app.render(w, r, http.StatusUnprocessableEntity,
			"signup.tmpl.html", data)
		if err != nil {
			app.serverError(w,r, err)
		}
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			form.AddFieldError("email", "Email address already in use")
			data := app.newTemplateData(r)
			data.Form = form

			err = app.render(w, r, http.StatusUnprocessableEntity,
				"signup.tmpl.html", data)
			if err != nil {
				app.serverError(w, r, err)
			}
		default:
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash",
		"Your signup was successful. Please log in.")
	
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}

	err := app.render(w, r, http.StatusOK, "login.tmpl.html", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email",
		"This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email",
		"This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password",
		"This field cannot be blank")
	form.CheckField(validator.MaxBytes(form.Password, 72), "password",
		"This field must not be more than 72 bytes long")
	
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		err := app.render(w, r, http.StatusUnprocessableEntity,
			"login.tmpl.html", data)
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidCredentials):
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form

			err := app.render(w, r, http.StatusUnprocessableEntity,
				"login.tmpl.html", data)
			if err != nil {
				app.serverError(w, r, err)
			}
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/games/new/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout")
}