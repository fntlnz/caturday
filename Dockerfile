FROM golang:1.8 as builder

ADD . /go/src/github.com/fntlnz/caturday

WORKDIR /go/src/github.com/fntlnz/caturday

RUN go build -a -tags netgo -ldflags '-w' .

FROM scratch

COPY --from=builder /go/src/github.com/fntlnz/caturday/caturday /caturday

EXPOSE 8080
ENTRYPOINT ["/caturday"]
