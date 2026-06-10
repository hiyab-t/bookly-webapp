package main

import (
	"errors"
	"net/http"
	"time"

	"cafe_store.hiyabnako/internal/data"
	"cafe_store.hiyabnako/internal/validator"
)

func (app *application) createAuthenticationHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w,r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePassword(v, input.Password)

	if !v.Valid() {
		app.FailedValidationResponse(w,r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentials(w,r)
		default:
			app.serverError(w,r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		switch {
		case !match:
			app.invalidCredentials(w,r)
		default:
			app.serverError(w,r,err)
		}
		return
	}

	token, err := app.models.Tokens.New(int64(user.User_id),24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverError(w, r,err)
		return
	}

	if err = app.writeJSON(w,http.StatusCreated, envelope{"authentication_token": token}, nil); err != nil {
		app.serverError(w,r,err)
	}
}

