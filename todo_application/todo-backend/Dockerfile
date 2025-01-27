FROM golang:1.23.4 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /todo-backend

FROM alpine:latest
COPY --from=builder /todo-backend /todo-backend

CMD ["/todo-backend"]
