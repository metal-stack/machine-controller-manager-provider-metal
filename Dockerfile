FROM golang:1.15 AS builder
WORKDIR /work
COPY . .
RUN make build

FROM alpine:3.12
RUN apk add --update bash curl tzdata
WORKDIR /
COPY --from=builder /work/bin/machine-controller /machine-controller
ENTRYPOINT ["/machine-controller"]
