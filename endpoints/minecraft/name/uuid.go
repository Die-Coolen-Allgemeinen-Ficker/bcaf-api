package name

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Uuid(path string, rest *gin.Engine) {
	rest.OPTIONS(path, func(ctx *gin.Context) {})

	rest.GET(path, func(ctx *gin.Context) {
		uuid := ctx.Param("uuid")

		response, err := http.Get("https://sessionserver.mojang.com/session/minecraft/profile/" + uuid)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "failed to make request",
			})
			return
		}

		defer response.Body.Close()
		var responseBody map[string]interface{}
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "couldnt parse response body",
			})
			return
		}

		if responseBody["errorMessage"] != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "invalid uuid",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"response": responseBody["name"].(string),
		})
	})
}