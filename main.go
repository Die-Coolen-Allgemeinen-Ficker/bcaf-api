package main

import (
	"bcaf-api/endpoints"
	"bcaf-api/endpoints/accounts"
	"bcaf-api/endpoints/accounts/lookup"
	"bcaf-api/endpoints/minecraft/name"
	"bcaf-api/endpoints/smp"
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")

		ctx.Next()
	}
}

func main() {
	godotenv.Load()

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_CONNECTION_STRING")))

	if err != nil {
		panic(err)
	}
	
	gin.SetMode(gin.ReleaseMode)
	rest := gin.Default()
	rest.SetTrustedProxies(nil)

	rest.Use(corsMiddleware())

	accounts.List("/v1/accounts/list", rest, mongoClient)
	accounts.Auth("/v1/accounts/auth", rest, mongoClient)
	accounts.Refresh("/v1/accounts/refresh", rest)
	lookup.Id("/v1/accounts/lookup/:id", rest, mongoClient)
	endpoints.Ping("/v1/ping", rest)
	name.Uuid("/v1/minecraft/name/:uuid", rest)
	smp.Info("/v1/smp/info", rest, mongoClient)
	smp.Worlds("/v1/smp/worlds", rest, mongoClient)

	rest.Run()
}
