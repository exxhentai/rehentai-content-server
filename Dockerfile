FROM golang:latest

COPY ./ /go/src/rehentai-content-server
WORKDIR /go/src/rehentai-content-server
ENV GOPATH=/go/src/rehentai-content-server

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 1234

# Command to run the executable
CMD ["./main"]