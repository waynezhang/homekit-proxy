# Build
FROM golang:latest as build

WORKDIR /go/src/app

COPY . .
RUN CGO_ENABLED=0 make build

# Run
FROM alpine:latest

COPY --from=build /go/src/app/bin/hkp /hkp

ENTRYPOINT ["/hkp", "serve", "-v", "-d", "/db", "-c", "/config/homekit.toml"]
