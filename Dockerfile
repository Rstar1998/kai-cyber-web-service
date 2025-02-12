# Useing a multi-stage build to reduce image size
FROM golang:1.23-alpine AS builder 

WORKDIR /app

# Installing build tools for CGO (musl-dev for Alpine)
ENV CGO_ENABLED=1

RUN apk add --no-cache musl-dev gcc

# Copy go mod and sum files for caching
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Create a minimal runtime image
FROM alpine:latest AS runner

WORKDIR /app

# Copy necessary files from the builder stage
COPY --from=builder /app/main /app/main

# Expose the port your application listens on
EXPOSE 8080

# Run the application
CMD ["./main"]