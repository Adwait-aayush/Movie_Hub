package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"project/models"
	"time"

	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)



func (app *application) Hometry(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}

	resp := response{
		Message: "Hello, World!",
		Status:  http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func (app *application) popularMovies(w http.ResponseWriter, r *http.Request) {
	type Movie struct {
		Adult         bool    `json:"adult"`
		BackdropPath  string  `json:"backdrop_path"`
		ID            int     `json:"id"`
		Title         string  `json:"title"`
		OriginalTitle string  `json:"original_title"`
		OriginalLang  string  `json:"original_language"`
		Overview      string  `json:"overview"`
		ReleaseDate   string  `json:"release_date"`
		PosterPath    string  `json:"poster_path"`
		Popularity    float64 `json:"popularity"`
		VoteAverage   float64 `json:"vote_average"`
		VoteCount     int     `json:"vote_count"`
	}

	type MoviesResponse struct {
		Page         int     `json:"page"`
		Results      []Movie `json:"results"`
		TotalPages   int     `json:"total_pages"`
		TotalResults int     `json:"total_results"`
	}

	var moviesResponse MoviesResponse
	client := &http.Client{}
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s&language=en-US", app.Apikey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal([]byte(body), &moviesResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	collection := app.DB.Database("MovieHub").Collection("MoviesPopular")

	for _, movie := range moviesResponse.Results {
		filter := bson.M{"id": movie.ID}
		err := collection.FindOne(ctx, filter).Err()
		if err == mongo.ErrNoDocuments {
			_, err = collection.InsertOne(ctx, movie)
			if err != nil {
				continue
			}
		} else if err != nil {
			continue
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(moviesResponse.Results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *application) Register(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("super-secret"))
	var user models.LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type Response struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
		Error   error  `json:"error"`
	}
	message, status, error := app.RegisterUser(&user)

	// Set the session for the user after successful registration
	session, err := store.Get(r, "session-id")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Setting username in the session
	session.Values["username"] = user.Username
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Message: message,
		Status:  status,
		Error:   error,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *application) LoginUser(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("super-secret"))
	var user models.LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type Response struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
		Error   error  `json:"error"`
	}
	message, status, error := app.Login(&user)

	// Find the user in the database
	existingUser, err := app.FindUser(&user)
	if err != nil {
		http.Error(w, "Failed to find user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the session for the logged-in user
	session, err := store.Get(r, "session-id")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Setting username in the session
	session.Values["username"] = existingUser.Username
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Message: message,
		Status:  status,
		Error:   error,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *application) GetUsername(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("super-secret"))
	session, err := store.Get(r, "session-id")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if username exists in session
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Error(w, "Username not found in session", http.StatusBadRequest)
		return
	}

	// Respond with username
	type response struct {
		Username string `json:"username"`
	}
	response1 := response{
		Username: username,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response1)
}
func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("super-secret"))
    session, err := store.Get(r, "session-id")
    if err != nil {
        http.Error(w, "Error retrieving session: "+err.Error(), http.StatusInternalServerError)
        return
    }

    session.Values = nil

    err = session.Save(r, w)
    if err != nil {
        http.Error(w, "Error saving session: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Logout successful",
    })
}
