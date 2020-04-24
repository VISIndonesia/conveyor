FROM golang:1.13
VOLUME /mnt

ENV DEBICRED_CONFIG=/mnt/config.yaml

ADD . /app
WORKDIR /app

RUN go build ./cmd/subscriber

CMD ./subscriber
