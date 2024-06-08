FROM golang:1.22.1
MAINTAINER Trinh Duc Chung <chungchuong6b@gmail.com>
WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download
ADD server.go /go/src/server.go
EXPOSE 1323
ENTRYPOINT [ "go", "run",  "server.go"]