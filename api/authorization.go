package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoggedInUser struct {
	ID    string
	Email string
}

func authorizeRoutes(appCtx *AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// If we dont have a token header, only accept the request if it is an anonymous route.
			if !isAnonymousRequest(ctx.Request) {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			ctx.Next()
			return
		}

		// Make sure anonymous routes cannot be accessed with authorized users
		if isAnonymousRequest(ctx.Request) {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Parse the token and set the loggedin user variables
		claims, err := parseToken(appCtx, extractToken(tokenHeader))
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("loggedInUser", LoggedInUser{
			ID:    claims.UserID,
			Email: claims.UserEmail,
		})

		ctx.Next()
	}
}

var anonymousPaths = []string{"/users", "/authenticate"}

func isAnonymousRequest(request *http.Request) bool {
	if request.Method == "POST" {
		for _, path := range anonymousPaths {
			if path == request.URL.Path {
				return true
			}
		}
	}

	return false
}

func extractToken(tokenHeader string) string {
	parts := strings.Split(tokenHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func getLoggedInUser(ctx *gin.Context) (*LoggedInUser, error) {
	value, ok := ctx.Get("loggedInUser")
	if !ok {
		return nil, errors.New("logged in user not found")
	}

	user, ok := value.(LoggedInUser)
	if !ok {
		return nil, errors.New("logged in user not found")
	}

	return &user, nil
}
