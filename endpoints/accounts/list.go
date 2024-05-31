package accounts

import (
	"bcaf-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func List(path string, rest *gin.Engine, mongoClient *mongo.Client) {
	rest.GET(path, func(ctx *gin.Context) {
		// Validate
		accessToken := ctx.Request.Header.Get("authorization")
		userId := util.Validate(accessToken, ctx)
		if userId == nil {
			return
		}

		// Get data

		results, err := util.GetData("accounts", bson.D{}, ctx, mongoClient)
		if err != nil {
			return
		}

		if results == nil {
			return
		}

		var response []map[string]interface{}
		for _, element := range results {
			json, _ := util.ToJSON(*element)
			response = append(response, json)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"response": response,
		})
	})
}