FROM golang:alpine AS build-env

ADD  ./main.go /go/src
ADD ./vendor/queue /go/src/queue

RUN apk add --no-cache git
RUN go get -d -v github.com/google/uuid
RUN go get -d -v github.com/sirupsen/logrus
RUN go get -d -v github.com/emirpasic/gods/lists/doublylinkedlist


RUN cd /go/src && CGO_ENABLED=0 go build -o locationserver

FROM alpine
ADD static /app/static
WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=build-env /go/src/locationserver /app/

RUN addgroup --gid 1000 go && adduser -D -G go -u 100 go
RUN chown go ./locationserver
USER go

EXPOSE 8081
ENTRYPOINT  ./locationserver

