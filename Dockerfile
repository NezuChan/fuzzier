FROM golang:1.20-alpine as build-stage

WORKDIR /tmp/build

COPY . .

# Build the project
RUN go build .

FROM alpine:3

LABEL name "NezuChan fuzzier"
LABEL maintainer "KagChi"

WORKDIR /app

# Install needed deps
RUN apk add --no-cache tini

COPY --from=build-stage /tmp/build/fuzzier main

ENTRYPOINT ["tini", "--"]
CMD ["/app/main"]