FROM golang:1.11


COPY . /go/src/event-filter-pool
WORKDIR /go/src/event-filter-pool

ENV GO111MODULE=on

RUN go build

EXPOSE 8080

CMD ./event-filter-pool