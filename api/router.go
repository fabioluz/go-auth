package api

import (
	"auth/domain/logs"
	"auth/domain/users"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Port        int
	JwtSecret   []byte
	UserService *users.UserService
	LogService  *logs.LogService
}

func (server *Server) Run() {
	router := gin.Default()
	router.Use(server.authorizeRoutes)
	router.POST("/authenticate", server.authenticateUser)
	router.POST("/users", server.createUser)
	router.GET("/users/:userId", server.getUser)
	router.PATCH("/users/:userId", server.updateUser)
	router.GET("/users/:userId/logs", server.getLogs)

	router.Run(":" + strconv.Itoa(server.Port))
}
