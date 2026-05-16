package grpcClient

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreategrpcClient(port int) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Printf("Cannot start grpc client:%v", err)
		return nil, err

	}

	return conn, nil
}
