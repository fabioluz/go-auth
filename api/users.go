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

		_, err := app.UserService.Create(request)
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

func getUser(appCtx *AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		loggedInUser, err := getLoggedInUser(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		if ctx.Param("userId") != loggedInUser.ID {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		user, err := appCtx.UserService.Get(loggedInUser.ID)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}

func updateUser(appCtx *AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		loggedInUser, err := getLoggedInUser(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		if ctx.Param("userId") != loggedInUser.ID {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		var input users.UpdateUser
		if err := ctx.BindJSON(&input); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err = appCtx.UserService.Update(loggedInUser.ID, input)
		if err != nil {
			var userErr *users.UpdateUserError
			if errors.As(err, &userErr) {
				ctx.JSON(http.StatusUnprocessableEntity, userErr)
				return
			}

			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}
