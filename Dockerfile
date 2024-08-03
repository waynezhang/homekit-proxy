FROM golang:latest

WORKDIR /app

COPY . .

RUN make build

ENTRYPOINT ["/app/bin/hkp", "serve", "-v", "-d", "/db", "-c", "/config/homekit.toml"]
