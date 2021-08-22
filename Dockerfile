## Get dependencies
FROM golang:alpine AS dependencies
WORKDIR /src
ENV CGO_ENABLED=0

# Copy Modd files
COPY go.* /
RUN go mod download
COPY . .

# Create build
FROM dependencies AS build
WORKDIR /src
# Build go application
RUN go build -v -o ./app ./cmd/server

## Create release build, only take go application into the release
FROM scratch AS bin
COPY --from=build src/app ./app
CMD ["./app"]

