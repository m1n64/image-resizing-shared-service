FROM golang:1.23-alpine AS builder

RUN apk add --no-cache build-base libwebp-dev

WORKDIR /app
COPY . .

ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GIN_MODE=release
RUN go build -ldflags "-X main.GinMode=release" -o imageresizer ./cmd/server

FROM alpine:latest

RUN apk add --no-cache libwebp

WORKDIR /app
COPY --from=builder /app/imageresizer .

EXPOSE 5689
EXPOSE 50066

ENTRYPOINT ["./imageresizer"]
