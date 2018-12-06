# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
COPY . $GOPATH/src/github.com/alldroll/suggest
WORKDIR $GOPATH/src/github.com/alldroll/suggest

# Build binaries (TODO replace with makefile)
RUN make build

# Make scripts executable
RUN chmod +x suggest_run.sh && mv suggest_run.sh suggest_run
RUN chmod +x indexer_run.sh && mv indexer_run.sh indexer_run
