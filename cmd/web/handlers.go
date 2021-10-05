package main

import (
	"GuillermoMajano/snippetbox/pkg/forms"
	"GuillermoMajano/snippetbox/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()
	fmt.Print(s)
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
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return

	}

	app.render(w, r, "show.page.tmpl", app.addDefaultData(&TemplateData{Snippet: s}, r))

	/*ts, err
	 := template.ParseFiles(files...)

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

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &TemplateData{
		Form: forms.New(nil),
	})
}

func (app application) createSnippet(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLenth("title", 100)
	form.PermittiedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &TemplateData{Form: form})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expirex"))

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

	w.Write([]byte("Create a new snippet..."))
}

/**func downloadHandler(w http.ResponseWriter, r *http.Request) {
	//fp := filepath.Clean("./ui/static/")
	http.FileServer(http.Dir("./ui/static/"))
}**/

func (app *application) showallT(w http.ResponseWriter, r *http.Request) {
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
