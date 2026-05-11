package awsCognito

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoClient struct {
	client       *cognitoidentityprovider.Client
	clientID     string
	clientSecret string
	userPoolID   string
}

func computeSecretHash(clientSecret, username, clientID string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func New() (*CognitoClient, error) {
	cfg, err := config.LoadDefaultConfig((context.Background()), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {

		return nil, err
	}

	return &CognitoClient{
		client:       cognitoidentityprovider.NewFromConfig(cfg),
		clientID:     os.Getenv("AWS_COGNITO_CLIENT_ID"),
		clientSecret: os.Getenv("AWS_COGNITO_CLIENT_SECRET"),
		userPoolID:   os.Getenv("AWS_COGNITO_USER_POOL_ID"),
	}, nil

}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserName string `json:"username"`
}

func (c *CognitoClient) SignUp(ctx context.Context, req SignUpRequest) (map[string]any, error) {
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

	signupOutput, err := c.client.SignUp(ctx, input)

	if err != nil {
		return nil, err
	}

	return map[string]any{
		"message":   "Sign-up successful. Check your email for the confirmation code.",
		"user_sub":  aws.ToString(signupOutput.UserSub),
		"confirmed": signupOutput.UserConfirmed,
	}, nil
}

type ConfirmRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
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

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *CognitoClient) SignIn(ctx context.Context, req SignInRequest) (map[string]any, error) {
	secretHash := computeSecretHash(c.clientSecret, req.Email, c.clientID)

	authOutput, err := c.client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(c.clientID),
		AuthParameters: map[string]string{
			"USERNAME":    req.Email,
			"PASSWORD":    req.Password,
			"SECRET_HASH": secretHash,
		},
	})

	if err != nil {
		return nil, err
	}

	return map[string]any{
		"access_token":  aws.ToString(authOutput.AuthenticationResult.AccessToken),
		"id_token":      aws.ToString(authOutput.AuthenticationResult.IdToken),
		"refresh_token": aws.ToString(authOutput.AuthenticationResult.RefreshToken),
		"expires_in":    authOutput.AuthenticationResult.ExpiresIn,
	}, nil
}
