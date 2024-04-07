# Development environment
# -----------------------
# pinned because of https://github.com/moby/moby/issues/45935
FROM golang:1.20.5-alpine
WORKDIR /porter

RUN apk update && apk add --no-cache gcc musl-dev git

# for live reloading of go container
RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN chmod +x /porter/docker/bin/*

CMD air -c .air.toml
