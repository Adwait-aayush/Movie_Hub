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
   mux.Get("/replycomments/{id}",app.commentbycmtid)
   mux.Get("/user",app.UserMovies)
   mux.Get("/Search",app.Searchmovies)
   mux.Post("/Logout",app.Logout)
   mux.Post("/Register",app.Register)
   mux.Post("/comments",app.PostComment)
   mux.Post("/Login",app.LoginUser)
   mux.Post("/addusermovies",app.addmovies)
   mux.Patch("/addreply",app.addreply)
   mux.Delete("/Deletecomment/{id}",app.DeleteComment)
   mux.Delete("/Delete/{id}",app.DeleteMovies)

	return mux
}
