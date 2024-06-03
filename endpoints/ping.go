package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(path string, rest *gin.Engine) {
	rest.OPTIONS(path, func(ctx *gin.Context) {})

	rest.GET(path, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"response": "hai :3",
		})
	})
}