package util

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Validate(accessToken string, ctx *gin.Context) *string {
	httpClient := &http.Client{}

	// Check if user is BCAF member

	guildRequest, err := http.NewRequest("GET", "https://discord.com/api/users/@me/guilds", nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"response": "internal server error",
		})
		return nil
	}
	guildRequest.Header.Set("authorization", "Bearer " + accessToken)
	guildResponse, err := httpClient.Do(guildRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"response": "failed to make request",
		})
		return nil
	}

	defer guildResponse.Body.Close()

	rawBody, err := io.ReadAll(guildResponse.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"response": "internal server error",
		})
		return nil
	}

	var bodyData []map[string]interface{}
	json.Unmarshal(rawBody, &bodyData)
	if bodyData == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"response": "invalid access token",
		})
		return nil
	}

	member := false
	for _, guild := range bodyData {
		if guild["id"] == "555729962188144660" {
			member = true
			break
		}
	}
	if !member {
		ctx.JSON(http.StatusForbidden, gin.H{
			"response": "you are not a bcaf member",
		})
		return nil
	}

	// Get user id

	userRequest, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"response": "internal server error",
		})
		return nil
	}
	userRequest.Header.Set("authorization", "Bearer " + accessToken)
	userResponse, err := httpClient.Do(userRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"response": "failed to make request",
		})
		return nil
	}

	defer userResponse.Body.Close()

	var userData map[string]interface{}
	err = json.NewDecoder(userResponse.Body).Decode(&userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"response": "couldnt parse response body",
		})
		return nil
	}

	userId := userData["id"].(string)
	return &userId
}