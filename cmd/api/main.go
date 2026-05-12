package main

import (
	"context"
	"fmt"
	"log"
	"time"

	awsCognito "github.com/Auction-Application/be-gateway/internal/aws_cognito"
	"github.com/Auction-Application/be-gateway/internal/database"
	"github.com/Auction-Application/be-gateway/internal/database/usersTableQuery"
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

	usersTable := usersTableQuery.New(app.databaseConn)

	err = usersTable.InsertUser(context.Background(), usersTableQuery.InsertUserParams{Username: "hello", Email: "hello@gmail.com", AuthID: "12345678"})

	if err != nil {
		fmt.Printf("error in inserting user,%v \n", err)

	}
	err = app.serve()

	if err != nil {
		fmt.Printf("error in starting the server,%v \n", err)

	}

}
