# syntax=docker/dockerfile:1

FROM gcr.io/distroless/base-debian10

COPY turborepo-s3-remote-cache /usr/bin/turborepo-s3-remote-cache
EXPOSE 8080

ENTRYPOINT ["/usr/bin/turborepo-s3-remote-cache"]
