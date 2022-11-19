# builder
FROM golang:1.18.8-alpine3.16 as builder

WORKDIR /go/service/src
COPY . /go/service/src

RUN GOOS=linux go build -o /go/ecom /go/service/src/main.go && \
    chmod +x /go/ecom

# runner
FROM alpine:3.16 as runner

ARG SERVICE
ENV ECOM_SERVICE=${SERVICE}

WORKDIR /ecom

COPY --from=builder /go/ecom ./ecom_bin

CMD [ "./ecom_bin" ]
