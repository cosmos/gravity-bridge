# Simple usage with a mounted data directory:
# > docker build -t simapp .
#
# Server:
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.simapp:/root/.simapp simapp gravity init test-chain
# TODO: need to set validator in genesis so start runs
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.simapp:/root/.simapp simapp gravity start
#
# Client: (Note the simapp binary always looks at ~/.simapp we can bind to different local storage)
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.simappcli:/root/.simapp simapp gravity keys add foo
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.simappcli:/root/.simapp simapp gravity keys list
# TODO: demo connecting rest-server (or is this in server now?)
FROM golang:alpine AS build-env

# Install minimum necessary dependencies,
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3
RUN apk add --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/althea-net/cosmos-gravity-bridge/module

# Get dependancies - will also be cached if we won't change mod/sum
COPY go.mod .
COPY go.sum .
RUN go mod download

# Add source files
COPY . .

# install simapp, remove packages
RUN go build -o build/gravity ./cmd/gravity/main.go

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add bash
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/althea-net/cosmos-gravity-bridge/module/build/gravity /usr/bin/gravity

EXPOSE 26656 26657 1317 9090

# Run gravity by default
CMD ["gravity", "--home", "home", "start", "--pruning=nothing"]