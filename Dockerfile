# build stage
FROM golang:1.15-buster AS build-env

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# RUN ssh-keyscan github.com >> known_hosts

ENV CGO_ENABLED=0
RUN go build -o swarmoperator

FROM alpine
RUN apk add --no-cache docker-cli
WORKDIR /app
COPY known_hosts /etc/ssh/ssh_known_hosts
COPY --from=build-env /src/swarmoperator /app/
ENTRYPOINT ["/app/swarmoperator"]
