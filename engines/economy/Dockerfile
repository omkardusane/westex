# STEP 1: Build Stage (using a development image to compile)
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies (this step is cached if deps don't change)
RUN go mod download

# Copy all source code
COPY . .

# Build the application for a Linux target
# -o app sets the output binary name to 'app'
# -ldflags -s -w reduces the binary size
# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o app .

# ----------------------------------------------------------------------

# STEP 2: Final/Release Stage (using a minimal image to deploy)
# FROM alpine:latest

# Set the working directory
# WORKDIR /root/

# Copy the compiled binary from the builder stage
# COPY --from=builder /app/app .

# Expose the port your server listens on (e.g., 8080)
# EXPOSE 8080

# Command to run the executable
# CMD ["./app"]