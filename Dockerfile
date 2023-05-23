# Use an official Golang runtime as the base image
FROM golang:1.19-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Set the environment variable for the port
ENV PORT=8080

# Expose the port on which your application will listen
EXPOSE 8080

# Start the Go application
CMD ["./main"]
