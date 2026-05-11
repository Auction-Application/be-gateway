package main

import (
	"fmt"
	"log"
	"time"

	awsCognito "github.com/Auction-Application/be-gateway/internal/aws_cognito"
	"github.com/Auction-Application/be-gateway/internal/data"
	"github.com/joho/godotenv"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}

	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	models data.Models
	auth   *awsCognito.CognitoClient
	config config
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
		fmt.Println(err)
		return
	}

	cognitoClient, err := awsCognito.New()
	if err != nil {
		log.Println("error in cognito initialization")
	}

	config := config{
		port: 8080,
	}

	app := &application{
		models: "sdfsdf",
		auth:   cognitoClient,
		config: config,
	}

	err = app.serve()

	if err != nil {
		fmt.Printf("error in starting the server,%v \n", err)

	}

}
