FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY .. .

RUN CGO_ENABLED=0 GOOS=linux go build -o nails_game ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/nails_game .
COPY --from=builder /app/.env .
COPY --from=builder /app/seed_players.json .

EXPOSE 8080

CMD ["/bin/sh", "-c", "sleep 5 && ./nails_game"]