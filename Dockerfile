# syntax=docker/dockerfile:1

FROM golang:1.25-bookworm AS build

ARG GOPROXY=https://proxy.golang.org,direct
ARG ANIFEED_VERSION=0.02
ENV GOPROXY=${GOPROXY} \
    CGO_ENABLED=1 \
    GOOS=linux

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags "-s -w -X github.com/zhifengle/anifeed/version.Version=${ANIFEED_VERSION}" \
    -o /out/anifeed .

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /data
COPY --from=build /out/anifeed /usr/local/bin/anifeed

EXPOSE 8115
VOLUME ["/data"]

ENTRYPOINT ["anifeed"]
CMD ["server"]
