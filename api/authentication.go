package api

import (
	"auth/domain/users"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticateUserOutput struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func authenticateUser(app *AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var request users.AuthenticateUserInput
		if err := ctx.BindJSON(&request); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		user, err := app.UserService.AuthenticateUser(request)
		if err != nil {
			var userErr *users.AuthenticateUserError
			if errors.As(err, &userErr) {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		token, err := generateToken(app, user)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		output := AuthenticateUserOutput{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Token: token,
		}

		ctx.JSON(http.StatusOK, output)
	}
}
