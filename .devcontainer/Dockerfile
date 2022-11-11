FROM golang:latest

RUN apt-get -y update
RUN apt-get -y install vim git
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 
RUN apt-get -y install protobuf-compiler
