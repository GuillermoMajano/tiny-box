package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.reciverPanic, app.logResquest, secureHeader)

	r := httprouter.New()

	r.GET("/", app.home)
	r.GET("/snippet/create", app.createSnippetForm)
	r.POST("/snippet/create", app.createSnippet)
	r.GET("/snippet/:id", app.showSnippet)

	r.ServeFiles("/static", http.Dir("./ui/static/"))

	return standardMiddleware.Then(r)

}
