# Use official Golang image
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Use minimal alpine image for final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tesseract-ocr

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

# Create necessary directories
RUN mkdir -p uploads tmp/ocr logs

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
