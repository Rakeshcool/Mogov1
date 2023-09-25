# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application
RUN go get -d -v ./...
RUN go install -v ./...

# Expose a port if your Go application listens on a specific port
EXPOSE 8080

# Command to run the application
CMD ["mogov1"]
