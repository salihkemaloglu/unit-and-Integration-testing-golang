# Stage 1 (to create a "build" image, ~850MB)
FROM golang:1.10.1 AS builder
RUN go version

COPY . /go/src/github.com/miguno/golang-docker-build-tutorial/
WORKDIR /go/src/github.com/miguno/golang-docker-build-tutorial/
RUN  go get gopkg.in/mgo.v2/bson
RUN  go get github.com/salihkemaloglu/UnitAndIntegrationTesting-Golang/operations
# If you need SSL certificates for HTTPS, replace `FROM SCRATCH` with:

EXPOSE 8081
ENTRYPOINT ["go","test","-v"]