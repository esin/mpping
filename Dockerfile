FROM golang:alpine as builder

WORKDIR /mpping

ADD mpping.go /mpping

RUN apk update && apk add git && \
    go get -u github.com/asaskevich/govalidator && \
    go get -u github.com/sethgrid/curse && \
    go build

FROM alpine:3.8

COPY --from=builder /mpping/mpping /app/mpping

ENTRYPOINT [ "/app/mpping" ]