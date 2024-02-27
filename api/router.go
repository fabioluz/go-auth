package api

import (
	"auth/domain/users"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AppContext struct {
	Port        int
	JwtSecret   []byte
	UserService *users.UserService
}

func Run(appCtx *AppContext) {
	router := gin.Default()
	router.Use(authorizeRoutes(appCtx))
	router.POST("/authenticate", authenticateUser(appCtx))
	router.POST("/users", createUser(appCtx))
	router.GET("/users/:userId", getUser(appCtx))

	router.Run(":" + strconv.Itoa(appCtx.Port))
}
