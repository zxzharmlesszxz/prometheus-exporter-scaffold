ARG PROJECT_NAME=exporter

FROM golang:1.26 AS build

WORKDIR /src

ARG VERSION=dev
ARG BRANCH=dev
ARG REVISION=dev
ARG BUILD_USER=docker
ARG BUILD_DATE=unknown
ARG PROJECT_NAME

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build \
  -trimpath \
  -ldflags "-s -w \
  -X github.com/prometheus/common/version.Version=${VERSION} \
  -X github.com/prometheus/common/version.Branch=${BRANCH} \
  -X github.com/prometheus/common/version.Revision=${REVISION} \
  -X github.com/prometheus/common/version.BuildUser=${BUILD_USER} \
  -X github.com/prometheus/common/version.BuildDate=${BUILD_DATE}" \
  -o "/out/${PROJECT_NAME}" \
  ./cmd

FROM debian:bookworm-slim

ARG PROJECT_NAME

COPY --from=build /out/${PROJECT_NAME} /usr/local/bin/exporter

ENTRYPOINT ["/usr/local/bin/exporter"]
