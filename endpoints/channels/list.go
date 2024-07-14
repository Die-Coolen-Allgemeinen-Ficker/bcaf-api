package channels

import (
	"bcaf-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func List(path string, rest *gin.Engine, mongoClient *mongo.Client) {
	rest.OPTIONS(path, func(ctx *gin.Context) {})

	rest.GET(path, func(ctx *gin.Context) {
		hideArchived, exists := ctx.GetQuery("hideArchived")
		if !exists {
			hideArchived = "true"
		} else if !(hideArchived == "true" || hideArchived == "false") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "invalid query param hideArchived",
			})
			return
		}

		accessToken := ctx.Request.Header.Get("authorization")
		userId := util.Validate(accessToken, ctx, true)
		if userId == nil {
			return
		}

		filter := bson.D{}
		if hideArchived == "true" {
			filter = append(filter, bson.E{Key: "archived", Value: false})
		}
		data, err := util.GetData[util.Channel]("channels", filter, ctx, mongoClient)
		if err != nil {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"response": data,
		})
	})
}