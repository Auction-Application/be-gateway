package main

import (
	"fmt"
	"log"
	"os"

	awsCognito "github.com/Auction-Application/be-gateway/internal/aws_cognito"
	"github.com/Auction-Application/be-gateway/internal/grpcClient"
	"google.golang.org/grpc"

	"github.com/joho/godotenv"
)

type trustedOrigins struct {
	trustedOrigins []string
}

type config struct {
	httpPort int
	grpcPort int
	env      string
	// db   struct {
	// 	dsn          string
	// 	maxOpenConns int
	// 	maxIdleConns int
	// 	maxIdleTime  time.Duration
	// }

	cors trustedOrigins
}

type application struct {
	// databaseConn *pgx.Conn
	auth      *awsCognito.CognitoClient
	config    config
	rpcClient *grpc.ClientConn
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
		httpPort: 8080,
		grpcPort: 9001,
		cors: trustedOrigins{
			trustedOrigins: []string{os.Getenv("TRUSTED_ORIGIN")},
		},
	}

	fmt.Println("Starting grpc client")

	grpcClient, err := grpcClient.CreategrpcClient(config.grpcPort)
	if err != nil {
		fmt.Printf("Cannot start grpc client:%v", err)
	}
	defer grpcClient.Close()

	app := &application{
		// databaseConn: database.ConnectToDB(),
		auth:      cognitoClient,
		config:    config,
		rpcClient: grpcClient,
	}
	// ctx := context.Background()

	// defer app.databaseConn.Close(ctx)

	fmt.Println("Starting rest server")
	err = app.serve()

	if err != nil {
		fmt.Printf("error in starting the server,%v \n", err)

	}

}
