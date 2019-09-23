FROM golang:alpine AS build-env


ADD  ./main.go /go/src
ADD ./vendor/queue /go/src/queue

RUN apk add --no-cache git
RUN go get -d -v github.com/google/uuid
RUN go get -d -v github.com/sirupsen/logrus


RUN cd /go/src && CGO_ENABLED=0 go build -o locationserver

FROM alpine
EXPOSE 8081
WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=build-env /go/src/locationserver /app/
ENTRYPOINT  ./locationserver

