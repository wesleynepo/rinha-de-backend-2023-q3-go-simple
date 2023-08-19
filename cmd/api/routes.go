package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)



func (app *application) routes() http.Handler {
    router := httprouter.New()

    router.HandlerFunc(http.MethodGet, "/pessoas/:id", app.showPersonHandler)
    router.HandlerFunc(http.MethodGet, "/pessoas", app.searchPeopleHandler)
    router.HandlerFunc(http.MethodGet, "/contagem-pessoas", app.countPeopleHandler)
    router.HandlerFunc(http.MethodPost, "/pessoas", app.createPeopleHandler)

    return router
}
