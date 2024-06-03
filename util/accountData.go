package util

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountData struct {
	UserId    string `json:"userId"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatarUrl"`
	Profile   struct {
		Level              float64 `json:"level"`
		Color              string  `json:"color"`
		BackgroundImageUrl string  `json:"backgroundImageUrl"`
		ForegroundImageUrl string  `json:"foregroundImageUrl"`
		MinecraftUuid      *string `json:"minecraftUuid"`
		SocialCredit       struct {
			Amount int64  `json:"amount"`
			Tier   string `json:"tier"`
		} `json:"socialCredit"`
		Games struct {
			SnakeHighscore int64 `json:"snakeHighscore"`
			TictactoeWins  int64 `json:"tictactoeWins"`
		} `json:"games"`
		MessageStats struct {
			NWordCount         int64 `json:"nWordCount"`
			BReactionCount     int64 `json:"bReactionCount"`
			MessageCount       int64 `json:"messageCount"`
			MessagesLast30Days int64 `json:"messagesLast30Days"`
		} `json:"messageStats"`
		Achievements []struct {
			Name string `json:"name"`
			Description string `json:"description"`
		} `json:"achievements"`
	} `json:"profile"`
	BcafCoin                 int64 `json:"bcafCoin"`
	HasBoostedBefore         bool  `json:"hasBoostedBefore"`
	HasPlayedLeagueOfLegends bool  `json:"hasPlayerLeagueOfLegends"`
	BcafJoinTimestamp        int64 `json:"bcafJoinTimestamp"`
	Legacy                   bool  `json:"legacy"`
	CreatedTimestamp         int64 `json:"createdTimestamp"`
	UpdatedTimestamp         int64 `json:"updatedTimestamp"`
}

func GetData(collection string, filter bson.D, ctx *gin.Context, mongoClient *mongo.Client) ([]*AccountData, error) {
	col := mongoClient.Database("bcaf-user-data").Collection(collection)
	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		switch err {
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "internal server error",
			})
		}
		return nil, err
	}

	var results []*AccountData
	for cursor.Next(context.TODO()) {
		var r bson.M
		cursor.Decode(&r)
		result, _ := FromJSONRaw[AccountData](r)
		results = append(results, result)
	}

	return results, nil
}

func UpdateData(account *AccountData, ctx *gin.Context, collection *mongo.Collection) error {
	account.UpdatedTimestamp = time.Now().UnixMilli()
	jsonData, _ := ToJSON(*account)
	_, err := collection.ReplaceOne(context.TODO(), bson.D{{Key: "userId", Value: account.UserId}}, jsonData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"response": "internal server error",
		})
		return err
	}
	return nil
}