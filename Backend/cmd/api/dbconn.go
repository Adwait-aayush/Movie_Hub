package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbtimeout=time.Second*5
func  ConnectToDB() (*mongo.Client, error) {
	ctx,cancel:=context.WithTimeout(context.Background(),dbtimeout)
	defer cancel()

	client, err := mongo.Connect(ctx,options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	
	 err=client.Ping(ctx,nil)
	 if err != nil {
		return nil, err
	 }
	return client, nil

}
