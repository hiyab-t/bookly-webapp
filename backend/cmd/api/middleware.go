package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"cafe_store.hiyabnako/internal/data"
	"cafe_store.hiyabnako/internal/validator"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
				
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
					
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
		})
	}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range app.config.cors.trustedOrigins {
			
				if origin == app.config.cors.trustedOrigins[i] {
				
					w.Header().Set("Access-Control-Allow-Origin", origin)

					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,PUT,PATCH,DELETE,GET")
						w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

						w.WriteHeader(http.StatusOK)
						return
				}
				break
			}
			}
		}

		next.ServeHTTP(w, r)
	})
}


func (app *application) Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Vary", "Authorization")

		authZHeader := r.Header.Get("Authorization")

		if authZHeader == "" {
			r = app.contextSetUser(r, data.AnonUser)
			next.ServeHTTP(w, r)
			return
		}


		headerParts := strings.Split(authZHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationToken(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateToken(v, token); !v.Valid() {
			app.invalidCredentials(w,r)
			return
		}

		user, err := app.models.Tokens.GetUserForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationToken(w,r)
			default:
				app.serverError(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(w,r)
	})

}

func (app *application) isAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := app.constextGetUser(r) 

		if user.IsAnon() {
			app.requireAuthenticationRes(w,r)
			return 
		}

		next.ServeHTTP(w,r)
	})
}

func (app *application) requirePermissions(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
			user := app.constextGetUser(r)

			p,err := app.models.Permissions.GetAllForUser(user.User_id)
			if err != nil {
				app.serverError(w,r,err)
				return 
			}

			if !p.Include(code) {
				app.notPermittedRes(w,r)
				return 
			}

			next.ServeHTTP(w,r)

	}

	return app.isAuthenticated(fn)
}
