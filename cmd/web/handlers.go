package main

import (
	"GuillermoMajano/snippetbox/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
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
	app.render(w, r, "create.page.tmpl", nil)
}

func (app application) createSnippet(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	errors := make(map[string]string)

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"

	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximun is 100 character)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "this field cannot be blank"
	}

	if strings.TrimSpace(expires) == "" {

	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	if len(errors) > 0 {
		fmt.Fprint(w, errors)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)

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
