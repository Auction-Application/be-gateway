package main

import (
	"context"
	"fmt"
	"log"
	"time"

	awsCognito "github.com/Auction-Application/be-gateway/internal/aws_cognito"
	"github.com/Auction-Application/be-gateway/internal/database"
	"github.com/jackc/pgx/v5"
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
	databaseConn *pgx.Conn
	auth         *awsCognito.CognitoClient
	config       config
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
		databaseConn: database.ConnectToDB(),
		auth:         cognitoClient,
		config:       config,
	}
	ctx := context.Background()

	defer app.databaseConn.Close(ctx)

	err = app.serve()

	if err != nil {
		fmt.Printf("error in starting the server,%v \n", err)

	}

}
