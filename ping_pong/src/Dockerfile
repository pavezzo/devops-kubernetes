FROM golang:1.23.4 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /ping-pong

FROM alpine:latest
COPY --from=builder /ping-pong /ping-pong

CMD ["/ping-pong"]
