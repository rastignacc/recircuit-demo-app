FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod ./
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /app /app
COPY migrations /migrations
EXPOSE 8080
ENTRYPOINT ["/app"]
