# A hello world example with Go
FROM golang:1.8-onbuild
MAINTAINER rafa87sch@gmail.com

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]