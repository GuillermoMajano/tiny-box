package main

import (
	"bytes"
	"fmt"
	"time"

	"net/http"
	"runtime/debug"
)

func (app *application) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {

	if td == nil {
		td = &TemplateData{}
	}

	td.CurrentYear = time.Now().Year()

	td.Flash = app.session.PopString(r, "flash")
	td.IsAuthenticated = app.isAuthenticated(r)
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *TemplateData) {
	ts, ok := app.templateCache[name]

	if !ok {
		app.serverError(w, fmt.Errorf("Template file %s no found", name))
		return
	}

	buf := new(bytes.Buffer)

	var err error

	if td != nil {
		td = app.addDefaultData(td, r)
		err = ts.Execute(buf, td)
	} else {
		err = ts.Execute(buf, nil)
	}

	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

// Return true if the current request is from authenticated user, otherwise return false.
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedUserID")
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorlog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
