# Use the official Golang base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Build the Go application
RUN go build -o myapp

# Expose port 8080 (if your app listens on this port)
EXPOSE 8080

# Define the entry point for your application
CMD ["./myapp"]
