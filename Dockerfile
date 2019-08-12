FROM golang:alpine as builder

WORKDIR /mpping

ADD mpping.go /mpping

RUN apk update && apk add --no-cache git && \
    go get -u github.com/asaskevich/govalidator && \
    go get -u github.com/sethgrid/curse && \
    GOOS=linux go build -ldflags="-w -s" && \
    adduser -D -g '' mppinger

FROM alpine:3.8

COPY --from=builder /mpping/mpping /app/mpping
COPY --from=builder /etc/passwd /etc/passwd

USER mppinger

WORKDIR /app

ENTRYPOINT [ "/app/mpping" ]