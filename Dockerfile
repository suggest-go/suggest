FROM golang:1.14.4-alpine as builder

RUN set -eux; \
  apk add --no-cache git \
  make

# Copy the local package files to the container's workspace.
COPY . /data
WORKDIR /data

# Build binaries
RUN CGO_ENABLED=0 make build-bin BUILD_FLAGS='-ldflags="-w -s"'

FROM scratch

COPY --from=builder /data/build /data/build

CMD ["/data/build/suggest"]
