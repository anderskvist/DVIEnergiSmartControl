# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang v1.11 base image
FROM golang:1.26-alpine AS builder

# Add Maintainer Info
LABEL maintainer="Anders Kvist <anderskvist@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

# Install the package
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X github.com/anderskvist/GoHelpers/version.Version=${BUILD_TIME}-${GIT_COMMIT}" \
    -o /app/DVIEnergiSmartControl .

FROM alpine:3.19

# Add root certificates for outbound SSL/TLS API requests
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder stage
COPY --from=builder /app/DVIEnergiSmartControl /app/DVIEnergiSmartControl

# Run the executable
CMD ["/app/DVIEnergiSmartControl","/config.ini"]
