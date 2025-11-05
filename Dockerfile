FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /wiki-voyage-app

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /wiki-voyage-app /wiki-voyage-app
COPY templates ./templates

USER nonroot:nonroot

ENTRYPOINT ["/wiki-voyage-app"]