package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable)

	mux := chi.NewRouter()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home).(http.HandlerFunc))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet).(http.HandlerFunc))
	mux.Get("/snippet/{id}", dynamicMiddleware.ThenFunc(app.showSnippet).(http.HandlerFunc))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm).(http.HandlerFunc))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm).(http.HandlerFunc))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser).(http.HandlerFunc))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm).(http.HandlerFunc))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser).(http.HandlerFunc))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser).(http.HandlerFunc))

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)

}
