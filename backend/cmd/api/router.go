package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)


func (app *application) router() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowed)

	router.HandlerFunc(http.MethodGet,"/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/books", app.createBookHandler)
	router.HandlerFunc(http.MethodGet, "/v1/books", app.GetAllBooks)
	router.HandlerFunc(http.MethodGet, "/v1/books/:book_id", app.GetBookID)
	router.HandlerFunc(http.MethodPut, "/v1/books", app.PutBook)


	return app.recoverPanic(app.enableCORS(router))
}

