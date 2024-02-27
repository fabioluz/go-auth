package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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

type LoggedInUser struct {
	ID    string
	Email string
}

func authorizeRoutes(appCtx *AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			if !isAnonymousRequest(ctx.Request) {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			ctx.Next()
			return
		}

		if isAnonymousRequest(ctx.Request) {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

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
