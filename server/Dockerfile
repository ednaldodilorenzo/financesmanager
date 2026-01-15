# Building the binary of the App
FROM golang:1.25-alpine3.22 AS build

# `boilerplate` should be replaced with your project name
WORKDIR /go/src/financeweb


COPY go.* ./

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Copy all the Code and stuff to compile everything
COPY . .

# Builds the application as a staticly linked one, to allow it to run on alpine
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOMAXPROCS=1
# -ldflags="-s -w": Strips debugging info, making the binary smaller
RUN cd cmd && go build -ldflags="-s -w" -o app .


# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest AS release

WORKDIR /app

# Create the `public` dir and copy all the assets into it

# `boilerplate` should be replaced here as well
COPY --from=build /go/src/financeweb .

# Add packages
RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod +x /app/cmd/app

# Exposes port 5000 because our program listens on that port
EXPOSE 5000

ENTRYPOINT ["/app/cmd/app"]
