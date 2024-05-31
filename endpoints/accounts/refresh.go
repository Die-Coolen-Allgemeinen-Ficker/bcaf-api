package accounts

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
)

func Refresh(path string, rest *gin.Engine) {
	rest.OPTIONS(path, func(ctx *gin.Context) {})

	rest.GET(path, func(ctx *gin.Context) {
		godotenv.Load()

		// Get access token from discord

		refreshToken := ctx.Request.Header.Get("authorization")

		var params = []byte("client_id=" + os.Getenv("CLIENT_ID") + "&client_secret=" + os.Getenv("CLIENT_SECRET") + "&grant_type=refresh_token&refresh_token=" + refreshToken)
		response, err := http.Post("https://discord.com/api/oauth2/token", "application/x-www-form-urlencoded", bytes.NewBuffer(params))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "no response from discord",
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

		if responseBody["error"] != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "invalid refresh token",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"response": responseBody,
		})
	})
}