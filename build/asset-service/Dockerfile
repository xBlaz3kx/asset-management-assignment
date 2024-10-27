FROM golang:1.23 AS build
WORKDIR /app/src

ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/root/.cache/go-build
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

# Leverage build cache by copying go.mod and go.sum first
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
RUN go mod verify

# Copy the rest of the source code
COPY . .

RUN --mount=type=cache,target="/root/.cache/go-build" go build -o ../asset ./cmd/asset-service

# Debug stage using Delve
FROM build AS debug

COPY --from=build /app/src /asset/src

WORKDIR /asset/src

# Install Delve and build the application
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go build -gcflags="all=-N -l" -o ./asset ./cmd/asset-service
CMD ["/go/bin/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./asset"]

FROM alpine AS app

# Add curl for health checks
RUN apk --update --no-cache add curl

SHELL ["/bin/sh", "-l", "-c"]

# Create a non-root user
ARG user=app
ARG group=docker
RUN addgroup $group
RUN adduser \
    --disabled-password \
    --gecos "" \
    --ingroup $group \
    $user
USER $user
WORKDIR /usr/local/asset-service

# Copy the binary and set executable permissions
COPY --chown=$user:$group --from=build /app/asset /usr/local/bin/asset
RUN chmod +x /usr/local/bin/asset

# Add a health check
HEALTHCHECK --start-period=10s --interval=5s --timeout=10s --retries=3 \
    CMD curl --fail http://localhost:80/healthz || exit 1

ENTRYPOINT ["/usr/local/bin/asset"]