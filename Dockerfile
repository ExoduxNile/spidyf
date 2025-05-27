# Use official Go image as the base image
FROM golang:1.21

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy source code, static files, and configuration
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY static/ ./static/
COPY config/ ./config/

# Build the application
RUN go build -o bin/spidy -v cmd/spidy/main.go

# Expose port 8080 for the web server
EXPOSE 8080

# Command to run the application
CMD ["./bin/spidy"]
