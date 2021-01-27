## Add Air for development live-reload
FROM golang:alpine  AS dev
WORKDIR /src
RUN apk add --no-cache curl
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

## Add the wait script to the image
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

CMD /wait && air

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
RUN echo NOPE!
CMD ["./app"]

