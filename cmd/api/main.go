package main

import (
	"fmt"
	"log"
	"os"

	awsCognito "github.com/Auction-Application/be-gateway/internal/aws_cognito"
	"github.com/Auction-Application/be-gateway/internal/grpcClient"

	"github.com/joho/godotenv"
)

type trustedOrigins struct {
	trustedOrigins []string
}

type config struct {
	httpPort    int
	grpcServers grpcClient.GrpcServerTargets
	env         string
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
	rpcClient *grpcClient.Clients
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
		// grpcPort: 9001,
		grpcServers: grpcClient.GrpcServerTargets{Targets: map[grpcClient.GrpcTargetServerName]string{grpcClient.UserService: "localhost:9001", grpcClient.LotService: "localhost:9002", grpcClient.LotImageService: "localhost:9002"}},
		cors: trustedOrigins{
			trustedOrigins: []string{os.Getenv("TRUSTED_ORIGIN")},
		},
	}

	fmt.Println("Starting grpc client")

	grpcClients, err := grpcClient.Dial(config.grpcServers.Targets)
	if err != nil {
		fmt.Println("Error creating grpc clients")
		return
	}

	defer grpcClients.Close()

	app := &application{
		// databaseConn: database.ConnectToDB(),
		auth:      cognitoClient,
		config:    config,
		rpcClient: grpcClients,
	}
	// ctx := context.Background()

	// defer app.databaseConn.Close(ctx)

	fmt.Println("Starting rest server")
	err = app.serve()

	if err != nil {
		fmt.Printf("error in starting the server,%v \n", err)

	}

}
