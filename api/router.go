package api

import (
	"auth/domain/logs"
	"auth/domain/users"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AppContext struct {
	Port        int
	JwtSecret   []byte
	UserService *users.UserService
	LogService  *logs.LogService
}

func Run(appCtx *AppContext) {
	router := gin.Default()
	router.Use(authorizeRoutes(appCtx))
	router.POST("/authenticate", authenticateUser(appCtx))
	router.POST("/users", createUser(appCtx))
	router.GET("/users/:userId", getUser(appCtx))
	router.GET("/users/:userId/logs", getLogs(appCtx))

	router.Run(":" + strconv.Itoa(appCtx.Port))
}
