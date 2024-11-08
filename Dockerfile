# Build stage
# Update to a newer Go version that is compatible with your dependencies
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

# Install necessary packages if any are required for dependencies
RUN go mod download

# Build the Go application
RUN go build -o main main.go

# Run stage
FROM alpine:3.13
WORKDIR /app

COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
CMD [ "/app/main" ]
