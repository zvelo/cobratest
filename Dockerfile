FROM golang:alpine as build
RUN apk add --no-cache git
RUN go get github.com/magefile/mage
COPY . /go/src/zvelo.io/cobratest
WORKDIR /go/src/zvelo.io/cobratest
RUN mage -v build

FROM alpine:latest
MAINTAINER Joshua Rubin <jrubin@zvelo.com>
ENTRYPOINT ["cobratest"]
COPY --from=build /go/src/zvelo.io/cobratest/cobratest /usr/local/bin/cobratest
