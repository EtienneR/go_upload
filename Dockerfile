FROM golang:latest
WORKDIR /go/src/github.com/EtienneR/go_upload
ADD . /go/src/github.com/EtienneR/go_upload
RUN go install /go/bin/go_upload
ENTRYPOINT /go/bin/go_upload
EXPOSE 3000