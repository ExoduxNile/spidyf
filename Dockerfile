# Use official Go image as the base image
FROM golang:1.21

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download and verify dependencies
RUN go mod tidy && go mod verify

# Explicitly fetch missing dependencies
RUN go get github.com/gorilla/websocket@v1.5.0
RUN go get github.com/PuerkitoBio/goquery@v1.8.1
RUN go get golang.org/x/net/publicsuffix@v0.17.0
RUN go get gopkg.in/yaml.v3@v3.0.1
# Note: twiny modules may require replace directives or authentication
RUN go get github.com/twiny/domaincheck || echo "Warning: github.com/twiny/domaincheck not found"
RUN go get github.com/twiny/flog || echo "Warning: github.com/twiny/flog not found"
RUN go get github.com/twiny/wbot || echo "Warning: github.com/twiny/wbot not found"
RUN go get github.com/twiny/carbon || echo "Warning: github.com/twiny/carbon not found"

# Copy source code, static files, and configuration
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY static/ ./static/
COPY config/ ./config/

# Debug: List copied files to verify structure
RUN ls -laR /app

# Debug: Check contents of key files
RUN cat /app/cmd/spidy/api/version || echo "version file missing or empty"
RUN cat /app/go.mod || echo "go.mod file missing or empty"

# Rename example.config.yaml to config.yaml to match GUI default
RUN if [ -f ./config/example.config.yaml ]; then mv ./config/example.config.yaml ./config/config.yaml; else echo "example.config.yaml not found"; fi

# Build the application with verbose output
RUN go build -v -o bin/spidy cmd/spidy/main.go || { echo "go build failed"; exit 1; }

# Expose port 8080 for the web server
EXPOSE 8080

# Command to run the application
CMD ["./bin/spidy"]
