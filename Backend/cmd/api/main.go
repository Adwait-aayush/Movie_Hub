package main

import (
	"flag"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type application struct {
	Domain string
	DB     *mongo.Client
	Apikey string
}

var err error

func main() {
	var app application
	flag.StringVar(&app.Apikey, "api_key", "0b124068ce1417120e20cb21c5e5213e", "api key")

	app.Domain = "example.com"
	app.DB, err = ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	err := http.ListenAndServe(":4000", app.router())
	if err != nil {
		log.Fatal(err)
	}
}
