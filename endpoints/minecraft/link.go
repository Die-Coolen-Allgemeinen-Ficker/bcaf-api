package minecraft

import (
	"bcaf-api/util"
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Link(path string, rest *gin.Engine, mongoClient *mongo.Client) {
	rest.OPTIONS(path, func(ctx *gin.Context) {})

	rest.POST(path, func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotImplemented, gin.H{})
		return

		// Check if code is given
		var bodyData map[string]interface{}
		err := json.NewDecoder(ctx.Request.Body).Decode(&bodyData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "invalid body",
			})
			return
		}
		code, exists := bodyData["code"].(string)
		if !exists {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"response": "missing code",
			})
			return
		}

		// Validate
		accessToken := ctx.Request.Header.Get("authorization")
		userId := util.Validate(accessToken, ctx, true)
		if userId == nil {
			return 
		}

		// Get Microsoft access token
		var params = []byte("client_id=" + os.Getenv("MICROSOFT_CLIENT_ID") + "&client_secret=" + os.Getenv("MICROSOFT_CLIENT_SECRET") + "&grant_type=authorization_code&code=" + code + "&redirect_uri=" + os.Getenv("MICROSOFT_REDIRECT_URI") + "&scope=XboxLive.signin")
		response, err := http.Post("https://login.microsoftonline.com/consumers/oauth2/v2.0/token", "application/x-www-form-urlencoded", bytes.NewBuffer(params))
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

		microsoftAccessToken, exists := responseBody["access_token"].(string)
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "couldnt get access token",
			})
			return
		}

		// Get Xbox Live access token
		httpClient := &http.Client{}

		postBody := map[string]interface{}{
			"Properties": map[string]interface{}{
				"AuthMethod": "RPS",
				"SiteName": "user.auth.xboxlive.com",
				"RpsTicket": "d=" + microsoftAccessToken,
			},
			"RelyingParty": "http://auth.xboxlive.com",
			"TokenType": "JWT",
		}
		postBodyBytes, _ := json.Marshal(postBody)
		request, _ := http.NewRequest("POST", "https://user.auth.xboxlive.com/user/authenticate", bytes.NewBuffer(postBodyBytes))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")

		response, err = httpClient.Do(request)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "failed to make request",
			})
			return
		}
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "couldnt parse response body",
			})
			return
		}

		xboxLiveToken, exists := responseBody["Token"].(string)
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "couldnt get xbox live token",
			})
			return
		}
		xboxUserHash := responseBody["DisplayClaims"].(map[string]interface{})["xui"].([]interface{})[0].(map[string]interface{})["uhs"].(string)
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "couldnt get xbox user hash",
			})
			return
		}

		// Get XSTS token
		postBody = map[string]interface{}{
			"Properties": map[string]interface{}{
				"SandboxId": "RETAIL",
				"UserTokens": []string{
					xboxLiveToken,
				},
			},
			"RelyingParty": "rp://api.minecraftservices.com/",
			"TokenType": "JWT",
		}
		postBodyBytes, _ = json.Marshal(postBody)
		request, _ = http.NewRequest("POST", "https://xsts.auth.xboxlive.com/xsts/authorize", bytes.NewBuffer(postBodyBytes))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")

		response, err = httpClient.Do(request)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "failed to make request",
			})
			return
		}
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "couldnt parse response body",
			})
			return
		}

		xstsToken, exists := responseBody["Token"].(string)
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "couldnt get xsts token",
			})
			return
		}

		// Get Minecraft access token
		postBody = map[string]interface{}{
			"identityToken": "XBL3.0 x=" + xboxUserHash + ";" + xstsToken,
		}
		postBodyBytes, _ = json.Marshal(postBody)
		request, _ = http.NewRequest("POST", "https://api.minecraftservices.com/authentication/login_with_xbox", bytes.NewBuffer(postBodyBytes))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")

		response, err = httpClient.Do(request)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "failed to make request",
			})
			return
		}
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "couldnt parse response body",
			})
			return
		}

		minecraftAccessToken, exists := responseBody["access_token"].(string)
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "couldnt get minecraft access token",
			})
			return
		}

		// FINALLY make the minecraft profile request
		request, _ = http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile", nil)
		request.Header.Set("Authorization", "Bearer " + minecraftAccessToken)
		request.Header.Set("Accept", "application/json")

		response, err = httpClient.Do(request)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "failed to make request",
			})
			return
		}
		defer response.Body.Close()
		err = json.NewDecoder(response.Body).Decode(&responseBody)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "couldnt parse response body",
			})
			return
		}
		/*minecraftName, exists := responseBody["name"].(string)
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"response": "couldnt get minecraft name",
			})
			return
		}*/

		ctx.JSON(http.StatusOK, gin.H{
			"debug": responseBody,
			"response": minecraftAccessToken,
		})
	})
}