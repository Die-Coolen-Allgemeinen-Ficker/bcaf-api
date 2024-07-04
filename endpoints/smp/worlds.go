package smp

import (
	"bcaf-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Worlds(path string, rest *gin.Engine, mongoClient *mongo.Client) {
	rest.OPTIONS(path, func(ctx *gin.Context) {})

	rest.GET(path, func(ctx *gin.Context) {
		accessToken := ctx.Request.Header.Get("authorization")
		userId := util.Validate(accessToken, ctx, true)
		if userId == nil {
			return 
		}

		data, err := util.GetData[util.SmpWorld]("smp-worlds", bson.D{}, ctx, mongoClient)
		if err != nil {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"response": data,
		})
	})
}