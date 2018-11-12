FROM golang:1.11 as builder

ADD . /go/src/github.com/fntlnz/caturday

WORKDIR /go/src/github.com/fntlnz/caturday

RUN GO111MODULE=on go build -a -tags netgo -ldflags '-w' .

FROM scratch

COPY --from=builder /go/src/github.com/fntlnz/caturday/caturday /caturday

USER 10000
EXPOSE 8080
ENTRYPOINT ["/caturday"]
