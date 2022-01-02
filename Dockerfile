# Go builder
FROM golang:1.17-buster AS build-env

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY internal ./internal
COPY cmd ./cmd
COPY main.go .

# RUN ssh-keyscan github.com >> known_hosts

ENV CGO_ENABLED=0
RUN go build -o swarmops

# final image
FROM alpine

RUN apk add --no-cache docker-cli

WORKDIR /app

COPY known_hosts /etc/ssh/ssh_known_hosts
COPY --from=build-env /src/swarmops /app/

ENTRYPOINT ["/app/swarmops"]
CMD ["serve"]
