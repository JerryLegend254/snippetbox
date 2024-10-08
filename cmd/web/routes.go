package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := chi.NewRouter()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home).(http.HandlerFunc))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnippet).(http.HandlerFunc))
	mux.Get("/snippet/{id}", dynamicMiddleware.ThenFunc(app.showSnippet).(http.HandlerFunc))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnippetForm).(http.HandlerFunc))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm).(http.HandlerFunc))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser).(http.HandlerFunc))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm).(http.HandlerFunc))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser).(http.HandlerFunc))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.logoutUser).(http.HandlerFunc))

	mux.Get("/ping", http.HandlerFunc(ping))

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)

}
