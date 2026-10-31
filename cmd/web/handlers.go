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

type userChangePasswordForm struct {
	CurrentPassword			string `form:"currentPassword"`
	NewPassword				string `form:"newPassword"`
	NewPasswordConfirm		string `form:"newPasswordConfirm"`
	validator.Validator		`form:"-"`
}

type playerCreateForm struct {
	Name					string `form:"name"`
	UserID					int `form:"userID"`
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

	err = app.render(w, r, http.StatusOK, "gameView.tmpl.html", data)
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

	path := app.sessionManager.PopString(r.Context(), "redirectPathAfterLogin")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/games/new/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash",
		"You've been logged out successfully")
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	if id < 1 {
		http.NotFound(w,r)
		return
	}

	user, err := app.users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoRecord):
			http.NotFound(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}
	
	data := app.newTemplateData(r)
	data.User = user

	err = app.render(w, r, http.StatusOK, "accountView.tmpl.html", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) changePass(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userChangePasswordForm{}
	
	err := app.render(w, r, http.StatusOK, "changePassword.tmpl.html", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) changePassPost(w http.ResponseWriter, r *http.Request) {
	var form userChangePasswordForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword",
		"This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword",
		"This field cannot be blank")
		form.CheckField(validator.NotBlank(form.NewPasswordConfirm), "newPasswordConfirm",
		"This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword",
		"This field must be at least 8 characters long")
	form.CheckField(validator.MinChars(form.NewPasswordConfirm, 8), "newPasswordConfirm",
		"This field must be at least 8 characters long")
	form.CheckField(validator.MaxBytes(form.NewPassword, 72), "newPassword",
		"This field must not be more than 72 bytes long")
	form.CheckField(validator.MaxBytes(form.NewPasswordConfirm, 72), "newPasswordConfirm",
		"This field must not be more than 72 bytes long")
	form.CheckField(validator.StringMatch(form.NewPassword, form.NewPasswordConfirm),
		"newPassword", "New passwords must match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		
		err := app.render(w, r, http.StatusUnprocessableEntity,
			"changePassword.tmpl.html", data)
		if err != nil {
			app.serverError(w,r, err)
		}
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	user, err := app.users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoRecord):
			http.NotFound(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	id, err = app.users.Authenticate(user.Email, form.CurrentPassword)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidCredentials):
			form.AddNonFieldError("Password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form

			err := app.render(w, r, http.StatusUnprocessableEntity,
				"changePassword.tmpl.html", data)
			if err != nil {
				app.serverError(w, r, err)
			}
		default:
			app.serverError(w, r, err)
		}
		return
	}

	ok, err := app.users.ChangePassword(id, form.CurrentPassword, form.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidCredentials):
			form.AddNonFieldError("Password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form

			err := app.render(w, r, http.StatusUnprocessableEntity,
				"changePassword.tmpl.html", data)
			if err != nil {
				app.serverError(w, r, err)
			}
		case errors.Is(err, models.ErrNoUser):
			form.AddNonFieldError("Could not find user")
			data := app.newTemplateData(r)
			data.Form = form

			err := app.render(w, r, http.StatusUnprocessableEntity,
				"changePassword.tmpl.html", data)
			if err != nil {
				app.serverError(w, r, err)
			}
		default:
			app.serverError(w, r, err)
		}
	}

	if !ok {
		form.AddNonFieldError("Unknown error")
		data := app.newTemplateData(r)
		data.Form = form
	
		err := app.render(w, r, http.StatusUnprocessableEntity,
			"changePassword.tmpl.html", data)
		if err != nil {
			app.serverError(w, r, err)
		}
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash",
		"Password changed successfully!")

	app.sessionManager.Put(r.Context(), "flash", "Password changed successfully!")
	
	http.Redirect(w, r, "/user/account", http.StatusSeeOther)

	
}

func (app *application) playerCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = playerCreateForm{}

	err := app.render(w, r, http.StatusOK, "addPlayer.tmpl.html", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) playerCreatePost(w http.ResponseWriter, r *http.Request) {
	var form playerCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	// Add check for scoring?
	// Add sensible checks for other fields (published player and time min < max)

	if form.UserID != 0 {
		user, err := app.users.Get(form.UserID)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrNoRecord):
				form.Validator.AddFieldError("userID", "User not found")

				data := app.newTemplateData(r)
				data.Form = form

				err := app.render(w, r, http.StatusUnprocessableEntity,
					"addPlayer.tmpl.html", data)
				if err != nil {
					app.serverError(w, r, err)
				}
		return
			default:
				app.serverError(w, r, err)
			}
		return
		}

		// This isn't properly checking!!!
		// confirmed the right user id is being entered into this function

		taken, err := app.players.UserIDTaken(user.ID)

		if err != nil {
			if taken {
				form.Validator.AddFieldError("userID", "User ID already in use")

				data := app.newTemplateData(r)
				data.Form = form

				err := app.render(w, r, http.StatusUnprocessableEntity,
					"addPlayer.tmpl.html", data)
				if err != nil {
					app.serverError(w, r, err)
				}
			}
		}
	}

	if !form.Valid() {
		form.UserID = 0
		data  := app.newTemplateData(r)
		data.Form = form
		
		err := app.render (w, r, http.StatusUnprocessableEntity, "addPlayer.tmpl.html", data)
		if err != nil {
			app.serverError(w, r, err)
		}
		return
	}


	// this currently fails when user id = 0 since it's set as a foreign key
	// 0 doesn't match a user in users, so errors out
	id, err := app.players.Insert(form.Name, form.UserID)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Player added successfully!")

	http.Redirect(w, r, fmt.Sprintf("/players/%d", id), http.StatusSeeOther)
}

func (app *application) playerView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w,r)
		return
	}

	player, err := app.players.Get(id)
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
	data.Player = player

	err = app.render(w, r, http.StatusOK, "playerView.tmpl.html", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) playersView(w http.ResponseWriter, r *http.Request) {
	players, err := app.players.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Players = players

	err = app.render(w, r, http.StatusOK, "playersView.tmpl.html", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}