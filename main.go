package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	allowedTokensEnv := os.Getenv("TUBOREPO_ALLOWED_TOKENS")
	allowedTokens := strings.Split(allowedTokensEnv, ",")

	router := gin.Default()

	v8 := router.Group("/v8", authTokenMiddleware(allowedTokens))
	{
		v8.GET("/artifacts/:hash", getArtifactHandler)
		v8.PUT("/artifacts/:hash", putArtifactHandler)
	}

	private := router.Group("/_/")
	{
		private.GET("/status", statusHandler)
	}

	router.Run()
}

func statusHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "ok",
	})
}

func putArtifactHandler(ctx *gin.Context) {
	teamId := ctx.Query("teamId")
	slug := ctx.Query("slug")
	hash := ctx.Param("hash")

	ctx.JSON(200, gin.H{
		"teamId": teamId,
		"slug":   slug,
		"hash":   hash,
	})
}

func getArtifactHandler(ctx *gin.Context) {
	teamId := ctx.Query("teamId")
	slug := ctx.Query("slug")
	hash := ctx.Param("hash")

	ctx.JSON(200, gin.H{
		"teamId": teamId,
		"slug":   slug,
		"hash":   hash,
	})
}

func authTokenMiddleware(allowedTokens []string) gin.HandlerFunc {
	tokenMap := make(map[string]bool)
	for _, token := range allowedTokens {
		tokenMap[token] = true
	}

	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "missing Authorization header",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			token := parts[1]
			if tokenMap[token] {
				ctx.Next()
				return
			}

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid Authorization header format",
		})
	}
}
