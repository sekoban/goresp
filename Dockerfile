# Use official Golang image as the base image
FROM golang:1.20-alpine

# Set working directory
WORKDIR /app

# Copy Go modules and build files
#COPY go.mod ./
#RUN go mod download

# Copy the source code
COPY main.go .

# Build the application
RUN go mod init goresp && go mod tidy && go build -o /goresp

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["/goresp"]
