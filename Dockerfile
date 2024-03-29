FROM golang:1.18-alpine AS build
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY . /app
WORKDIR /app
COPY . .

RUN apk add --no-cache git ca-certificates

RUN go get -d
RUN go build -o main .
FROM alpine:3.16 AS runtime

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist
COPY --from=build /app/main ./

# Export necessary port
EXPOSE 8080

ENTRYPOINT ["/dist/main"]
