FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o spamhaus
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/spamhaus .
COPY config.json .
EXPOSE 8080
CMD ["./spamhaus"]