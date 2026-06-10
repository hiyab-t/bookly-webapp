package main

import (
	"context"
	"net/http"

	"cafe_store.hiyabnako/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.Users) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) constextGetUser(r *http.Request) *data.Users {
	user, ok := r.Context().Value(userContextKey).(*data.Users)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}