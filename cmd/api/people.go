package main

import (
	"net/http"
	"strconv"
)

func (app *application) countPeopleHandler(w http.ResponseWriter, r *http.Request) {
    count, err := app.models.People.Count()

    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte(strconv.Itoa(count)))
    return
}
