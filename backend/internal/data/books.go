package data

import (
	"context"
	"errors"

	//"fmt"
	"time"

	//"github.com/jackc/pgx/pgtype"
	//"github.com/upper/db/v4/adapter/postgresql"
	//"github.com/jackc/pgx/v5/pgtype"
	//"github.com/lib/pq"
	//"github.com/jackc/pgtype"

	"cafe_store.hiyabnako/internal/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Book struct {
	Book_id int `json:"book_id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Genres []string `json:"genres"`
	//Year int32 `json:"year"`
	Version int32 `json:"version"`
}

type BookModel struct {
	DBpool *pgxpool.Pool
}

func ValidateBook(v *validator.Validator, b *Book) {
	v.Check(b.Title != "", "title", "must not be empty")
	v.Check(len(b.Title) > 400, "title", "must not be more than 400 bytes long")

	v.Check(b.Author != "", "author", "must not be empty")
	v.Check(len(b.Author) > 300, "author", "must not be more than 300 bytes long")

	v.Check(b.Genres != nil, "genres", "must not be empty")
	v.Check(len(b.Genres) >= 1, "genres","must contain at least 1 genre")
	v.Check(len(b.Genres) <= 5,"genres", "must not contain more than 5 genres")

	v.Check(validator.Unique(b.Genres), "genres", "must not have duplicate vlaues")


} 

func (m *BookModel) AllBooks(title string, author string, genres []string, filters Filters) ([]*Book,Metadata, error) {

	qry := `select count(*) over(), book_id,title, author, genres, version 
	from books 
	where (lower(title) = lower($1) OR $1 = '') 
	and (lower(author) = lower($2) or $2 = '')
	and (genres @> $3 OR $3 = '{}')
	order by book_id
	limit $4 offset $5;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DBpool.Query(ctx, qry, title, author,genres, filters.limit(), filters.Offset() )
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0

	books := []*Book{}

	for rows.Next() {

		b := &Book{}

		err = rows.Scan(&totalRecords, &b.Book_id,&b.Title, &b.Author, &b.Genres, &b.Version)
		if err != nil {
			return nil,Metadata{}, err
		}

		books = append(books, b)

	}

	if rows.Err() != nil {
		return nil,Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	
	return books, metadata,nil
	
}

func (m *BookModel) GetBookID(book_id int) (*Book, error) {

	if book_id < 1 {
		return nil, ErrRecordNotFound
	}

	qry := `select * from books where book_id = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := m.DBpool.QueryRow(ctx, qry, book_id)

	b := &Book{}

	if err := row.Scan(&b.Book_id,&b.Title, &b.Author, &b.Genres, &b.Version); err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	return b, nil
}

func (m *BookModel) CreateBook(book *Book) error {
	
	qry := `insert into books (title, author, genres) values ($1, $2, $3) returning book_id, version;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{book.Title, book.Author, book.Genres}

	err := m.DBpool.QueryRow(ctx, qry, args...).Scan(&book.Book_id, &book.Version)
	if err != nil {
		return err 
	}

	return nil
}

func (m *BookModel) UpdateBook(book *Book) error {

	qry := `update books 
			set title = $1, author = $2, genres = $3, version = version + 1 
			where book_id = $4;`

	args := []any{book.Title, book.Author, book.Genres, book.Book_id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DBpool.Exec(ctx, qry, args...)
	if err != nil {
		return err
	}

	return nil
}