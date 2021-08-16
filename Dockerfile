FROM golang:1.14.4

COPY ./web /go/src/app
WORKDIR /go/src/app

RUN set -eux && \
  curl -fLo /go/bin/air https://git.io/linux_air && \
  chmod +x /go/bin/air

RUN go build -o /go/app/main

CMD /go/app/main $PORT

