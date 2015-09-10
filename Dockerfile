FROM golang:latest
WORKDIR /go/src/github.com/EtienneR/go_upload
ADD . /go/src/github.com/EtienneR/go_upload
RUN go install github.com/EtienneR/go_upload
ENTRYPOINT /go/bin/go_upload
EXPOSE 3000