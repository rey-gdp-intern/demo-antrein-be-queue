# Build Stage
FROM golang:1.22-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/binary application/*.go

# Final Stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/binary .
COPY --from=builder /app/files/secrets/secrets.config.json ./files/secrets/secrets.config.json

EXPOSE 8080
EXPOSE 9090

ENTRYPOINT ["./binary"]

# trigger pipeline
# trigger pipeline