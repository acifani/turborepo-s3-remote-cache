# turborepo-s3-remote-cache

HTTP server to store [Turborepo](https://turborepo.org/) remote cache to an AWS S3 bucket.

> :warning: Do not use in production :warning:

## Configuration

| Environment variable       | Description                                                                                                              |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| `TURBOREPO_ALLOWED_TOKENS` | **Required**. Allowed tokens for Turborepo authentication. Comma separated list                                          |
| `AWS_REGION`               | **Required**. Region of the S3 Bucket                                                                                    |
| `AWS_ACCESS_KEY_ID`        | **Required**. Access key ID                                                                                              |
| `AWS_SECRET_ACCESS_KEY`    | **Required**. Secret access key                                                                                          |
| `AWS_ENDPOINT`             | Leave empty for default AWS endpoint. Customize for S3 compatible storage (e.g. [min.io](https://min.io/))               |
| `AWS_DISABLE_SSL`          | Disable SSL in the AWS SDK                                                                                               |
| `AWS_S3_BUCKET`            | Default: `turborepo-cache`                                                                                               |
| `AWS_S3_FORCE_PATH_STYLE`  | Use legacy path for S3 objects                                                                                           |
| `GIN_TRUSTED_PROXIES`      | See [Gin docs](https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetTrustedProxies) for more info. Comma separated list |

For more info take a look at the
[AWS guide on configuring the SDK](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/)

## Running

### With Docker

```sh
docker run ghcr.io/acifani/turborepo-s3-remote-cache:latest
```

Check out all available tags in the [Packages page](https://github.com/acifani/turborepo-s3-remote-cache/pkgs/container/turborepo-s3-remote-cache).

### From pre-built binaries

Download a binary from the [Release page](https://github.com/acifani/turborepo-s3-remote-cache/releases), unpackage it and run it.

Binaries are available for Windows, MacOS, Linux in amd64 and arm64

### From source

```sh
git clone https://github.com/acifani/turborepo-s3-remote-cache.git
cd turborepo-s3-remote-cache
go build ./...
./turborepo-s3-remote-cache
```

## Usage

You will need to configure the API endpoint and the auth token.

E.g.

```sh
export $TURBOREPO_TOKEN="some_t0k3n"
turbo run build --api="http://localhost:8080" --team="my-team" --token=$TURBOREPO_TOKEN
```

You will see this message if the remote cache has been correctly enabled.

> `• Remote computation caching enabled`

If you provide a team name (recommended), the cache will be stored in
a dedicated directory within the bucket.

### Config file

You can also create a `.turbo/config.json` file and configure the API server and team id/slug there.

E.g.

```json
{
  "apiUrl": "http://localhost:8080",
  "teamSlug": "my-team"
}
```

and then run

```sh
turbo run build --token=$TURBOREPO_TOKEN
```
