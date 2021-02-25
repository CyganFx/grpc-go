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

ENTRYPOINT cd /go/src/grpc-test3/greet/greet_server && go run server.go

#FROM golang AS builder
#COPY . /app
#COPY ../../go.mod /app
#COPY ../../go.sum /app
#COPY ../greetpb /app
#WORKDIR /app
#RUN go mod download
#RUN CGO_ENABLED=0 GOOS=linux go build -o myApp ../../greet
#
#FROM alpine:latest
#WORKDIR /root
#RUN apk --no-cache add ca-certificates
#COPY --from=builder /app .
#
#CMD["./myApp"]