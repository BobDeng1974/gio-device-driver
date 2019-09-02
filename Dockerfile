FROM golang:alpine AS builder

WORKDIR /devicedriver

# Install git for fetching dependencies
RUN apk update && apk add --no-cache git

COPY go.mod .

RUN go mod download

COPY . .

# Build the binary.
RUN go build -o /go/bin/devicedriver cmd/devicedriver/main.go

## Build lighter image
FROM alpine:latest
LABEL Name=gio-device-driver-go Version=1.0.0

# Copy our static executable.
COPY --from=builder /go/bin/devicedriver /devicedriver

EXPOSE 8080

# Run the binary.
ENTRYPOINT /devicedriver