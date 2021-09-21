package main

import (
	"bytes"
	"fmt"
	"time"

	"net/http"
	"runtime/debug"
)

func (app *application) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {

	/*	if td != nil {
		td = &TemplateData{}
	}*/
	td.CurrentYear = time.Now().Year()

	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *TemplateData) {
	ts, ok := app.templateCache[name]

	if !ok {
		app.serverError(w, fmt.Errorf("Template file %s no found", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, td.Snippets)

	if nil != err {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
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
