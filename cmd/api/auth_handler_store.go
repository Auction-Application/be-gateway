package main

import "github.com/Auction-Application/be-gateway/internal/grpcClient"

func (app *application) insertUser(userCredentials grpcClient.IdTokenClaims) error {
	// ctx := context.Background()
	// usersTable := usersTableQuery.New(app.databaseConn)
	// err := usersTable.InsertUser(ctx, usersTableQuery.InsertUserParams{Username: userCredentials.Username, Email: userCredentials.Email, AuthID: userCredentials.Sub})

	// if err != nil {
	// 	return err
	// }

	return nil

}
