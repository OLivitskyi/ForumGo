# Use the official Golang image that supports CGO
FROM golang:alpine as builder

# Install dependencies for CGO and SQLite
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Remove `CGO_ENABLED=0` to allow the use of CGO
RUN go build -o forum ./cmd/forum/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/forum .

# Install libraries needed for SQLite
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/configs/config.json ./configs/
COPY --from=builder /app/internal/db/forum.db ./internal/db/
COPY --from=builder /app/web/ ./web/

EXPOSE 8080

CMD ["./forum"]