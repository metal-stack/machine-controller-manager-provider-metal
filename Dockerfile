#############      builder                                  #############
FROM golang:1.14 AS builder

WORKDIR /work
COPY . .

RUN go mod download
RUN make build

#############      base                                     #############
FROM alpine:3.12 as base

RUN apk add --update bash curl tzdata
WORKDIR /

#############      machine-controller               #############
FROM base AS machine-controller

COPY --from=builder /work/bin/machine-controller /machine-controller
ENTRYPOINT ["/machine-controller"]
