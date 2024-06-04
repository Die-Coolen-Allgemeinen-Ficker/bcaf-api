package accounts

import (
	"bcaf-api/util"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Auth(path string, rest *gin.Engine, mongoClient *mongo.Client) {
	rest.OPTIONS(path, func(ctx *gin.Context) {})

	rest.GET(path, func(ctx *gin.Context) {
		godotenv.Load()

		// Get access token from discord

		code := ctx.Request.Header.Get("authorization")

		var params = []byte("client_id=" + os.Getenv("CLIENT_ID") + "&client_secret=" + os.Getenv("CLIENT_SECRET") + "&grant_type=authorization_code&code=" + code + "&redirect_uri=" + os.Getenv("REDIRECT_URI") + "&scope=identify+guilds")
		response, err := http.Post("https://discord.com/api/oauth2/token", "application/x-www-form-urlencoded", bytes.NewBuffer(params))
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

		// Validate

		if responseBody["error"] != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"response": "invalid code",
			})
			return
		}
		accessToken := responseBody["access_token"].(string)

		userId := util.Validate(accessToken, ctx)
		if userId == nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"response": "you are not a bcaf member",
			})
			return
		}

		// Create new account entry in database if no account exists

		results, err := util.GetData("accounts", bson.D{{Key: "userId", Value: userId}}, ctx, mongoClient)
		if err != nil {
			return
		}

		// Get user data

		httpClient := &http.Client{}

		userRequest, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "internal server error",
			})
			return
		}
		userRequest.Header.Set("authorization", "Bearer " + accessToken)
		userResponse, err := httpClient.Do(userRequest)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "failed to make request",
			})
			return
		}

		defer userResponse.Body.Close()

		var userData map[string]interface{}
		err = json.NewDecoder(userResponse.Body).Decode(&userData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "couldnt parse response body",
			})
			return
		}

		if len(results) == 0 {
			// Init new account

			var newAccount util.AccountData
			newAccount.UserId = *userId
			newAccount.Name = userData["global_name"].(string)
			newAccount.AvatarUrl = "https://cdn.discordapp.com/avatars/" + *userId + "/" + userData["avatar"].(string)
			newAccount.Profile.Level = 0
			newAccount.Profile.Color = "#000000"
			newAccount.Profile.BackgroundImageUrl = "https://die-coolen-allgemeinen-ficker.github.io/assets/images/wallpapers/3.png"
			newAccount.Profile.ForegroundImageUrl = ""
			newAccount.Profile.SocialCredit.Amount = 1000
			newAccount.Profile.SocialCredit.Tier = "A"
			newAccount.Profile.Games.SnakeHighscore = 0
			newAccount.Profile.Games.TictactoeWins = 0
			newAccount.Profile.MessageStats.NWordCount = 0
			newAccount.Profile.MessageStats.MessageCount = 0
			newAccount.Profile.MessageStats.MessagesLast30Days = 0
			newAccount.Profile.Achievements = []struct{Name string "json:\"name\""; Description string "json:\"description\""}{}
			newAccount.BcafCoin = 0
			newAccount.HasBoostedBefore = false
			newAccount.HasPlayedLeagueOfLegends = false
			newAccount.BcafJoinTimestamp = 0
			newAccount.Legacy = false
			newAccount.CreatedTimestamp = time.Now().UnixMilli()
			newAccount.UpdatedTimestamp = time.Now().UnixMilli()

			jsonData, err := util.ToJSON(newAccount)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"response": "failed to create account",
				})
				return
			}

			col := mongoClient.Database("bcaf-user-data").Collection("accounts")
			_, err = col.InsertOne(context.TODO(), jsonData)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"response": "failed to create account",
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"response": gin.H{
				"token": responseBody,
				"user": userData,
			},
		})
	})
}