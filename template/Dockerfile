ARG PROJECT_NAME=exporter

FROM golang:1.26 AS build

WORKDIR /src

ARG PROJECT_NAME
ARG LDFLAGS

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build \
  -trimpath \
  -ldflags "${LDFLAGS}" \
  -o "/out/${PROJECT_NAME}" \
  ./cmd

FROM debian:bookworm-slim

ARG PROJECT_NAME

COPY --from=build /out/${PROJECT_NAME} /usr/local/bin/exporter

ENTRYPOINT ["/usr/local/bin/exporter"]
