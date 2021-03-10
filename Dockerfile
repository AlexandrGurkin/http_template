FROM quay.io/goswagger/swagger:v0.25.0 AS swagger

ADD . /app
WORKDIR /app
RUN swagger generate server --exclude-main -f "./api/swagger.yaml"

FROM golang:1.16 AS builder

COPY --from=swagger /app /app
WORKDIR /app
ENV CGO_ENABLED=0
RUN make build

FROM alpine
# ENV GODEBUG netdns=cgo
WORKDIR /app
ENV CGO_ENABLED=0
COPY --from=builder /app/test_http /usr/local/bin/test_http
CMD ["/usr/local/bin/test_http", "run"]