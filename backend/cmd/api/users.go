package main

import (
	"net/http"

	"cafe_store.hiyabnako/internal/data"
	"cafe_store.hiyabnako/internal/validator"
)

func (app *application) RegisterUsers(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user := &data.Users{
		Name: input.Name,
		Email: input.Email,
		Active: false,
	}

	if err = user.Password.Set(input.Password); err != nil {
		app.serverError(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUsers(v, user); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	if err = app.models.Users.InsertUser(user); err != nil {
		app.serverError(w, r, err)
		return
	}

}

func (app *application) GetUsers(w http.ResponseWriter, r *http.Request) {


	

}