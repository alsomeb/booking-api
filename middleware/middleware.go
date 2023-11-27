package middleware

import (
	"booking-api/auth"
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func TokenVerificationMiddleware(firebaseClient *firebase.App) gin.HandlerFunc {
	// anonymous func
	return func(c *gin.Context) {
		// Get the Auth client
		authClient, err := auth.GetAuthClient(firebaseClient)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Firebase Auth Error"})
			return
		}

		// Get the token from the request header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
			return
		}

		tokenParams := strings.Fields(token)
		if len(tokenParams) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		// Verify the token
		decodedToken, err := authClient.VerifyIDToken(context.Background(), tokenParams[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims := decodedToken.Claims

		// Store user information in the context for later use in handlers
		c.Set("userClaims", claims)

		//log.Print("TokenVerificationMiddleware was used")

		// Continue to the next middleware or route handler
		c.Next()
	}
}
