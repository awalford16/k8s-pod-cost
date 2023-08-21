# Use an official Golang runtime as a parent image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Build the Golang application
RUN go build -o myapp

# Expose a port (if your application listens on a specific port)
EXPOSE 8080

# Define the command to run when the container starts
CMD ["./myapp"]
