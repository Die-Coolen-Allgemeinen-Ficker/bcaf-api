package util

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetData[T any](collection string, filter bson.D, ctx *gin.Context, mongoClient *mongo.Client) ([]*T, error) {
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

	var results []*T
	for cursor.Next(context.TODO()) {
		var r bson.M
		cursor.Decode(&r)
		if r["_hidden"] == false {
			result, _ := FromJSONRaw[T](r)
			results = append(results, result)
		}
	}

	return results, nil
}