# syntax=docker/dockerfile:1

FROM golang:1.25-bookworm AS build

ARG GOPROXY=https://proxy.golang.org,direct
ARG ANIFEED_VERSION=0.02
ENV GOPROXY=${GOPROXY} \
    CGO_ENABLED=1 \
    GOOS=linux

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates musl-tools tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CC=musl-gcc go build -trimpath \
    -ldflags "-s -w -extldflags -static -X github.com/huohen92/anifeed/version.Version=${ANIFEED_VERSION}" \
    -o /out/anifeed .

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /out/anifeed /usr/local/bin/anifeed

ENV TZ=Asia/Shanghai
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

WORKDIR /data
VOLUME ["/data"]
EXPOSE 8115

ENTRYPOINT ["/usr/local/bin/anifeed"]
CMD ["server"]
