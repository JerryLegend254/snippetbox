package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := chi.NewRouter()

	mux.Get("/", app.home)
	mux.Post("/snippet/create", app.createSnippet)
	mux.Get("/snippet/{id}", app.showSnippet)
	mux.Get("/snippet/create", app.createSnippetForm)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)

}
