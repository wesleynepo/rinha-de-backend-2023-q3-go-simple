package main

import (
	"errors"
	"fmt"
	"gopherinha/internal/data"
	"gopherinha/internal/validator"
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

func (app *application) createPeopleHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Apelido string `json:"apelido"`
        Nome string `json:"nome"`
        Nascimento string `json:"nascimento"`
        Stack []string `json:"stack"`
    }

    err := app.readJSON(w, r, &input)
    if err != nil {
        app.badRequestResponse(w, r, err)
        return
    }

    person := &data.Person{
        Nome: input.Nome,
        Apelido: input.Apelido,
        Nascimento: input.Nascimento,
        Stack: input.Stack,
    }

    v := validator.New()

    if data.ValidatePerson(v, person); !v.Valid() {
        app.failedValidationResponse(w, r, v.Errors)
        return
    }

    err = app.models.People.Insert(person)

    if err != nil {
        switch {
        case errors.Is(err, data.ErrDuplicateApelido):
            app.failedUnprocessable(w, r)
        default:
            app.serverErrorResponse(w, r, err)
        }
        return
    }

    headers := make(http.Header)
    headers.Set("Location", fmt.Sprintf("/pessoas/%s", person.UUID))

    err = app.writeJSON(w, http.StatusCreated, envelope{"pessoa": person}, headers)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}

func (app *application) showPersonHandler(w http.ResponseWriter, r *http.Request) {
    id := app.readIDParam(r)

    person, err := app.models.People.Get(id)

    if err != nil {
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            app.notFoundResponse(w, r)
        default:
            app.serverErrorResponse(w, r, err)
        }
        return
    }

    err = app.writeJSON(w, http.StatusOK, person, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    } 
}
