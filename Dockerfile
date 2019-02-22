# Simple usage with a mounted data directory:
# > docker build -t qosd .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.qosd:/root/.qosd -v ~/.qoscli:/root/.qoscli qosd qosd init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.qosd:/root/.qosd -v ~/.qoscli:/root/.qoscli qosd gstarsd start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES git bash gcc make libc-dev linux-headers eudev-dev
#make libc-dev bash gcc linux-headers eudev-dev

# Set working directory for the build
WORKDIR /go/src/github.com/QOSGroup/qos

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
  #    make tools && \
  #    make vendor-deps && \
  #    make build && \
  #    make install
  cd cmd/qosd && \
  export GOPROXY=http://192.168.1.177:8081 && \
  export GO111MODULE=on && \
  go build && \
  cd ../qoscli && \
  go build

# Final image
FROM alpine:edge

# Install ca-certificates
#RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/QOSGroup/qos/cmd/qosd/qosd /usr/bin/qosd
COPY --from=build-env /go/src/github.com/QOSGroup/qos/cmd/qoscli/qoscli /usr/bin/qoscli

# Run gaiad by default, omit entrypoint to ease using container with gaiacli
CMD ["qosd"]
