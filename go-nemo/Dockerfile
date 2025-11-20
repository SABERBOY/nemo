FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o nemo-server ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/nemo-server .
# Assuming config is in internal/config or root config folder. 
# The code uses viper to look in "config" path.
COPY --from=builder /app/config ./config

EXPOSE 9981

CMD ["./nemo-server"]
