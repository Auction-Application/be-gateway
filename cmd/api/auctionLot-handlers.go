package main

import (
	"net/http"
)

func (app *application) createAuctionLot(w http.ResponseWriter, r *http.Request) {

	// if req.Username == "" {
	// 	writeError(w, http.StatusBadRequest, "username is required")
	// 	return
	// }
	// if req.Password == "" {
	// 	writeError(w, http.StatusBadRequest, "password is required")
	// 	return
	// }

	// if err != nil {
	// 	fmt.Printf("Error getting signInOutput,%v", err)
	// 	writeError(w, http.StatusBadRequest, "error signing in")
	// 	return
	// }

	// if !ok {
	// 	fmt.Println("Error in conversion of jwt from any type to string type")
	// 	writeError(w, http.StatusInternalServerError, "Something went wrong")
	// 	return
	// }

	// if err != nil {
	// 	fmt.Printf("Error extracting IdToken,%w", err)
	// 	writeError(w, http.StatusInternalServerError, "Something went wrong")
	// 	return
	// }

	// // err = grpcClient.InsertUserInUserService(userClaim, app.rpcClient)

	// if err != nil {
	// 	fmt.Printf("error in inserting user,%v \n", err)
	// 	writeError(w, http.StatusInternalServerError, err.Error())
	// 	return

	// }

	// fmt.Printf("%+v", userClaim)

	type hello struct {
		Name string `json:"name"`
	}

	writeJSON(w, http.StatusOK, hello{Name: "hello"})
}
