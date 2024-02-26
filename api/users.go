package api

import (
	"auth/domain/users"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createUser(app *AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var request users.CreateUser
		if err := ctx.BindJSON(&request); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		_, err := app.UserService.CreateUser(request)
		if err != nil {
			var userErr *users.CreateUserError
			if errors.As(err, &userErr) {
				ctx.JSON(http.StatusUnprocessableEntity, userErr)
				return
			}

			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		response := gin.H{
			"message": "User created. Authorization required to view information.",
		}
		ctx.JSON(http.StatusCreated, response)
	}
}
