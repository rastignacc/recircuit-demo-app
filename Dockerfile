# Stage 1: Build Vue frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install
COPY frontend/ .
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.22-alpine AS backend-builder
WORKDIR /build
COPY backend/go.mod ./
COPY backend/ .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app ./cmd/server

# Stage 3: Final minimal image
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=backend-builder /app /app
COPY backend/migrations /migrations
COPY --from=frontend-builder /frontend/dist /static
EXPOSE 8080
ENTRYPOINT ["/app"]
