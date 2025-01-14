FROM golang:1.23.4 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /log-output-writer

FROM alpine:latest
COPY --from=builder /log-output-writer /log-output-writer

CMD ["/log-output-writer"]
