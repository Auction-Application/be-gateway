package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/joho/godotenv"
)


type CognitoClient struct{
	client *cognitoidentityprovider.Client
	clientID string
	clientSecret string
	userPoolID string
}

type SignUpRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
	UserName string `json:"username"`
}

type ConfirmRequest struct{
	Email string `json:"email"`
	Code string `json:"code"`
}

type SignInRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func computeSecretHash(clientSecret,username,clientID string) string{
	mac:=hmac.New(sha256.New,[]byte(clientSecret))
	mac.Write([]byte(username+clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func NewCognitoClient()(*CognitoClient,error){
	cfg,err:=config.LoadDefaultConfig((context.TODO()),config.WithRegion(os.Getenv("AWS_REGION")))
	if err!=nil{
		return nil,err
	}

	return &CognitoClient{
		client:cognitoidentityprovider.NewFromConfig(cfg) ,
		clientID: os.Getenv("AWS_COGNITO_CLIENT_ID"),
		clientSecret: os.Getenv("AWS_COGNITO_CLIENT_SECRET"),
		userPoolID: os.Getenv("AWS_COGNITO_USER_POOL_ID"),

	},nil

	
	

}

func (c *CognitoClient) SignUp(ctx context.Context, req SignUpRequest) (*cognitoidentityprovider.SignUpOutput, error) {
	secretHash := computeSecretHash(c.clientSecret, req.Email, c.clientID)

	input := &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(c.clientID),
		SecretHash: aws.String(secretHash),
		Username:   aws.String(req.Email),
		Password:   aws.String(req.Password),
		// UserAttributes: []types.AttributeType{
		// 	{Name: aws.String("email"), Value: aws.String(req.Email)},
		// 	{Name: aws.String("userName"), Value: aws.String(req.UserName)},
		// },
	}

	return c.client.SignUp(ctx, input)
}

func (c *CognitoClient) ConfirmSignUp(ctx context.Context, req ConfirmRequest) error {
	secretHash := computeSecretHash(c.clientSecret, req.Email, c.clientID)

	_, err := c.client.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientID),
		SecretHash:       aws.String(secretHash),
		Username:         aws.String(req.Email),
		ConfirmationCode: aws.String(req.Code),
	})
	return err
}

func (c *CognitoClient) SignIn(ctx context.Context, req SignInRequest) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	secretHash := computeSecretHash(c.clientSecret, req.Email, c.clientID)

	return c.client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(c.clientID),
		AuthParameters: map[string]string{
			"USERNAME":    req.Email,
			"PASSWORD":    req.Password,
			"SECRET_HASH": secretHash,
		},
	})
}


// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func decodeJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func validEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

// --- Handlers ---

func signUpHandler(c *CognitoClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req SignUpRequest
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

		out, err := c.SignUp(r.Context(), req)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"message":   "Sign-up successful. Check your email for the confirmation code.",
			"user_sub":  aws.ToString(out.UserSub),
			"confirmed": out.UserConfirmed,
		})
	}
}

func confirmHandler(c *CognitoClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req ConfirmRequest
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

		if err := c.ConfirmSignUp(r.Context(), req); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"message": "Account confirmed successfully"})
	}
}

func signInHandler(c *CognitoClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req SignInRequest
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

		out, err := c.SignIn(r.Context(), req)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"access_token":  aws.ToString(out.AuthenticationResult.AccessToken),
			"id_token":      aws.ToString(out.AuthenticationResult.IdToken),
			"refresh_token": aws.ToString(out.AuthenticationResult.RefreshToken),
			"expires_in":    out.AuthenticationResult.ExpiresIn,
		})
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cognito, err := NewCognitoClient()
	if err != nil {
		log.Fatalf("Failed to initialize Cognito client: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/signup", signUpHandler(cognito))
	mux.HandleFunc("/confirm", confirmHandler(cognito))
	mux.HandleFunc("/signin", signInHandler(cognito))

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}