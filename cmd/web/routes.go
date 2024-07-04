package main

import (
	"net/http"

	"Service/ui"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.PathPrefix("/static/").Handler(fileServer).Methods(http.MethodGet)
	//router.Handle("/static/", fileServer).Methods(http.MethodGet)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handle("/users/signup", dynamic.ThenFunc(app.userSignUp)).Methods(http.MethodGet)
	router.Handle("/users/signup", dynamic.ThenFunc(app.userSignUpPost)).Methods(http.MethodPost)
	router.Handle("/users/login", dynamic.ThenFunc(app.userLogin)).Methods(http.MethodGet)
	router.Handle("/users/login", dynamic.ThenFunc(app.userLoginPost)).Methods(http.MethodPost)

	// Protected (authenticated-only) application routes
	protected := dynamic.Append(app.requireAuthentication)

	router.Handle("/days/view/{id}", protected.ThenFunc(app.dayView)).Methods(http.MethodGet)
	router.Handle("/", protected.ThenFunc(app.home)).Methods(http.MethodGet)
	router.Handle("/days/create", protected.ThenFunc(app.dayCreate)).Methods(http.MethodGet)
	router.Handle("/days/create", protected.ThenFunc(app.dayCreatePost)).Methods(http.MethodPost)
	router.Handle("/users/logout", protected.ThenFunc(app.userLogoutPost)).Methods(http.MethodPost)
	router.Handle("/users/account", protected.ThenFunc(app.accountView)).Methods(http.MethodGet)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
