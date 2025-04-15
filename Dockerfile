FROM golang:1.23

WORKDIR ${GOPATH}/avito-pvz/
COPY . ${GOPATH}/avito-pvz/

RUN go build -o /build ./cmd/server \
    && go clean -cache -modcache

EXPOSE 8080
EXPOSE 9000

CMD ["/build"]