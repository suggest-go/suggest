# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.12

# Copy the local package files to the container's workspace.
COPY . /data
WORKDIR /data

# Build binaries
RUN make build
