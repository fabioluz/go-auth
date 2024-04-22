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

func (server *Server) authenticateUser(ctx *gin.Context) {
	var request users.AuthenticateUserInput
	if err := ctx.BindJSON(&request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := server.UserService.Authenticate(request)
	if err != nil {
		var userErr *users.AuthenticateUserError
		if errors.As(err, &userErr) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := generateToken(server, user)
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
