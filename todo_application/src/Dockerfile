FROM golang:1.23.4 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /todo-application

FROM alpine:latest
COPY --from=builder /todo-application /todo-application

CMD ["/todo-application"]
