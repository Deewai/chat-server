FROM golang:latest

LABEL maintainer="Innocent Abdullahi <deewai48@gmail.com>"

RUN apt-get update


#for debugging purposes
# RUN go get github.com/derekparker/delve/cmd/dlv

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/Deewai/chat-server

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
ENV GO111MODULE=on
RUN go mod download
# RUN go mod vendor

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .
RUN go install -v ./...

# This container exposes port 8000 to the outside world
EXPOSE 8080

# Run the executable
CMD ["chat-server"]

#for debugging purposes
# CMD [ "dlv", "debug", "Bridge/backend", "--listen=:40000", "--headless=true", "--api-version=2", "--log" ]