
FROM golang:latest as builder

WORKDIR /go/src/ygo
COPY . .

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
# RUN go env -w CGO_ENABLED=0


RUN go mod tidy
RUN go build -o server

FROM golang:latest

WORKDIR /go/src/ygo

COPY --from=0 /go/src/ygo/server .

EXPOSE 8888

ENTRYPOINT ./server 