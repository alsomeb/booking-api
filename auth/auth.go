package auth

import (
	"context"
	"errors"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"log"
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

func GetAuthClient(client *firebase.App) (*auth.Client, error) {
	// Get the Auth client
	authClient, err := client.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return authClient, nil
}

// GetUserData retrieves user data based on the provided Firebase client and Gin context.
func GetUserData(client *firebase.App, ctx *gin.Context) (*auth.UserRecord, error) {
	authClient, err := GetAuthClient(client)

	// Extract user information
	claims, ok := ctx.Get("userClaims")
	if !ok {
		return nil, errors.New("failed to fetch user value from context")
	}

	claimsMap := claims.(map[string]interface{}) // dynamic we don't know which type key value thus interface{}
	email := claimsMap["email"].(string)         // type assert email as string

	// Get user record by email
	userRecord, err := authClient.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, errors.New("failed to retrieve User record")
	}

	return userRecord, nil
}
