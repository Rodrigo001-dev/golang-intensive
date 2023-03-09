FROM golang:1.20

WORKDIR /go/goapp

RUN apt-get && apt-get install -y librdkafka-dev

CMD ["tail", "-f", "/dev/null"]