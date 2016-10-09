FROM scratch

ADD dist/caturday /caturday

ENTRYPOINT ["/caturday"]
