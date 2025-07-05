# Build
FROM golang:alpine as build

WORKDIR /go/src/app

# Install build dependencies for CGO
RUN apk add --no-cache make gcc musl-dev

COPY . .
RUN CGO_ENABLED=1 make build

# Run
FROM alpine:latest
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache tzdata
COPY --from=build /go/src/app/bin/hkp /app/hkp
COPY web ./web

ENTRYPOINT ["/app/hkp", "serve", "-v", "-d", "/db", "-c", "/config"]
