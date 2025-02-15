# Use an official Go image as the base image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the current directory (i.e. the directory with the Dockerfile) into the container at /app
COPY . /app

# Install dependencies
RUN go mod download

# Build the Go application
RUN go build -o main .

# Expose port 8080 for the web application
EXPOSE 8080

# Run the Go web application
CMD ["./main"]
