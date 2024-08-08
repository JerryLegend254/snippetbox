package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JerryLegend254/snippetbox/pkg/models"
	"github.com/go-chi/chi/v5"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w, err)
		return
	}

	a.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})

}

func (a *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		a.notFound(w)
		return
	}
	s, err := a.snippets.Get(id)
	if err == models.ErrNoRecord {
		a.notFound(w)
		return
	} else if err != nil {
		a.serverError(w, err)
		return
	}

	a.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})

}

func (a *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	id, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.serverError(w, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}
