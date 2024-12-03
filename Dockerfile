# Start from the official Go image to build the application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project
COPY ./src .

# Build the Go app with the correct settings for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Final stage: use a minimal base image
FROM alpine:latest

# Set the working directory inside the minimal image
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose port 8080 (or any other port the server is using)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
