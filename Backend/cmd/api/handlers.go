package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"project/models"
	"strconv"
	"time"

	"github.com/go-chi/chi"
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

func (app *application) GetMovbyid(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid movie id", http.StatusBadRequest)
		return
	}

	movie, err := app.Gmbid(id)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(movie); err != nil {
		http.Error(w, "Error Displaying Movie", http.StatusInternalServerError)
		return
	}
}

func (app *application) PostComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comments
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
	}
	comment.CommentID = app.GenerateCommentID()
	message, err := app.postcomment(&comment)
	type response struct {
		Message string `json:"message"`
		Error   error  `json:"error"`
	}
	response1 := response{
		Message: message,
		Error:   err,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response1); err != nil {
		return
	}
}

func (app *application) Getcomsbid(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	comments, err := app.getcomments(id)
	if err != nil {
		http.Error(w, "Invalid comment id", http.StatusBadRequest)
	}
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		http.Error(w, "Error Displaying Comments", http.StatusInternalServerError)
	}

}

func (app *application) commentbycmtid(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	comment, err := app.getcommentbyid(id)
	if err != nil {
		http.Error(w, "Invalid comment id", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		http.Error(w, "Error Displaying Comment", http.StatusInternalServerError)
	}

}

func (app *application) addreply(w http.ResponseWriter, r *http.Request) {

	var comment models.Comments
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	for i := range comment.Replies {
		if comment.Replies[i].CommentID == "" {
			comment.Replies[i].CommentID = app.GenerateCommentID()
		}
	}
	message, err := app.addcomments(comment)
	if err != nil {
		http.Error(w, "Error Adding Comment", http.StatusInternalServerError)
		return
	}

	type response struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	response1 := response{
		Message: message,
		Error:   fmt.Sprintf("%v", err),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response1)
	if err != nil {
		http.Error(w, "Error Displaying Comment", http.StatusInternalServerError)
	}
}

func (app *application) Searchmovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("name")

	if len(query) < 3 {
		http.Error(w, "Invalid query", http.StatusBadRequest)
	}
	movies, err := app.Searching(query)
	if err != nil {
		http.Error(w, "Error Searching Movies", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(movies)
	if err != nil {
		http.Error(w, "Error Displaying Movies", http.StatusInternalServerError)
	}
}

func (app *application) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "id")
	if len(commentID) == 0 {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	err := app.Deletecmnt(commentID)
	if err != nil {
		http.Error(w, "Error Deleting Comment", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(err)
	if err != nil {
		http.Error(w, "Error Displaying Comment", http.StatusInternalServerError)
	}

}

func (app *application) addmovies(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	type rresp struct {
		Page    int `json:"page"`
		Results []struct {
			PosterPath string `json:"poster_path"`
		} `json:"results"`
	}

	client := &http.Client{}
	url := fmt.Sprintf("http://api.themoviedb.org/3/search/movie?api_key=%s&query=%s", app.Apikey, url.QueryEscape(movie.Title))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	cont, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	var responseobject rresp
	err = json.Unmarshal(cont, &responseobject)
	if err != nil {
		http.Error(w, "Error unmarshaling response", http.StatusInternalServerError)
		return
	}

	if len(responseobject.Results) > 0 {
		movie.PosterPath = responseobject.Results[0].PosterPath
	}

	id := app.GenerateMovieid()


	movie.ID = id

	response, err := app.ADDmovie(movie)
	if err != nil {
		http.Error(w, "Error creating movie", http.StatusInternalServerError)
		return
	}

	type response1 struct {
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}

	resp1 := response1{
		Message: response,
	
	}

	

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp1)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
