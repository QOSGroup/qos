# Base build image
FROM golang:1.11.5-alpine3.9 as builder

# Install some dependencies needed to build the project
RUN apk add make gcc git libc-dev bash linux-headers eudev-dev

# Set up GOPATH & PATH
ENV GOPATH       /root/go
ENV BASE_PATH    $GOPATH/src/github.com/
ENV REPO_PATH    $BASE_PATH/qos
ENV PATH         $GOPATH/bin:$PATH

# Force the go compiler to use modules
ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR $REPO_PATH

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

#This is the ‘magic’ step that will download all the dependencies that are specified in
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction this line)
RUN go mod download

# Add source files
COPY . $REPO_PATH/

# build qosd & qoscli
RUN go build -o qosd $REPO_PATH/cmd/qosd/main.go
RUN go build -o qoscli $REPO_PATH/cmd/qoscli/main.go

# new stage
FROM alpine:3.9

# p2p port
EXPOSE 26656
# rpc port
EXPOSE 26657

COPY --from=builder /root/go/src/github.com/qos/qosd /usr/local/bin/
COPY --from=builder /root/go/src/github.com/qos/qoscli /usr/local/bin/
