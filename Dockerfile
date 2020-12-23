# Go builder
FROM golang:1.15-buster AS build-env

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY internal ./internal
COPY cmd ./cmd
COPY main.go .

# RUN ssh-keyscan github.com >> known_hosts

ENV CGO_ENABLED=0
RUN go build -o swarmops

# UI builder
FROM node:14 AS ui-build-env

WORKDIR /src

COPY ./ui/package.json ./ui/yarn.lock ./
RUN yarn install --frozen-lockfile

COPY ./ui .
RUN yarn build

# final image
FROM alpine

RUN apk add --no-cache docker-cli

WORKDIR /app

COPY known_hosts /etc/ssh/ssh_known_hosts
COPY --from=build-env /src/swarmops /app/
COPY --from=ui-build-env /src/dist /app/ui

ENTRYPOINT ["/app/swarmops"]
CMD ["serve"]
