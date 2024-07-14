package ngram

import (
	"bcaf-api/util"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Search(path string, rest *gin.Engine, mongoClient *mongo.Client) {
	rest.OPTIONS(path, func(ctx *gin.Context) {})

	rest.GET(path, func(ctx *gin.Context) {
		// Validate
		accessToken := ctx.Request.Header.Get("authorization")
		userId := util.Validate(accessToken, ctx, true)
		if userId == nil {
			return
		}

		// Get params
		query, queryExists := ctx.GetQuery("query")
		if !queryExists {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "missing query",
			})
			return
		}
		var after int64
		afterQuery, afterExists := ctx.GetQuery("after")
		if afterExists {
			var err error
			after, err = strconv.ParseInt(afterQuery, 10, 64)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"response": "invalid query param after",
				})
				return
			}
		}
		var before int64
		beforeQuery, beforeExists := ctx.GetQuery("before")
		if beforeExists {
			var err error
			before, err = strconv.ParseInt(beforeQuery, 10, 64)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"response": "invalid query param before",
				})
				return
			}
		}
		authorId, authorExists := ctx.GetQuery("author")
		channelId, channelExists := ctx.GetQuery("channel")
		interval, intervalExists := ctx.GetQuery("interval")
		if !intervalExists {
			interval = "daily"
		} else if !(interval == "daily" || interval == "weekly") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "invalid query param interval",
			})
			return
		}

		// Search
		filter := bson.D{{Key: "content", Value: bson.D{{Key: "$regex", Value: query}, {Key: "$options", Value: "i"}}}}
		if afterExists {
			filter = append(filter, bson.E{Key: "createdTimestamp", Value: bson.D{{Key: "$gte", Value: after}}})
		}
		if beforeExists {
			filter = append(filter, bson.E{Key: "createdTimestamp", Value: bson.D{{Key: "$lte", Value: before}}})
		}
		if authorExists {
			filter = append(filter, bson.E{Key: "authorId", Value: authorId})
		}
		if channelExists {
			filter = append(filter, bson.E{Key: "channelId", Value: channelId})
		}

		matchedMessages, err := util.GetData[util.Message]("messages", filter, ctx, mongoClient)
		if err != nil {
			return
		}
		messageCounts, err := util.GetData[util.MessageCount]("message-counts", bson.D{{Key: "_id", Value: interval}}, ctx, mongoClient)
		if err != nil {
			return
		}

		matchedMessageCounts := map[string]int64{}
		for _, message := range matchedMessages {
			var timestamp string
			if interval == "daily" {
				timestamp = strconv.FormatInt(message.CreatedTimestamp - (message.CreatedTimestamp % 86400000), 10);
			} else {
				timestamp = strconv.FormatInt(message.CreatedTimestamp - (message.CreatedTimestamp % 604800000), 10);
			}
			matchedMessageCounts[timestamp]++;
		}
		relative := map[string]float64{}
		absolute := map[string]int64{}
		messageCount := messageCounts[0].Counts
		messageCountInterval := map[string]struct{Count int64 "json:\"count\""; Characters int "json:\"characters\""}{}
		if afterExists || beforeExists {
			for timestamp, count := range messageCount {
				timestampNumber, _ := strconv.ParseInt(timestamp, 10, 64)
				if afterExists && beforeExists && !(timestampNumber >= after && timestampNumber <= before) {
					continue;
				}
				if afterExists && timestampNumber < after {
					continue;
				}
				if beforeExists && timestampNumber > before {
					continue;
				}
				messageCountInterval[timestamp] = count
			}
		} else {
			messageCountInterval = messageCount
		}

		for timestamp, count := range messageCountInterval {
			intervalCount, exists := matchedMessageCounts[timestamp]
			if !exists {
				intervalCount = 0
			}
			relative[timestamp] = float64(len(query)) * float64(intervalCount) / float64(count.Characters)
			if math.IsNaN(relative[timestamp]) {
				relative[timestamp] = 0.0
			}
			absolute[timestamp] = intervalCount
		}

		ctx.JSON(http.StatusOK, gin.H{
			"response": gin.H{
				"relative": relative,
				"absolute": absolute,
			},
		})
	})
}