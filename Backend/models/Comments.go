package models

type Comments struct {
	MovieID   string     `json:"movie_id" bson:"movieid"`
	CommentID string     `json:"comment_id" bson:"commentid"`
	Message   string     `json:"message" bson:"message"`
	Likes     int        `json:"likes" bson:"likes"`
	Dislikes  int        `json:"dislikes" bson:"dislikes"`
	Replies   []Comments `json:"replies" bson:"replies"`
	Author    string     `json:"author" bson:"author"`
}
