# build stage
FROM golang:1.15-buster AS build-env

ADD . /src
ENV CGO_ENABLED=0
WORKDIR /src

# RUN ssh-keyscan github.com >> known_hosts

ARG TARGETOS
ARG TARGETARCH
RUN go build -o swarmoperator

# final stage
FROM alpine
RUN apk add --no-cache docker-cli
WORKDIR /app
COPY known_hosts /etc/ssh/ssh_known_hosts
COPY --from=build-env /src/swarmoperator /app/
ENTRYPOINT ["/app/swarmoperator"]
