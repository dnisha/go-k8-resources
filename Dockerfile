# Stage 1: Build the Go application
FROM golang:1.20-alpine AS build

WORKDIR /app

# Copy the application source code
COPY main.go go.mod .

# Download dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o myapp

# Stage 2: Create the production image
FROM golang:1.20-alpine AS production

WORKDIR /app

# Copy the built application from the previous stage
COPY --from=build /app/myapp .

# Expose the desired port (replace 8080 with your application's port if different)
EXPOSE 8080

# Start the Go application
CMD ["./myapp"]
