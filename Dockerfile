# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git (required for go build with git info)
RUN apk add --no-cache git

# Copy the entire project
COPY . .

# Build the binary with ldflags
RUN go build -ldflags="-s -w -X github.com/prometheus/common/version.Version=$(git describe --abbrev=0 --tags) \
    -X github.com/prometheus/common/version.Revision=$(git rev-parse --short HEAD) \
    -X github.com/prometheus/common/version.Branch=$(git rev-parse --abbrev-ref HEAD) \
    -X github.com/prometheus/common/version.BuildUser=$(git config user.name | tr -d ' ') \
    -X github.com/prometheus/common/version.BuildDate=$(TZ=UTC date +%FT%T%z)" \
    -o spectrum-virtualize-exporter

# Final stage
FROM alpine:3.23.4

# Install gcompat (required for Go binaries on Alpine)
RUN apk add --no-cache gcompat

# Copy the binary and config from the builder
COPY --from=builder /app/spectrum-virtualize-exporter /opt/spectrumVirtualize/spectrum-virtualize-exporter
COPY spectrumVirtualize.yml /opt/spectrumVirtualize/spectrumVirtualize.yml

# Port of Prometheus exporter endpoint
EXPOSE 9119

ENTRYPOINT ["/opt/spectrumVirtualize/spectrum-virtualize-exporter"]
CMD ["--config.file=/opt/spectrumVirtualize/spectrumVirtualize.yml"]
