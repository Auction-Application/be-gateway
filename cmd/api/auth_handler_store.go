package main

import (
	"context"

	"github.com/Auction-Application/be-gateway/internal/database/usersTableQuery"
)

func (app *application) insertUser(userCredentials idTokenClaims) error {
	ctx := context.Background()
	usersTable := usersTableQuery.New(app.databaseConn)
	err := usersTable.InsertUser(ctx, usersTableQuery.InsertUserParams{Username: userCredentials.Username, Email: userCredentials.Email, AuthID: userCredentials.Sub})

	if err != nil {
		return err
	}

	return nil

}
