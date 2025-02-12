# Building the binary of the App
FROM golang:1.23 AS build

# `boilerplate` should be replaced with your project name
WORKDIR /go/src/taskexecutor

# Copy all the Code and stuff to compile everything
COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Builds the application as a staticly linked one, to allow it to run on alpine
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOMAXPROCS=1
RUN cd cmd && go build -o taskexec .


# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest AS release

WORKDIR /app

# Create the `public` dir and copy all the assets into it

# `boilerplate` should be replaced here as well
COPY --from=build /go/src/taskexecutor .

# Add packages
RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod +x /app/cmd/taskexec

ENTRYPOINT ["/app/cmd/taskexec"]