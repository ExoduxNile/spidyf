# Build stage
FROM golang:1.21 as builder

# Install Fyne dependencies and build tools
RUN apt-get update && apt-get install -y \
    libgl1-mesa-dev \
    xorg-dev \
    libgles2-mesa-dev \
    gcc \
    g++ \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

# Download Go modules
RUN go mod download

# Explicitly install Fyne (required for some systems)
RUN go get fyne.io/fyne/v2@latest

# Build with CGO enabled (required for Fyne)
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o /spidy-gui ./cmd/spidy-gui

# Final stage
FROM ubuntu:22.04

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    libgl1 \
    libx11-6 \
    libegl1 \
    xvfb \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /spidy-gui /spidy-gui

# Set up virtual display and run
CMD ["/bin/bash", "-c", "Xvfb :99 -screen 0 1024x768x24 & export DISPLAY=:99 && /spidy-gui"]
