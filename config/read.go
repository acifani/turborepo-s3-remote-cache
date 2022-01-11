package config

import (
	"os"
	"strings"
)

type Config struct {
	AWSRegion           string
	AWSEndpoint         string
	AWSAccessKeyID      string
	AWSSecretAccessKey  string
	AWSDisableSSL       bool
	AWSS3Bucket         string
	AWSS3ForcePathStyle bool
	TurboAllowedTokens  []string
}

func Read() *Config {
	allowedTokensEnv := os.Getenv("TURBOREPO_ALLOWED_TOKENS")
	allowedTokens := strings.Split(allowedTokensEnv, ",")

	awsRegion := os.Getenv("AWS_REGION")
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	awsAccessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsDisableSSL := os.Getenv("AWS_DISABLE_SSL") == "true"
	awsS3Bucket := readEnvWithDefault("AWS_S3_BUCKET", "turborepo-cache")
	awsS3ForcePathStyle := os.Getenv("AWS_S3_FORCE_PATH_STYLE") == "true"

	return &Config{
		AWSRegion:           awsRegion,
		AWSEndpoint:         awsEndpoint,
		AWSAccessKeyID:      awsAccessKeyId,
		AWSSecretAccessKey:  awsSecretAccessKey,
		AWSDisableSSL:       awsDisableSSL,
		AWSS3Bucket:         awsS3Bucket,
		AWSS3ForcePathStyle: awsS3ForcePathStyle,
		TurboAllowedTokens:  allowedTokens,
	}
}

func readEnvWithDefault(variable, defaultValue string) string {
	value, found := os.LookupEnv(variable)
	if !found {
		return defaultValue
	}
	return value
}
