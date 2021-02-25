FROM golang AS builder
COPY . /app
WORKDIR /app
RUN go get -u google.golang.org/grpc && \
    go get -u github.com/golang/protobuf/protoc-gen-go

FROM alpine:latest
WORKDIR /root
RUN apk --no-cache add ca-certificates
COPY --from=builder /app .

EXPOSE 50051

ENTRYPOINT cd /go/src/grpc-test3/greet/greet_server && go run client.go