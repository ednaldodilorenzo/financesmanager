# Building the binary of the App
FROM golang:1.24 AS development

# `boilerplate` should be replaced with your project name
WORKDIR /app

# Install Air for live-reloading
RUN go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

# Ensure that Air is in PATH
ENV PATH="/go/bin:$PATH"

# Copy all the Code and stuff to compile everything
# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

EXPOSE 5000 2345

# Use Air as the entrypoint
CMD ["air"]

# Building the binary of the App
FROM golang:1.24 AS build

# `boilerplate` should be replaced with your project name
WORKDIR /go/src/financeweb

# Copy all the Code and stuff to compile everything
COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Builds the application as a staticly linked one, to allow it to run on alpine
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOMAXPROCS=1
# -ldflags="-s -w": Strips debugging info, making the binary smaller
RUN cd cmd && go build -o app -ldflags="-s -w" .


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