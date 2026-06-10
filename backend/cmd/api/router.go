package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)


func (app *application) router() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowed)

	// route for health
	router.HandlerFunc(http.MethodGet,"/v1/healthcheck", app.healthcheckHandler)
	
	// route for books
	router.HandlerFunc(http.MethodPost, "/v1/books", app.createBookHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books", app.GetAllBooks)
	router.HandlerFunc(http.MethodGet, "/v1/books/:book_id", app.GetBookID)
	router.HandlerFunc(http.MethodPut, "/v1/books", app.PutBook)
	
	// route for users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.RegisterUsers)

	// route for auth
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationHandler)


	return app.recoverPanic(app.enableCORS(app.Authenticated(router)))
}

