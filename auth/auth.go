package auth

import (
	"context"
	firebase "firebase.google.com/go/v4"
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
