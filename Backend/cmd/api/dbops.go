package main

import (
	"context"
	"project/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	filter:=bson.M{"email":user.Email}
	var existingUser models.LoginUser
	err:=collection.FindOne(ctx,filter).Decode(&existingUser)
	if err!=nil{
		return existingUser, err
	}
	return existingUser,nil
}
