FROM golang:1.23.4 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /log-output-reader

FROM alpine:latest
COPY --from=builder /log-output-reader /log-output-reader

CMD ["/log-output-reader"]
