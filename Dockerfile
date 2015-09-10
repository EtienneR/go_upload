FROM golang:latest
WORKDIR /go/src/github.com/etienner/goupload
ADD . /go/src/github.com/etienner/goupload
RUN go install /go/bin/goupload
ENTRYPOINT /go/bin/goupload
EXPOSE 3000