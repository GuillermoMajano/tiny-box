package main

import (
	"GuillermoMajano/snippetbox/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	s, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}
	tp := &TemplateData{Snippets: s}

	app.render(w, r, "home.page.tmpl", tp)

	/*data := &TemplateData{Snippets: s}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, err)

		http.Error(w, "Internal Server Error", 404)
		return
	}

	err = ts.Execute(w, data)

	if err != nil {
		app.serverError(w, err)
		return
	}*/

}

func (app application) showSnippet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("you request to showsnippext")
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest()

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return

	}

	app.render(w, r, "show.page.tmpl", &TemplateData{Snippets: s})

	/*
		ts, err := template.ParseFiles(files...)

		if err != nil {
			app.serverError(w, err)
			return
		}

		err = ts.Execute(w, data)

		if err != nil {
			app.serverError(w, err)
			return
		}*/
}

func (app application) createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "0 snail"
	content := "0 snail\nClim Mount Fuji,\nBut slowly,slowly!\n\n- kobatashi issa"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

	w.Write([]byte("Create a new snippet..."))
}

/**func downloadHandler(w http.ResponseWriter, r *http.Request) {
	//fp := filepath.Clean("./ui/static/")
	http.FileServer(http.Dir("./ui/static/"))
}**/

func (app *application) showallT(w http.ResponseWriter, R *http.Request) {
	getids, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	jresp, err := json.Marshal(getids)

	if err != nil {
		app.errorlog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jresp)

}
