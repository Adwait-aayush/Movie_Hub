package models

type UserMovies struct {
	Adult         bool    `json:"adult" bson:"adult"`
	BackdropPath  string  `json:"backdrop_path" bson:"backdroppath"`
	ID            int     `json:"id" bson:"id"`
	Title         string  `json:"title" bson:"title"`
	OriginalTitle string  `json:"original_title" bson:"originaltitle"`
	OriginalLang  string  `json:"original_language" bson:"originallang"`
	Overview      string  `json:"overview" bson:"overview"`
	ReleaseDate   string  `json:"release_date" bson:"releasedate"`
	PosterPath    string  `json:"poster_path" bson:"posterpath"`
	Popularity    float64 `json:"popularity" bson:"popularity"`
	VoteAverage   float64 `json:"vote_average" bson:"voteaverage"`
	VoteCount     int     `json:"vote_count" bson:"votecount"`
	Author        string  `json:"author" bson:"author"`
}
