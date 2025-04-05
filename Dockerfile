# Build
FROM golang:latest as build

WORKDIR /go/src/app

COPY . .
RUN CGO_ENABLED=0 make build

# Run
FROM alpine:latest
WORKDIR /app

RUN apk add tzdata
COPY --from=build /go/src/app/bin/hkp /app/hkp
COPY views ./views

ENTRYPOINT ["/app/hkp", "serve", "-v", "-d", "/db", "-c", "/config/homekit.toml"]
