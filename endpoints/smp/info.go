package smp

import (
	"bcaf-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Info(path string, rest *gin.Engine, mongoClient *mongo.Client) {
	rest.OPTIONS(path, func(ctx *gin.Context) {})

	rest.GET(path, func(ctx *gin.Context) {
		data, err := util.GetData[util.Smp]("smps", bson.D{}, ctx, mongoClient)
		if err != nil {
			return
		}

		accessToken := ctx.Request.Header.Get("authorization")
		userId := util.Validate(accessToken, ctx, false)
		if userId == nil {
			for _, smp := range data {
				smp.Ip = nil
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"response": data,
		})
	})
}