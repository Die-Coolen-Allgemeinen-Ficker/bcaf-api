package lookup

import (
	"bcaf-api/util"
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func isValidKey(key string) bool {
	var validMutableKeys = []string{"color", "backgroundImageUrl", "foregroundImageUrl"}
	for _, element := range validMutableKeys {
		if key == element {
			return true
		}
	}
	return false
}

func Id(path string, rest *gin.Engine, mongoClient *mongo.Client) {
	rest.OPTIONS(path, func(ctx *gin.Context) {})

	rest.GET(path, func(ctx *gin.Context) {
		// Validate
		searchId := ctx.Param("id")
		accessToken := ctx.Request.Header.Get("authorization")
		userId := util.Validate(accessToken, ctx)
		if userId == nil {
			return
		}

		// Get Data

		results, err := util.GetData("accounts", bson.D{{Key: "userId", Value: searchId}}, ctx, mongoClient)
		if err != nil {
			return
		}

		if len(results) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{
				"response": "user not found",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"response": results[0],
		})
	})

	rest.POST(path, func(ctx *gin.Context) {
		// Validate
		searchId := ctx.Param("id")
		match, _ := regexp.MatchString("^[0-9]+$", searchId)
		if !match || len(searchId) < 17 || len(searchId) > 18 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "invalid user id",
			})
			return
		}

		accessToken := ctx.Request.Header.Get("authorization")
		userId := util.Validate(accessToken, ctx)
		if userId == nil {
			return
		}

		if searchId != *userId {
			ctx.JSON(http.StatusForbidden, gin.H{
				"response": "you may only edit your own user data",
			})
			return
		}

		// Get data

		results, err := util.GetData("accounts", bson.D{{Key: "userId", Value: searchId}}, ctx, mongoClient)
		if err != nil {
			return
		}

		if len(results) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{
				"response": "user not found",
			})
			return
		}

		var bodyData map[string]interface{}
		err = json.NewDecoder(ctx.Request.Body).Decode(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "invalid body",
			})
			return
		}

		// Check if attemped modification is valid

		user := results[0]

		for key := range bodyData {
			if !isValidKey(key) {
				hasAchievement := false
				for _, achievement := range user.Profile.Achievements {
					if achievement.Name == "Hackerman" {
						hasAchievement = true
						break
					}
				}
				if !hasAchievement {
					user.Profile.Achievements = append(user.Profile.Achievements, struct{Name string "json:\"name\""; Description string "json:\"description\""; Timestamp int64 "json:\"timestamp\""}{Name: "Hackerman", Description: "Schicke eine POST Request an die BCAF REST API um Daten zu verändern, die du nicht verändern darfst.", Timestamp: time.Now().UnixMilli()})
					err = util.UpdateData(user, ctx, mongoClient.Database("bcaf-user-data").Collection("accounts"))
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{
							"response": "internal server error",
						})
						return
					}

					ctx.JSON(http.StatusBadRequest, gin.H{
						"response": "invalid body (nice try +hackerman achievement)",
					})
					return
				}

				ctx.JSON(http.StatusBadRequest, gin.H{
					"response": "invalid body",
				})
				return
			}
		}

		// Modify data

		var ok bool
		if bodyData["color"] != nil {
			var color string
			color, ok = bodyData["color"].(string)
			user.Profile.Color = color
		}
		if bodyData["backgroundImageUrl"] != nil {
			var backgroundImageUrl string
			backgroundImageUrl, ok = bodyData["backgroundImageUrl"].(string)
			user.Profile.BackgroundImageUrl = backgroundImageUrl
		}
		if bodyData["foregroundImageUrl"] != nil {
			var foregroundImageUrl string
			foregroundImageUrl, ok = bodyData["foregroundImageUrl"].(string)
			user.Profile.ForegroundImageUrl = foregroundImageUrl
		}
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "invalid property type",
			})
			return
		}

		// Update data

		err = util.UpdateData(user, ctx, mongoClient.Database("bcaf-user-data").Collection("accounts"))
		if err != nil {
			return
		}

		response, _ := util.ToJSON(*user)
		ctx.JSON(http.StatusOK, gin.H{
			"response": response,
		})
	})
}