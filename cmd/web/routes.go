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
	mux.Handle("GET /players/{id}", dynamic.ThenFunc(app.playerView))
	mux.Handle("GET /sessions", dynamic.ThenFunc(app.sessionsView))
	mux.Handle("GET /sessions/{id}", dynamic.ThenFunc(app.sessionView))

	// for protected routes (require autentication)
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /games/new", protected.ThenFunc(app.gameCreate))
	mux.Handle("POST /games/new", protected.ThenFunc(app.gameCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	mux.Handle("GET /user/account", protected.ThenFunc(app.accountView))
	mux.Handle("GET /user/account/change-password", protected.ThenFunc(app.changePass))
	mux.Handle("POST /user/account/change-password", protected.ThenFunc(app.changePassPost))
	mux.Handle("GET /user/account/change-name", protected.ThenFunc(app.userChangeName))
	mux.Handle("POST /user/account/change-name", protected.ThenFunc(app.userChangeNamePost))
	mux.Handle("GET /players/new", protected.ThenFunc(app.playerCreate))
	mux.Handle("POST /players/new", protected.ThenFunc(app.playerCreatePost))
	mux.Handle("GET /players/assign", protected.ThenFunc(app.playerAssignUser))
	mux.Handle("POST /players/assign", protected.ThenFunc(app.playerAssignUserPost))
	mux.Handle("GET /players/change-name/{id}", protected.ThenFunc(app.playerChangeName))
	mux.Handle("POST /players/change-name/{id}", protected.ThenFunc(app.playerChangeNamePost))
	// mux.Handle("GET /sessions/new", protected.ThenFunc(app.sessionCreate))
	// mux.Handle("POST /sessions/new", protected.ThenFunc(app.sessionCreatePost))
	mux.Handle("GET /sessions/new/{id}", protected.ThenFunc(app.sessionCreate))
	mux.Handle("POST /sessions/new/{id}", protected.ThenFunc(app.sessionCreatePost))
	mux.Handle("GET /games/{id}/scoresheet", protected.ThenFunc(app.gameScoresheet))
	mux.Handle("POST /games/{id}/scoresheet", protected.ThenFunc(app.gameScoresheetPost))
	mux.Handle("GET /games/{id}/scoresheet/edit", protected.ThenFunc(app.gameScoresheetEdit))
	mux.Handle("POST /games/{id}/scoresheet/edit", protected.ThenFunc(app.gameScoresheetEditPost))


	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}