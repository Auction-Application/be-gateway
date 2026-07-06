package grpcClient

import (
	"context"
	"fmt"

	userPb "github.com/Auction-Application/be-gateway/rpc/gen/user/v1"
)

type IdTokenClaims struct {
	Sub      string `json:"sub"`
	Email    string `json:"email"`
	Username string `json:"cognito:username"`
}

// func InsertUserInUserService(userCredentials IdTokenClaims, conn *grpc.ClientConn) error {
func InsertUserInUserService(userCredentials IdTokenClaims, userClient userPb.UserServiceClient) error {
	// userClient := userPb.NewUserServiceClient(conn)
	insertUserResponse, err := userClient.InsertUser(context.Background(), &userPb.InsertUserPayload{
		User: &userPb.User{
			Username: &userCredentials.Username,
			Email:    &userCredentials.Email,
			AuthID:   &userCredentials.Sub,
			// AuthID: new(string),
		},
	})

	fmt.Println(insertUserResponse)

	if err != nil {
		return err
	}

	return nil

}
