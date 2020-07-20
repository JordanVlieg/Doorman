FROM golang:alpine

# Opt into Go modules available in Golang version 1.11
ENV GO111MODULE=on

# Move to working directory /build
WORKDIR /build

# Copy Go modules files so dependencies aren't redownloaded each time the image is built
COPY go.mod .
COPY go.sum .

# Download Go dependencies
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Command to run when starting the container
CMD ["./main"]

# Export necessary port
EXPOSE 80