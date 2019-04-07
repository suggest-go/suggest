# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
COPY . $GOPATH/src/github.com/alldroll/suggest
WORKDIR $GOPATH/src/github.com/alldroll/suggest

# Build binaries (TODO replace with makefile)
RUN make build
