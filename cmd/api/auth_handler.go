package main

import (
	"fmt"
	"net/http"

	awsCognito "github.com/Auction-Application/be-gateway/internal/aws_cognito"
	"github.com/Auction-Application/be-gateway/internal/grpcClient"
)

func (app *application) signUpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Sign Up handler")

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req awsCognito.SignUpRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Basic validation
	if !validEmail(req.Email) {
		writeError(w, http.StatusBadRequest, "invalid email")
		return
	}
	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}
	if req.UserName == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	response, err := app.auth.SignUp(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, response)

}

func (app *application) confirmHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Confirm handler")

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req awsCognito.ConfirmRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if req.Username == "" {
		writeError(w, http.StatusBadRequest, "username is required")
		return
	}
	if req.Code == "" {
		writeError(w, http.StatusBadRequest, "code is required")
		return
	}

	if err := app.auth.ConfirmSignUp(r.Context(), req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Account confirmed successfully"})

}

func (app *application) signInHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req awsCognito.SignInRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if req.Username == "" {
		writeError(w, http.StatusBadRequest, "username is required")
		return
	}
	if req.Password == "" {
		writeError(w, http.StatusBadRequest, "password is required")
		return
	}

	signInOutput, err := app.auth.SignIn(r.Context(), req)

	if err != nil {
		fmt.Printf("Error getting signInOutput,%v", err)
		writeError(w, http.StatusBadRequest, "error signing in")
		return
	}

	id_token, ok := (signInOutput["id_token"]).(string)

	if !ok {
		fmt.Println("Error in conversion of jwt from any type to string type")
		writeError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	fmt.Println(id_token)
	userClaim, err := extractIDTokenClaims(id_token)

	if err != nil {
		fmt.Printf("Error extracting IdToken,%w", err)
		writeError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// err = app.insertUser(userClaim)
	err = grpcClient.InsertUserInUserService(userClaim, app.rpcClient)

	if err != nil {
		fmt.Printf("error in inserting user,%v \n", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return

	}

	fmt.Printf("%+v", userClaim)

	writeJSON(w, http.StatusOK, signInOutput)
}
