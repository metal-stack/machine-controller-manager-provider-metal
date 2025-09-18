FROM golang:1.25 AS builder
WORKDIR /work
COPY . .
RUN make build

FROM alpine:3.22
RUN apk add --update bash curl tzdata
WORKDIR /
COPY --from=builder /work/bin/machine-controller /machine-controller
ENTRYPOINT ["/machine-controller"]
