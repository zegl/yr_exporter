# syntax=docker/dockerfile:1
FROM golang:1.20-bullseye as builder

WORKDIR /app
ADD go.mod /app/go.mod
ADD go.sum /app/go.sum
RUN go mod download
ADD ./ /app
RUN --mount=type=cache,target=/root/.cache/go-build,id=go-build \
    --mount=type=cache,target=/root/.cache/go-mod,id=go-cache \
    export GOMODCACHE=/root/.cache/go-mod && \
    CGO_ENABLED=0 go build -tags netgo -o /bin/yr_exporter

FROM debian:bullseye as runtime
RUN apt-get update && apt-get install ca-certificates -y
COPY --from=builder /bin/yr_exporter /bin/yr_exporter
EXPOSE      9367
ENTRYPOINT  [ "/bin/yr_exporter" ]

LABEL org.opencontainers.image.source https://github.com/zegl/yr_exporter