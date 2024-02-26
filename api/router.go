package api

import (
	"auth/domain/users"

	"github.com/gin-gonic/gin"
)

type AppContext struct {
	Port        int
	UserService *users.UserService
}

func Run(appCtx *AppContext) {
	router := gin.Default()
	router.POST("/users", createUser(appCtx))

	router.Run(":8080")
}
