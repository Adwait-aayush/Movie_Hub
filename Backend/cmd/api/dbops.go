package main

import (
	"context"
	"project/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) RegisterUser(user *models.LoginUser) (string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	collection := app.DB.Database("MovieHub").Collection("Users")

	var existingUser models.LoginUser
	err := collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		return "Username is already taken", 400, nil
	}

	err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return "Email is already taken", 400, nil
	}

	hashedp, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return "Error generating password hash", 500, err
	}
	user.Password = string(hashedp)
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return "Failed to register user", 500, err
	}

	return "User registered successfully", 200, nil
}

func (app *application) Login(user *models.LoginUser) (string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	collection := app.DB.Database("MovieHub").Collection("Users")
	var existingUser models.LoginUser

	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "User Does not Exist", 401, nil
		}
		return "Database error", 500, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return "Wrong Password", 401, nil
	}

	return "Login successful", 200, nil
}
func (app *application) FindUser(user *models.LoginUser) (models.LoginUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	collection := app.DB.Database("MovieHub").Collection("Users")
	filter := bson.M{"email": user.Email}
	var existingUser models.LoginUser
	err := collection.FindOne(ctx, filter).Decode(&existingUser)
	if err != nil {
		return existingUser, err
	}
	return existingUser, nil
}
func (app *application) Gmbid(id int) (models.Movie, error) {
	var movie models.Movie
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	collection := app.DB.Database("MovieHub").Collection("MoviesPopular")
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&movie)
	if err != nil {
		return movie, err
	}
	return movie, nil
}

func (app *application) GenerateCommentID() string {
	return uuid.New().String()
}

func (app *application) postcomment(comment *models.Comments) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	collection := app.DB.Database("MovieHub").Collection("Comments")
	_, err := collection.InsertOne(ctx, comment)
	if err != nil {
		return "Unable to insert comment", err
	}
	return "Comment posted successfully", nil
}
func (app *application) getcomments(id string) ([]models.Comments, error) {
	var comments []models.Comments
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	collection := app.DB.Database("MovieHub").Collection("Comments")
	filter := bson.M{"movieid": id}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return comments, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var comment models.Comments
		err := cursor.Decode(&comment)
		if err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}
	if err := cursor.Err(); err != nil {
		return comments, err
	}
	return comments, nil

}

func (app *application) getcommentbyid(id string) (models.Comments, error) {
	var comment models.Comments
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	collection := app.DB.Database("MovieHub").Collection("Comments")
	filter := bson.M{
		"$or": []bson.M{
			{"commentid": id},
			{"replies.commentid": id},
		},
	}

	err := collection.FindOne(ctx, filter).Decode(&comment)
	if err != nil {
		return comment, err
	}

	return comment, nil

}
func (app *application) addcomments(comment models.Comments) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	collection := app.DB.Database("MovieHub").Collection("Comments")

	// Create a filter to find the comment by movie ID and comment ID
	filter := bson.M{
		"movieid":   comment.MovieID,
		"commentid": comment.CommentID,
	}

	// Find the existing document
	var existingComment models.Comments
	err := collection.FindOne(ctx, filter).Decode(&existingComment)
	if err != nil && err != mongo.ErrNoDocuments {
		return "Unable to fetch existing comment", err
	}

	// Create a map of existing reply IDs for quick lookup
	existingReplyIDs := make(map[string]bool)
	for _, reply := range existingComment.Replies {
		existingReplyIDs[reply.CommentID] = true
	}

	// Filter out duplicate replies
	uniqueReplies := []models.Comments{}
	for _, reply := range comment.Replies {
		if !existingReplyIDs[reply.CommentID] {
			uniqueReplies = append(uniqueReplies, reply)
		}
	}

	// If there are no new unique replies, return early
	if len(uniqueReplies) == 0 {
		return "No new replies to add", nil
	}

	// Use $push with $each to add only unique replies
	update := bson.M{
		"$push": bson.M{
			"replies": bson.M{
				"$each": uniqueReplies, // Add only unique replies
			},
		},
	}

	// Perform the update operation with upsert enabled
	result, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return "Unable to update or insert comments", err
	}

	// Check if the document was matched or inserted
	if result.MatchedCount == 0 {
		return "No matching comment found, a new one was inserted", nil
	}

	return "Comment updated with unique replies successfully", nil
}

func (app *application) Searching(name string) ([]models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	collection := app.DB.Database("MovieHub").Collection("MoviesPopular")
	filter := bson.M{
		"$or": []interface{}{
			bson.M{"name": bson.M{"$regex": name, "$options": "i"}},
			bson.M{"originaltitle": bson.M{"$regex": name, "$options": "i"}},
		},
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []models.Movie
	for cursor.Next(ctx) {
		var movie models.Movie
		if err := cursor.Decode(&movie); err != nil {
			return nil, err

		}
		movies = append(movies, movie)
	}
	return movies, nil
}
func (app *application) Deletecmnt(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	collection := app.DB.Database("MovieHub").Collection("Comments")
	filter := bson.M{
		"$or": []bson.M{
			{"commentid": id},
			{"replies.commentid": id},
		},
	}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil

}


