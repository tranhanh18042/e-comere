FROM golang:1.18.8-alpine3.16 as build

RUN mkdir -p /go/services/item
COPY . /go/services/item
RUN GOOS="linux" go build -o /go/services/item/svc_item /go/services/item

FROM alpine:3.16
RUN mkdir
COPY --from=build
