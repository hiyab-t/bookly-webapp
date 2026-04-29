package main

import (

	"net/http"
)

func (app *application) healthcheckHandler(w  http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "Available",
		"system_info":map[string]string{
		"environment": app.config.env,
		"version": version,
		},
	}

	err := app.writeJSON(w, env, nil)
	if err != nil {
		app.serverError(w, r, err)
	}

}

