# Use the official Go image
FROM golang:1.21 as builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /spidy-gui ./cmd/spidy-gui

# Final stage
FROM ubuntu:22.04

# Install Fyne dependencies
RUN apt-get update && apt-get install -y \
    libgl1-mesa-dev \
    xorg-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy the binary from builder
COPY --from=builder /spidy-gui /spidy-gui

# Set the entrypoint
ENTRYPOINT ["/spidy-gui"]
