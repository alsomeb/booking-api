package auth

import (
	"context"
	"errors"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"log"
	"strings"
)

func InitFireBaseApp() *firebase.App {
	opt := option.WithCredentialsFile("private_key.json")
	config := &firebase.Config{ProjectID: "fir-auth-react-af3e5"}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	log.Println("Firebase App Initialized")
	return app
}

// GetUserData retrieves user data based on the provided Firebase client and Gin context.
func GetUserData(client *firebase.App, ctx *gin.Context) (*auth.UserRecord, error) {
	// The gin Context in func params == carries information about the current HTTP request

	// Get the Auth client
	authClient, err := client.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	// Get the token from the request header
	token := ctx.GetHeader("Authorization")

	if token == "" {
		return nil, errors.New("authorization token is missing")
	}

	tokenParams := strings.Fields(token)
	if len(tokenParams) != 2 {
		return nil, errors.New("invalid token or format")
	}

	// Verify the token
	decodedToken, err := authClient.VerifyIDToken(context.Background(), tokenParams[1])
	if err != nil {
		return nil, errors.New("invalid token or format")
	}

	// Extract user information
	email, ok := decodedToken.Claims["email"].(string)
	if !ok {
		return nil, errors.New("failed to extract user email")
	}

	// Get user record by email
	userRecord, err := authClient.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}

	return userRecord, nil
}
