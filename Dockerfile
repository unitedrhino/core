FROM golang:1.19-alpine3.16 as go-builder
WORKDIR /ithings/
COPY ./go.mod ./go.mod
RUN go mod download
COPY ./ ./
RUN cd ./service/apisvr && go mod tidy && go build .

FROM alpine:3.16  as web-builder
WORKDIR /ithings/
RUN mkdir front
RUN cd front&& wget -O front.tgz ${frontFIle}
RUN cd front && tar -xvzf front.tgz
RUN cd front && rm -rf front.tgz


FROM alpine:3.16
LABEL homepage="https://github.com/i-Things/iThings"
ENV TZ Asia/Shanghai
RUN apk add tzdata

WORKDIR /ithings/
COPY --from=go-builder /ithings/service/apisvr/apisvr ./apisvr
COPY --from=go-builder /ithings/service/deploy/conf/core/etc/ ./etc
COPY --from=web-builder /ithings/assets/front/* ./dist/app

ENTRYPOINT ["./apisvr"]
