package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getPageSizeParam(ctx *gin.Context) int {
	pageSizeParam := ctx.Param("pageSize")
	if pageSizeParam == "" {
		return 10
	}

	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil {
		return 10
	}

	if pageSize > 100 {
		return 100
	}

	return pageSize
}

func getAfterParam(ctx *gin.Context) string {
	return ctx.Param("after")
}

func (server *Server) getLogs(ctx *gin.Context) {
	loggedInUser, err := getLoggedInUser(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if ctx.Param("userId") != loggedInUser.ID {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	pageSize := getPageSizeParam(ctx)
	after := getAfterParam(ctx)

	logs, err := server.LogService.Get(loggedInUser.ID, pageSize, after)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, logs)
}
