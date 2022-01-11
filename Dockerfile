# syntax=docker/dockerfile:1

FROM golang:1.17-buster AS build

WORKDIR /app
COPY . .

RUN go mod download && go build -o /turborepo-s3-remote-cache

FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /turborepo-s3-remote-cache /turborepo-s3-remote-cache

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/turborepo-s3-remote-cache"]
