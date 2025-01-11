FROM golang:1.23.4 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /todo-application

FROM alpine:latest
COPY --from=builder /todo-application /app/todo-application
COPY public /app/public

WORKDIR /app
CMD ["/app/todo-application"]
