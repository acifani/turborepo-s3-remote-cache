package main

import (
	"net/http"
	"strings"

	"github.com/acifani/turborepo-s3-remote-cache/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

func main() {
	config := config.Read()
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.SetTrustedProxies(config.TrustedProxies)

	v8 := router.Group("/v8", authTokenMiddleware(config.TurboAllowedTokens), configMiddleware(config), awsMiddleware(config))
	{
		v8.GET("/artifacts/:hash", getArtifactHandler)
		v8.PUT("/artifacts/:hash", putArtifactHandler)
	}

	private := router.Group("/_/")
	{
		private.GET("/status", statusHandler)
	}

	err := router.Run()
	if err != nil {
		panic(err)
	}
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

	uploader := ctx.MustGet("uploader").(*s3manager.Uploader)
	config := ctx.MustGet("config").(*config.Config)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &config.AWSS3Bucket,
		Key:    getArtifactPath(teamId, slug, hash),
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

	buffer := aws.NewWriteAtBuffer([]byte{})
	downloader := ctx.MustGet("downloader").(*s3manager.Downloader)
	config := ctx.MustGet("config").(*config.Config)
	_, err := downloader.Download(buffer, &s3.GetObjectInput{
		Bucket: &config.AWSS3Bucket,
		Key:    getArtifactPath(teamId, slug, hash),
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

func awsMiddleware(config *config.Config) gin.HandlerFunc {
	session := createAwsSession(config)
	return func(ctx *gin.Context) {
		ctx.Set("awsSession", session)
		ctx.Set("downloader", s3manager.NewDownloader(session))
		ctx.Set("uploader", s3manager.NewUploader(session))
		ctx.Next()
	}
}

func createAwsSession(config *config.Config) *session.Session {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:           &config.AWSRegion,
		Endpoint:         &config.AWSEndpoint,
		DisableSSL:       &config.AWSDisableSSL,
		S3ForcePathStyle: &config.AWSS3ForcePathStyle,
		Credentials: credentials.NewStaticCredentials(
			config.AWSAccessKeyID,
			config.AWSSecretAccessKey,
			"",
		),
	}))

	return sess
}

func configMiddleware(config *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("config", config)
		ctx.Next()
	}
}

func getArtifactPath(teamId, slug, hash string) *string {
	bucketDirectory := teamId
	if slug != "" {
		bucketDirectory = slug
	}

	if bucketDirectory == "" {
		return &hash
	}

	return aws.String(bucketDirectory + "/" + hash)
}
