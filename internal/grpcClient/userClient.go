package grpcClient

import (
	"context"
	"fmt"

	userPb "github.com/Auction-Application/be-gateway/rpc/gen/user/v1"
	"google.golang.org/grpc"
)

type IdTokenClaims struct {
	Sub      string `json:"sub"`
	Email    string `json:"email"`
	Username string `json:"cognito:username"`
}

func InsertUserInUserService(userCredentials IdTokenClaims, conn *grpc.ClientConn) error {
	userClient := userPb.NewUserServiceClient(conn)
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
	// ctx := context.Background()
	// usersTable := usersTableQuery.New(app.databaseConn)
	// err := usersTable.InsertUser(ctx, usersTableQuery.InsertUserParams{Username: userCredentials.Username, Email: userCredentials.Email, AuthID: userCredentials.Sub})

	// if err != nil {
	// 	return err
	// }

}
