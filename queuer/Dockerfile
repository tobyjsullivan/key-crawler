FROM golang:1

ADD . /go/src/github.com/tobyjsullivan/key-crawler/
RUN go install github.com/tobyjsullivan/key-crawler/queuer

CMD /go/bin/queuer
