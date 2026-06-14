package main

import (
	"fmt"
	"net/http"

	"golang.org/x/text/message"
)

func (app *application) logError(r *http.Request, err error) {

	
	method := r.Method
	uri := r.URL.RequestURI()


	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errResponse(w, r, http.StatusBadRequest, err.Error())
} 

func (app *application) errResponse(w http.ResponseWriter, r *http.Request, status int, message any) {

	env := envelope{"error": message}

	err := app.writeJSON(w,status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {

	app.logError(r, err)

	message := "sorry, could not process your request"

	app.errResponse(w, r, http.StatusInternalServerError, message)

}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {

	message := "resource for your request not found"

	app.errResponse(w, r, http.StatusNotFound, message)

}

func (app *application) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {

	message := fmt.Sprintf("%s method request not allowed for this response", r.Method)

	app.errResponse(w, r, http.StatusMethodNotAllowed, message)

}

func (app *application) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {

	app.errResponse(w ,r, http.StatusUnprocessableEntity, errors)
}

func (app *application) invalidCredentials(w http.ResponseWriter, r *http.Request) {
	
	message := "invalid authentication credentials"

	app.errResponse(w, r, http.StatusUnauthorized, message)

}

func (app *application) invalidAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"

	app.errResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) requireAuthenticationRes(w http.ResponseWriter, r *http.Request) {
	
	message := "you must be authenticated to access this resource"

	app.errResponse(w,r,http.StatusUnauthorized, message)
}

func (app *application) notPermittedRes(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"

	app.errResponse(w,r, http.StatusForbidden,message)

}