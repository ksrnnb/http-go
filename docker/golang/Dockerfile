FROM golang:1.17-alpine

RUN apk update && \
    apk add \
    curl \
    gcc \
    git \
    musl-dev \
    && curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

WORKDIR /go/src

CMD ["air"]