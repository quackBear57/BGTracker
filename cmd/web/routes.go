package main

import (
	"net/http"

	"bgtracker.ndr54.com/assets"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(assets.Files))

	// for unprotected routes (non-authenticated users)
	dynamic := alice.New(app.sessionManager.LoadAndSave, preventCSRF,
		app.authenticate)

	mux.Handle("GET /", dynamic.ThenFunc(app.home))
	mux.Handle("GET /games", dynamic.ThenFunc(app.gamesView))
	mux.Handle("GET /games/{id}", dynamic.ThenFunc(app.gameView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("GET /players", dynamic.ThenFunc(app.playersView))

	// for protected routes (require autentication)
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /games/new", protected.ThenFunc(app.gameCreate))
	mux.Handle("POST /games/new", protected.ThenFunc(app.gameCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	mux.Handle("GET /user/account", protected.ThenFunc(app.accountView))
	mux.Handle("GET /user/account/change-password", protected.ThenFunc(app.changePass))
	mux.Handle("POST /user/account/change-password", protected.ThenFunc(app.changePassPost))
	mux.Handle("GET /players/new", protected.ThenFunc(app.playerCreate))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}