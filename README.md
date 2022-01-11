# turborepo-s3-remote-cache

HTTP server to store [Turborepo](https://turborepo.org/) remote cache to an AWS S3 bucket.

> :warning: Do not use in production :warning:

## Configuration

| Environment variable     | Description                                                                                                |
| ------------------------ | ---------------------------------------------------------------------------------------------------------- |
| TURBOREPO_ALLOWED_TOKENS | Allowed tokens for Turborepo authentication. Comma separated list                                          |
| AWS_REGION               | Region of the S3 Bucket                                                                                    |
| AWS_ENDPOINT             | Leave empty for default AWS endpoint. Customize for S3 compatible storage (e.g. [min.io](https://min.io/)) |
| AWS_ACCESS_KEY_ID        | Access key ID                                                                                              |
| AWS_SECRET_ACCESS_KEY    | Secret access key                                                                                          |
| AWS_DISABLE_SSL          | Disable SSL in the AWS SDK                                                                                 |
| AWS_S3_BUCKET            | Defaults to `turborepo-cache`                                                                              |
| AWS_S3_FORCE_PATH_STYLE  | Use legacy path for S3 objects                                                                             |

For more info take a look at the
[AWS guide on configuring the SDK](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/)

## Running

Either build from source or use the docker image

```
ghcr.io/acifani/turborepo-s3-remote-cache:latest
```
