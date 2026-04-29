package main

import (
	//"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"cafe_store.hiyabnako/internal/data"
	"cafe_store.hiyabnako/internal/validator"
)

// create book handler

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title string `json:"title"`
		Author string `json:"author"`
		Genres []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Error(err.Error())
		app.badRequest(w, r, err)
		return
	}

	book := &data.Book{
		Title: input.Title,
		Author: input.Author,
		Genres: input.Genres,
	}

	v := validator.New()

	if data.ValidateBook(v, book); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
	} 

	err = app.models.Books.CreateBook(book)
	if err != nil {
		app.serverError(w,r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.Book_id))

	err = app.writeJSON(w, envelope{"book": book}, headers)
	if err != nil {
		app.serverError(w,r,err)
		return
	}

}

func (app *application) GetAllBooks(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title string
		Author string
		Genres []string
		data.Filters
	}

	qs := r.URL.Query()

	v := validator.New()

	input.Title = app.readString(qs, "title", "")
	input.Author = app.readString(qs, "author", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)

	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "book_id")

	input.Filters.SortSafeList = []string{"title", "author", "-title", "book_id", "-book_id", "-author"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.errResponse(w, r, http.StatusUnprocessableEntity, v.Errors)
		return
	}

	books,  metadata, err := app.models.Books.AllBooks(input.Title, input.Author, input.Genres, input.Filters)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.writeJSON(w, envelope{"metadata": metadata,"books": books,}, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

}

func (app *application) GetBookID(w http.ResponseWriter, r *http.Request) {
	
	book_id, err := app.readIDparam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/:%v", book_id))

	book, err := app.models.Books.GetBookID(int(book_id))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if err = app.writeJSON(w, envelope{"book": book}, headers); err != nil {
		app.serverError(w, r, err)
		return
	}

}

func (app *application) PutBook(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Book_id string `json:"book_id"`
		Title string `json:"title"`
		Author string `json:"author"`
		Genres []string `json:"genres"`
	}

	//book_id, err := app.readIDparam(r)
	//if err != nil {
	//	app.notFoundResponse(w, r)
	//	return
	//}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	book_id, err := strconv.Atoi(input.Book_id)
	if err != nil {
		app.errResponse(w, r, http.StatusNotFound, "invalid Book ID")
	}

	_, err = app.models.Books.GetBookID(book_id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverError(w, r, err)
		}
		return 
	}

	book := &data.Book{
		Book_id: book_id,
		Title: input.Title, 
		Author: input.Author,
		Genres: input.Genres,}

	v := validator.New()

	if data.ValidateBook(v, book); v.Valid() {
		app.errResponse(w, r, http.StatusUnprocessableEntity, v.Errors)
		return
	}

	if err = app.models.Books.UpdateBook(book); err != nil {
		app.serverError(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.Book_id))

	err = app.writeJSON(w, envelope{"book":book}, headers)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

