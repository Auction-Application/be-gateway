package main

import (
	"context"

	"github.com/Auction-Application/be-gateway/internal/database/usersTableQuery"
)

func (app *application) insertUser() error {
	ctx := context.Background()
	usersTable := usersTableQuery.New(app.databaseConn)
	err := usersTable.InsertUser(ctx, usersTableQuery.InsertUserParams{Username: "hello", Email: "hello@gmail.com", AuthID: "12345678"})

	if err != nil {
		return err
	}

	return nil

}
