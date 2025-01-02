# Start with a lightweight Go image with necessary tools for CGO
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk update && apk add --no-cache git gcc musl-dev

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Set GOPROXY to direct to bypass module proxy
ENV GOPROXY=direct
ENV GOPRIVATE=*

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with CGO enabled
RUN CGO_ENABLED=1 go build -o app .

# Use a lightweight runtime image for the final stage
FROM alpine:latest

# Install runtime dependencies (if needed, like ca-certificates)
RUN apk add --no-cache ca-certificates

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/app ./

# Expose the application port
EXPOSE 8080

# Set environment variables with defaults
ENV DB_FILE="app.db"

# Run the application
CMD ["./app"]
