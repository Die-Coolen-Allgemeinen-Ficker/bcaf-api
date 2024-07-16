package util

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var tokenCache = map[string]struct{UserId string; Timeout int64}{}

func Validate(accessToken string, ctx *gin.Context, denyIfInvalid bool) *string {
	for token, cache := range tokenCache {
		if time.Now().UnixMilli() >= cache.Timeout {
			delete(tokenCache, token)
		} else if accessToken == token {
			return &cache.UserId
		}
	}

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
		if denyIfInvalid {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "invalid access token",
			})
		}
		return nil
	} else {
		member := false
		for _, guild := range bodyData {
			if guild["id"] == "555729962188144660" {
				member = true
				break
			}
		}
		if !member {
			if denyIfInvalid {
				ctx.JSON(http.StatusForbidden, gin.H{
					"response": "you are not a bcaf member",
				})
			}
			return nil
		}
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
	tokenCache[accessToken] = struct{UserId string; Timeout int64}{
		UserId: userId,
		Timeout: time.Now().UnixMilli() + time.Minute.Milliseconds(),
	}
	return &userId
}