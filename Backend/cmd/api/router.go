package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) router() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(app.enableCORS)
	mux.Get("/Username",app.GetUsername)
   mux.Get("/Home",app.Hometry)
   mux.Get("/pop",app.popularMovies)
   mux.Get("/movie/{id}",app.GetMovbyid)
   mux.Get("/comments/{id}",app.Getcomsbid)
   mux.Post("/Logout",app.Logout)
   mux.Post("/Register",app.Register)
   mux.Post("/comments",app.PostComment)
   mux.Post("/Login",app.LoginUser)

	return mux
}
