FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o nemo-server ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/nemo-server .
COPY config/application.yml .

EXPOSE 9981

CMD ["./nemo-server"]
