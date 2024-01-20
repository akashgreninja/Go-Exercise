# Use the official Golang image as the base
FROM golang:1.21.5

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 go build -o /app/app

# Expose port 4000
EXPOSE 4000

# Set the entrypoint for the container to run the binary
ENTRYPOINT ["/app/app"]
