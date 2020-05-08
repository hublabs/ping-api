FROM pangpanglabs/golang:builder-beta AS builder
WORKDIR /go/src/github.com/hublabs/ping-api
ADD . /go/src/github.com/hublabs/ping-api
ENV CGO_ENABLED=0
RUN go build -o ping-api

FROM pangpanglabs/alpine-ssl
WORKDIR /go/src/github.com/hublabs/ping-api
COPY --from=builder /go/src/github.com/hublabs/ping-api/ping-api /go/src/github.com/hublabs/ping-api/
EXPOSE 8000
CMD ["./ping-api"]