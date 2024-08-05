package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JerryLegend254/snippetbox/pkg/models"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFound(w)
		return
	}
	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.
	//files := []string{
	//	"./ui/html/home.page.tmpl",
	//	"./ui/html/base.layout.tmpl",
	//	"./ui/html/footer.partial.tmpl",
	//}
	//// Use the template.ParseFiles() function to read the files and store the
	//// templates in a template set. Notice that we can pass the slice of file p
	//// as a variadic parameter?
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	a.errorLog.Println(err.Error())
	//	a.serverError(w, err)
	//	return
	//}
	//err = ts.Execute(w, nil)
	//if err != nil {
	//	a.errorLog.Println(err.Error())
	//	http.Error(w, "Internal Server Error", 500)
	//}
	s, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w, err)
		return
	}
	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}
}

func (a *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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
	}

	fmt.Fprintf(w, "%v", s)
}

func (a *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		a.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := "7"
	id, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.serverError(w, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
