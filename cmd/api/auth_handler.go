package main

import (
	"fmt"
	"net/http"

	awsCognito "github.com/Auction-Application/be-gateway/internal/aws_cognito"
)


func (app *application) signUpHandler(w http.ResponseWriter, r *http.Request)  {
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

func (app *application) confirmHandler(w http.ResponseWriter, r *http.Request)  {
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

		if !validEmail(req.Email) {
			writeError(w, http.StatusBadRequest, "invalid email")
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
	fmt.Println("SignIn handler")

		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req awsCognito.SignInRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}

		if !validEmail(req.Email) {
			writeError(w, http.StatusBadRequest, "invalid email")
			return
		}
		if req.Password == "" {
			writeError(w, http.StatusBadRequest, "password is required")
			return
		}

		signInOutput, err := app.auth.SignIn(r.Context(), req)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, signInOutput)
}