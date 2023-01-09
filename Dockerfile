
FROM golang:latest

WORKDIR /go/src/ygo
COPY . .

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod tidy
RUN go build -o server

EXPOSE 8888

ENTRYPOINT ./server app start
