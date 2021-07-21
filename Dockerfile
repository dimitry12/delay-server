FROM golang:alpine as builder
WORKDIR /app
COPY  . /app/
RUN GO111MODULE=auto CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/delay-server

FROM scratch
COPY --from=builder /go/bin/delay-server /go/bin/delay-server
EXPOSE 8080
CMD ["/go/bin/delay-server"]
