package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

func main() {
	allowedTokensEnv := os.Getenv("TURBOREPO_ALLOWED_TOKENS")
	allowedTokens := strings.Split(allowedTokensEnv, ",")

	awsRegion := os.Getenv("AWS_REGION")
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	awsAccessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsSession := createAwsSession(awsRegion, awsAccessKeyId, awsSecretAccessKey, awsEndpoint)

	router := gin.Default()

	v8 := router.Group("/v8", authTokenMiddleware(allowedTokens), awsMiddleware(awsSession))
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
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func putArtifactHandler(ctx *gin.Context) {
	teamId := ctx.Query("teamId")
	slug := ctx.Query("slug")
	hash := ctx.Param("hash")

	if teamId == "" && slug == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "teamId or slug are required",
		})
		return
	}

	bucketDirectory := teamId
	if slug != "" {
		bucketDirectory = slug
	}

	uploader := ctx.MustGet("uploader").(*s3manager.Uploader)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("turborepo-cache"),
		Key:    aws.String(bucketDirectory + "/" + hash),
		Body:   ctx.Request.Body,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func getArtifactHandler(ctx *gin.Context) {
	teamId := ctx.Query("teamId")
	slug := ctx.Query("slug")
	hash := ctx.Param("hash")

	if teamId == "" && slug == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "teamId or slug are required",
		})
		return
	}

	bucketDirectory := teamId
	if slug != "" {
		bucketDirectory = slug
	}

	buffer := aws.NewWriteAtBuffer([]byte{})
	downloader := ctx.MustGet("downloader").(*s3manager.Downloader)
	_, err := downloader.Download(buffer, &s3.GetObjectInput{
		Bucket: aws.String("turborepo-cache"),
		Key:    aws.String(bucketDirectory + "/" + hash),
	})

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())
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

func createAwsSession(awsRegion, awsAccessKeyId, awsSecretAccessKey, awsEndpoint string) *session.Session {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   &awsRegion,
		Endpoint: &awsEndpoint,
		Credentials: credentials.NewStaticCredentials(
			awsAccessKeyId,
			awsSecretAccessKey,
			"",
		),
	}))

	return sess
}

func awsMiddleware(session *session.Session) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("awsSession", session)
		ctx.Set("downloader", s3manager.NewDownloader(session))
		ctx.Set("uploader", s3manager.NewUploader(session))
		ctx.Next()
	}
}
